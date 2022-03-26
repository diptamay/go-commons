package doggie

import (
	Chance "github.com/ZeFort/chance"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var (
	dogstatsdHost string
	dogstatdPort  int
)

func TestGetValuesFromEnv(t *testing.T) {
	chance := Chance.New()
	dogstatsdHost = chance.Word()
	dogstatdPort = chance.IntBtw(1000, 9999)
	os.Setenv("DOGSTATSD_HOST", dogstatsdHost)
	os.Setenv("DOGSTATSD_PORT", strconv.Itoa(dogstatdPort))

	dataDogHost = nil
	dataDogPort = nil

	env := getValuesFromEnv()

	assert.Equal(t, true, env.hasHost, "Host should be set from env DOGSTATSD_HOST value")
	assert.Equal(t, true, env.hasPort, "Port should be set from env DOGSTATSD_PORT value")

	assert.Equal(t, dogstatsdHost, *dataDogHost, "Host values should match")
	assert.Equal(t, dogstatdPort, *dataDogPort, "Port values should match")
}

func writeInvalidFacts() {
	chance := Chance.New()
	if _, err := os.Stat(filepath.Dir(factsFilePath)); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Dir(factsFilePath), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(factsFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(chance.String()))
	if err != nil {
		panic(err)
	}
	err = file.Sync()
	if err != nil {
		panic(err)
	}
}

func writeValidFacts(jsonstr []byte) {
	if _, err := os.Stat(filepath.Dir(factsFilePath)); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Dir(factsFilePath), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(factsFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(jsonstr)
	if err != nil {
		panic(err)
	}
	err = file.Sync()
	if err != nil {
		panic(err)
	}
}

func removeFacts() {
	err := os.RemoveAll(filepath.Dir(factsFilePath))
	if err != nil {
		panic(err)
	}
}

func TestParseValuesFromFacts(t *testing.T) {
	factsFilePath = "test/facts.json"
	assert.Equal(t, false, parseValuesFromFacts(), "should return false if facts file cant be opened")

	writeInvalidFacts()
	defer removeFacts()
	assert.Equal(t, false, parseValuesFromFacts(), "should return false if facts file is invalid json")
	removeFacts()
	json := `{
		"identity_document": {
			"privateIp": "127.0.0.1"
		}
	}`
	writeValidFacts([]byte(json))
	assert.Equal(t, true, parseValuesFromFacts(), "should be able to parse datadog host value from facts")
	assert.Equal(t, "127.0.0.1", *dataDogHost, "host should match value defined in facts file")
	removeFacts()

	json = `{
		"identity_document": {}
	}`
	writeValidFacts([]byte(json))
	assert.Equal(t, false, parseValuesFromFacts(), "should return false if \"privateIp\" value is not defined in facts")
	removeFacts()
	json = `{}`
	writeValidFacts([]byte(json))
	assert.Equal(t, false, parseValuesFromFacts(), "should return false if \"identity_document\" value is not defined in facts")
}
