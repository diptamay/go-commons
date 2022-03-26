package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
)

const (
	GCMTagSize     = 16
	IVSize         = 12
	TokenSeparator = "|$|"
)

func getIV(value string) ([]byte, []byte, []byte) {
	var encrypted, iv, tag []byte
	for index, v := range strings.Split(value, TokenSeparator) {
		switch index {
		case 0:
			encrypted, _ = base64.StdEncoding.DecodeString(v)
		case 1:
			iv, _ = base64.StdEncoding.DecodeString(v)
		case 2:
			tag, _ = base64.StdEncoding.DecodeString(v)
		}
	}
	return encrypted, iv, tag
}

func setIV(value []byte, iv []byte, tag []byte) string {
	return base64.StdEncoding.EncodeToString(value) + TokenSeparator + base64.StdEncoding.EncodeToString(iv) + TokenSeparator + base64.StdEncoding.EncodeToString(tag)
}

type CryptKeeperInterface interface {
	EncryptPayload(*map[string]interface{}, *[]string) (*map[string]interface{}, error)
	DecryptPayload(*map[string]interface{}) (*map[string]interface{}, error)
	Encrypt([]byte) (string, error)
	Decrypt(string) ([]byte, error)
}

type CryptKeeper struct {
	Cipher cipher.Block
}

func (keeper *CryptKeeper) Encrypt(toEnc []byte) (string, error) {
	iv := make([]byte, IVSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	gcmCipher, err := cipher.NewGCM(keeper.Cipher)
	if err != nil {
		return "", err
	}
	encrypted := gcmCipher.Seal(nil, iv, toEnc, nil)
	tag := encrypted[len(encrypted)-GCMTagSize:]
	encrypted = encrypted[:len(encrypted)-GCMTagSize]
	return setIV(encrypted, iv, tag), nil
}

func (keeper *CryptKeeper) EncryptPayload(payload *map[string]interface{}, whitelist *[]string) (*map[string]interface{}, error) {
	uniqueKeys := &map[string]bool{}
	for _, key := range *whitelist {
		if _, ok := (*uniqueKeys)[key]; !ok {
			(*uniqueKeys)[key] = true
		}
	}
	toEncrypt := &map[string]interface{}{}
	finalPayload := &map[string]interface{}{}
	for key, value := range *payload {
		if _, ok := (*uniqueKeys)[key]; ok {
			(*finalPayload)[key] = value
		} else {
			(*toEncrypt)[key] = value
		}
	}
	if len(*toEncrypt) > 0 {
		jsonstr, err := json.Marshal(*toEncrypt)
		if err != nil {
			return nil, err
		}
		encrypted, err := keeper.Encrypt(jsonstr)
		if err != nil {
			return nil, err
		}
		(*finalPayload)["ENCRYPTED_PAYLOAD"] = encrypted
	}
	return finalPayload, nil
}

func (keeper *CryptKeeper) Decrypt(value string) ([]byte, error) {
	toDecrypt, iv, tag := getIV(value)
	gcmDecipher, err := cipher.NewGCM(keeper.Cipher)
	if err != nil {
		return []byte{}, err
	}
	decrypted, err := gcmDecipher.Open(nil, iv, append(toDecrypt, tag...), nil)
	if err != nil {
		return []byte{}, err
	}
	return decrypted, nil
}

func (keeper *CryptKeeper) DecryptPayload(payload *map[string]interface{}) (*map[string]interface{}, error) {
	if _, ok := (*payload)["ENCRYPTED_PAYLOAD"]; !ok {
		return payload, nil
	}
	toDecrypt := (*payload)["ENCRYPTED_PAYLOAD"]
	decrypted, err := keeper.Decrypt(toDecrypt.(string))
	if err != nil {
		return nil, err
	}
	result := &map[string]interface{}{}
	err = json.Unmarshal(decrypted, result)
	if err != nil {
		return nil, err
	}
	final := *payload
	for key, value := range *result {
		final[key] = value
	}
	delete(final, "ENCRYPTED_PAYLOAD")
	return &final, nil
}

func makeCipher(key string) (cipher.Block, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return aes.NewCipher(decoded)
}

func MakeCryptKeeper(key string) (*CryptKeeper, error) {
	block, err := makeCipher(key)
	if err != nil {
		return nil, err
	}
	return &CryptKeeper{block}, nil
}
