package doggie

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	factsFilePath = "/etc/ansible/facts.d/aws.json"
	dataDogHost   *string
	dataDogPort   *int
	configuration *dataDogConfig
)

type dataDogConfig struct {
	Host string
	Port int
}

func getValuesFromEnv() *struct {
	hasPort bool
	hasHost bool
} {
	hasValues := new(struct {
		hasPort bool
		hasHost bool
	})
	for _, env := range os.Environ() {
		keyval := strings.Split(env, "=")
		if keyval[0] == "DOGSTATSD_HOST" {
			dataDogHost = &keyval[1]
			hasValues.hasHost = true
		} else if keyval[0] == "DOGSTATSD_PORT" {
			port, _ := strconv.Atoi(keyval[1])
			dataDogPort = &port
			hasValues.hasPort = true
		}
		if dataDogHost != nil && dataDogPort != nil {
			break
		}
	}
	return hasValues
}

func parseValuesFromFacts() bool {
	file, err := os.Open(factsFilePath)
	if err != nil {
		log.Println("aws facts file does not exist", err)
		return false
	}
	defer file.Close()
	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("failed to read aws facts file data", err)
		return false
	}
	config := make(map[string]interface{})
	err = json.Unmarshal(filedata, &config)
	if err != nil {
		log.Println("failed to read aws facts file data", err)
		return false
	}
	if _, ok := config["identity_document"]; ok {
		identityDocument := config["identity_document"].(map[string]interface{})
		if value, ok := identityDocument["privateIp"]; ok {
			host := value.(string)
			dataDogHost = &host
		} else {
			return false
		}
	} else {
		return false
	}
	return true
}

func getDataDogConfig(wg *sync.WaitGroup) *dataDogConfig {
	defer wg.Done()
	if configuration == nil || dataDogHost == nil || dataDogPort == nil {
		hasValues := getValuesFromEnv()
		if !hasValues.hasHost {
			hasHost := parseValuesFromFacts()
			if !hasHost {
				host := "router.service.consul"
				dataDogHost = &host
			}
		}
		if !hasValues.hasPort {
			port := 8125
			dataDogPort = &port
		}
		configuration = &dataDogConfig{
			*dataDogHost,
			*dataDogPort,
		}
	}
	return configuration
}
