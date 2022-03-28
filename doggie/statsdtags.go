package doggie

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	dockerHostnameFilePath = "/etc/hostname"
	dockerHostname         *string
	imageName              *string
	imageOwner             *string
	workstream             *string
)

var envValueMap = map[string]func(string, *[3]string){
	"image_name": func(value string, result *[3]string) {
		imageName = &value
		result[0] = value
	},
	"image_owner": func(value string, result *[3]string) {
		imageOwner = &value
		result[1] = value
	},
	"workstream": func(value string, result *[3]string) {
		workstream = &value
		result[2] = value
	},
}

func getDockerHostname() string {
	if dockerHostname != nil {
		return *dockerHostname
	}
	file, err := os.Open(dockerHostnameFilePath)
	if err != nil {
		log.Println("failed to read docker hostname file", err)
		return "unknown"
	}
	defer file.Close()
	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("failed to read docker hostname file", err)
		return "unknown"
	}
	dockerHost := strings.TrimSpace(string(filedata))
	dockerHostname = &dockerHost
	return *dockerHostname
}

func getTagValuesFromEnv() *[3]string {
	var result [3]string
	if imageName != nil && imageOwner != nil && workstream != nil {
		result = [3]string{
			*imageName,
			*imageOwner,
			*workstream,
		}
	} else {
		result = [3]string{}
		for _, env := range os.Environ() {
			keyval := strings.Split(env, "=")
			if setter, ok := envValueMap[keyval[0]]; ok {
				setter(keyval[1], &result)
			}
			if imageName != nil && imageOwner != nil && workstream != nil {
				break
			}
		}
	}
	return &result
}

func getDefaultTags(wg *sync.WaitGroup) *[]string {
	defer wg.Done()
	result := []string{}
	tagNames := []string{
		"image.name:%s",
		"image.owner:%s",
		"workstream:%s",
	}
	for index, tag := range *getTagValuesFromEnv() {
		if len(tag) > 0 {
			result = append(result, fmt.Sprintf(tagNames[index], tag))
		}
	}
	dhostname := getDockerHostname()
	dockerHostname = &dhostname
	processHostname, _ := os.Hostname()
	tags := []string{
		fmt.Sprintf("docker_host:%s", *dockerHostname),
		fmt.Sprintf("docker_instance:%s", processHostname),
	}
	tags = append(tags, result...)
	return &tags
}
