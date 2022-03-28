package doggie

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
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

func getDataDogConfig(wg *sync.WaitGroup) *dataDogConfig {
	defer wg.Done()
	if configuration == nil || dataDogHost == nil || dataDogPort == nil {
		hasValues := getValuesFromEnv()
		if !hasValues.hasPort {
			port := 8125
			dataDogPort = &port
		}
		if !hasValues.hasHost {
			host := "127.0.0.1"
			dataDogHost = &host
		}
		configuration = &dataDogConfig{
			*dataDogHost,
			*dataDogPort,
		}
	}
	return configuration
}
