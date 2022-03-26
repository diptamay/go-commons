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
	dockerHostnameFilePath = "/etc/docker_hostname"
	dockerHostname         *string
	fullName               *string
	imageName              *string
	imageOwner             *string
	workstream             *string
)

var envValueMap = map[string]func(string, *[4]string){
	"image_full_name": func(value string, result *[4]string) {
		fullName = &value
		result[1] = value
	},
	"image_name": func(value string, result *[4]string) {
		imageName = &value
		result[0] = value
	},
	"image_owner": func(value string, result *[4]string) {
		imageOwner = &value
		result[2] = value
	},
	"workstream": func(value string, result *[4]string) {
		workstream = &value
		result[3] = value
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

func getTagValuesFromEnv() *[4]string {
	var result [4]string
	if fullName != nil && imageName != nil && imageOwner != nil && workstream != nil {
		result = [4]string{
			*imageName,
			*fullName,
			*imageOwner,
			*workstream,
		}
	} else {
		result = [4]string{}
		for _, env := range os.Environ() {
			keyval := strings.Split(env, "=")
			if setter, ok := envValueMap[keyval[0]]; ok {
				setter(keyval[1], &result)
			}
			if fullName != nil && imageName != nil && imageOwner != nil && workstream != nil {
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
		"image.full_name:%s",
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
