package crypt

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	Chance "github.com/ZeFort/chance"
	"github.com/stretchr/testify/assert"
)

var TestSecretKey string

func init() {
	hexbytes, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	TestSecretKey = base64.StdEncoding.EncodeToString(hexbytes)
}

func TestGetIV(t *testing.T) {
	chance := Chance.New()
	value := base64.StdEncoding.EncodeToString([]byte(chance.Word()))
	iv := base64.StdEncoding.EncodeToString([]byte(chance.String()))
	tag := base64.StdEncoding.EncodeToString([]byte(chance.String()))
	val1, val2, val3 := getIV(value + TokenSeparator + iv + TokenSeparator + tag)

	assert.Equal(t, value, base64.StdEncoding.EncodeToString(val1), "should parse value out of string and return base64 encoded byte array")
	assert.Equal(t, iv, base64.StdEncoding.EncodeToString(val2), "should parse value out of string and return base64 encoded byte array")
	assert.Equal(t, tag, base64.StdEncoding.EncodeToString(val3), "should parse value out of string and return base64 encoded byte array")
}

func TestSetIV(t *testing.T) {
	chance := Chance.New()
	value, _ := base64.StdEncoding.DecodeString(chance.Word())
	iv, _ := base64.StdEncoding.DecodeString(chance.String())
	tag, _ := base64.StdEncoding.DecodeString(chance.String())
	result := setIV(value, iv, tag)
	assert.Equal(
		t,
		base64.StdEncoding.EncodeToString(value)+TokenSeparator+base64.StdEncoding.EncodeToString(iv)+TokenSeparator+base64.StdEncoding.EncodeToString(tag),
		result,
		fmt.Sprintf("should create a %s delimited string", TokenSeparator),
	)
}

func TestMakeCryptKeeper(t *testing.T) {
	_, err := MakeCryptKeeper(string(TestSecretKey))
	assert.Equal(t, nil, err, "No error is thrown")
	_, err = MakeCryptKeeper("abc")
	assert.NotEqual(t, nil, err, "Error is thrown due to invalid secret key")
}

func TestEncryptAndDecrypt(t *testing.T) {
	chance := Chance.New()
	crypter, err := MakeCryptKeeper(string(TestSecretKey))
	if err != nil {
		panic(err)
	}
	toEncrypt := []byte(chance.Word())
	encrypted, err := crypter.Encrypt(toEncrypt)
	if err != nil {
		panic(err)
	}
	decrypted, err := crypter.Decrypt(encrypted)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, toEncrypt, decrypted, "should be able to encrypt a byte array and decrypt back to original value")
}

func TestEncryptAndDecryptErrorHandling(t *testing.T) {
	chance := Chance.New()
	crypter, err := MakeCryptKeeper(string(TestSecretKey))
	if err != nil {
		panic(err)
	}
	toEncrypt := []byte(chance.Word())
	encrypted, err := crypter.Encrypt(toEncrypt)

	hexbytes, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726566")
	fakeSecretKey := base64.StdEncoding.EncodeToString(hexbytes)
	crypter.Cipher, err = makeCipher(fakeSecretKey)
	if err != nil {
		panic(err)
	}

	_, err = crypter.Decrypt(encrypted)
	assert.NotEqual(t, nil, err, "An error should occur due to incorrect cipher for decrypt")
}

func TestEncryptAndDecryptPayload(t *testing.T) {
	chance := Chance.New()
	crypter, err := MakeCryptKeeper(string(TestSecretKey))
	if err != nil {
		panic(err)
	}
	payload := &map[string]interface{}{
		"should-encrypt": map[string]interface{}{
			chance.Word(): chance.String(),
		},
		"should-not-encrypt": map[string]interface{}{
			chance.Word(): chance.String(),
		},
		"also-should-not-encrypt": chance.Word(),
	}
	whitelist := &[]string{
		"should-not-encrypt",
		"also-should-not-encrypt",
	}
	encrypted, err := crypter.EncryptPayload(payload, whitelist)
	if err != nil {
		panic(err)
	}
	assert.IsType(t, string(""), (*encrypted)["ENCRYPTED_PAYLOAD"], "should have an encrypted payload")
	_, ok := (*encrypted)["should-encrypt"]
	assert.Equal(t, false, ok, "should not have \"should-encrypt\" field")
	_, ok = (*encrypted)["also-should-not-encrypt"]
	assert.Equal(t, true, ok, "should have \"also-should-not-encrypt\" field")
	_, ok = (*encrypted)["also-should-not-encrypt"]
	assert.Equal(t, true, ok, "should have \"also-should-not-encrypt\" field")

	decrypted, err := crypter.DecryptPayload(encrypted)
	_, ok = (*decrypted)["ENCRYPTED_PAYLOAD"]
	assert.Equal(t, false, ok, "should not have \"ENCRYPTED_PAYLOAD\" field")
	_, ok = (*decrypted)["should-encrypt"]
	assert.Equal(t, true, ok, "should have \"should-encrypt\" field")

	payload = &map[string]interface{}{
		"hello": "world",
	}
	whitelist = &[]string{
		"hello",
	}
	encrypted, err = crypter.EncryptPayload(payload, whitelist)
	_, ok = (*encrypted)["ENCRYPTED_PAYLOAD"]
	assert.Equal(t, false, ok, "should not have \"ENCRYPTED_PAYLOAD\" field")
}
