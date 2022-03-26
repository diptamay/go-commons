package secrets

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"testing"

	Chance "github.com/ZeFort/chance"
	"github.com/stretchr/testify/assert"
)

var tempSecretsDir string

func init() {
	chance := Chance.New()
	dir, err := ioutil.TempDir("", "testsecrets")
	if err != nil {
		panic(err)
	}
	tempSecretsDir = dir
	SecretsDir = dir
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	filesToMake := [2][]byte{
		[]byte(chance.String()),
		[]byte(`{"version": "latest", "name": "ENCRYPTION_KEY", "value": "enckeyvalue"}`),
	}
	for _, filedata := range filesToMake {
		wg.Add(1)
		name := chance.Word()
		go func(d []byte, name string) {
			err := ioutil.WriteFile(path.Join(dir, name+".json"), d, 0777)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(filedata, name)
	}
	wg.Wait()
}

func TestFetchSecretFile(t *testing.T) {
	fileChan := make(chan *Secret)
	go fetchSecretFile("nonexistent.json", fileChan)
	result := <-fileChan
	assert.Equal(t, true, result == nil, "should result in nil when file does not exist")
}

func TestFetchSecrets(t *testing.T) {
	secrets, err := FetchSecrets()
	assert.IsType(t, Secret{}, (*(*secrets)["ENCRYPTION_KEY"]), "should return an instance of Secret when file is parsed successfully")
	if tempSecretsDir == "" {
		panic("temp secrets dir not created")
	}
	err = os.RemoveAll(tempSecretsDir)
	if err != nil {
		panic(err)
	}
	_, err = FetchSecrets()
	assert.Equal(t, true, err != nil, "should return an error when secrets direcrtory not mounted")
}

func TestGetSecret(t *testing.T) {
	chance := Chance.New()
	secret, err := GetSecret("ENCRYPTION_KEY")
	assert.IsType(t, Secret{}, *secret, "should retrieve secret when initialized")
	key := chance.Word()
	_, err = GetSecret(key)
	assert.Errorf(t, err, fmt.Sprintf("Secret for key %s does not exist", key))
}

func TestDetermineSecretName(t *testing.T) {
	os.Setenv("SECRET_NAME", "name")
	res := DetermineSecretName()
	assert.Equal(t, "name", res, "the returned secret name should equal to the environment variable SECRET_NAME")
	os.Setenv("SECRET_NAME", "")
	res = DetermineSecretName()
	assert.Equal(t, "service.ENCRYPTION_KEY", res, "the returned secret name should equal to the logiflows service encryption key")
	os.Setenv("SERVICE_8080_NAME", "name")
	res = DetermineSecretName()
	assert.Equal(t, "name.ENCRYPTION_KEY", res, "the returned secret name should equal to the environment variable plus encryption key")
}

func TestGetSecretFromAWSSecretManager(t *testing.T) {
	_, err := GetSecretsFromAWSSecretManager("cgi.ENCRYPTION_KEY")
	assert.Error(t, err, "Error should exist because the aws session will not work")
}
