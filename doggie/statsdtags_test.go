package doggie

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	Chance "github.com/ZeFort/chance"
	"github.com/stretchr/testify/assert"
)

var (
	hostfileContent string
)

func makeTestDockerHostFile() {
	chance := Chance.New()
	hostfileContent = chance.Word()
	if _, err := os.Stat(filepath.Dir(dockerHostnameFilePath)); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Dir(dockerHostnameFilePath), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(dockerHostnameFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write([]byte(hostfileContent))
	if err != nil {
		panic(err)
	}
	err = file.Sync()
	if err != nil {
		panic(err)
	}
}

func removeTestDockerHostFile(fp string) {
	err := os.RemoveAll(filepath.Dir(dockerHostnameFilePath))
	if err != nil {
		panic(err)
	}
	dockerHostnameFilePath = fp
}

func TestGetDockerHostname(t *testing.T) {
	dockerHostname = nil
	originalDockerHostnameFilepath := dockerHostnameFilePath
	dockerHostnameFilePath = "test/docker_hostname"
	assert.Equal(t, "unknown", getDockerHostname(), "should return \"unknown\" when docker host file does not exist")

	makeTestDockerHostFile()

	assert.Equal(t, hostfileContent, getDockerHostname(), "should return content of docker host file")
	assert.Equal(t, *dockerHostname, getDockerHostname(), "should return assigned value of dockerHostname pointer")

	removeTestDockerHostFile(originalDockerHostnameFilepath)
}

func setEnvVars() *map[string]string {
	chance := Chance.New()
	result := make(map[string]string)
	vars := []string{
		"image_full_name",
		"image_name",
		"image_owner",
		"workstream",
	}

	for _, v := range vars {
		value := chance.Word()
		err := os.Setenv(v, value)
		if err != nil {
			panic(err)
		}
		result[v] = value
	}
	return &result
}

func TestGetTagValuesFromEnv(t *testing.T) {
	chance := Chance.New()
	testValue := chance.Word()
	fullName = &testValue
	imageName = &testValue
	imageOwner = &testValue
	workstream = &testValue
	expectedValues := []string{
		*imageName,
		*fullName,
		*imageOwner,
		*workstream,
	}
	tags := getTagValuesFromEnv()
	for index, tag := range *tags {
		assert.Equal(t, expectedValues[index], tag, "should resolve tag values if they are already defined")
	}
	fullName = nil
	imageName = nil
	imageOwner = nil
	workstream = nil

	vars := *setEnvVars()
	tags = getTagValuesFromEnv()
	assertionTags := *tags

	assert.Equal(
		t,
		vars["image_name"],
		assertionTags[0],
		fmt.Sprintf("imageName should equal %s", vars["image_name"]),
	)
	assert.Equal(
		t,
		vars["image_full_name"],
		assertionTags[1],
		fmt.Sprintf("fullName should equal %s", vars["image_full_name"]),
	)
	assert.Equal(
		t,
		vars["image_owner"],
		assertionTags[2],
		fmt.Sprintf("imageOwner should equal %s", vars["image_owner"]),
	)
	assert.Equal(
		t,
		vars["workstream"],
		assertionTags[3],
		fmt.Sprintf("workstream should equal %s", vars["workstream"]),
	)
}

func TestGetDefaultTags(t *testing.T) {
	var wg sync.WaitGroup
	tags := make(chan *[]string)
	wg.Add(1)
	go func() {
		tags <- getDefaultTags(&wg)
	}()
	wg.Wait()
	processHostname, _ := os.Hostname()
	tagList := []string{
		fmt.Sprintf("docker_host:%s", *dockerHostname),
		fmt.Sprintf("docker_instance:%s", processHostname),
		fmt.Sprintf("image.name:%s", *imageName),
		fmt.Sprintf("image.full_name:%s", *fullName),
		fmt.Sprintf("image.owner:%s", *imageOwner),
		fmt.Sprintf("workstream:%s", *workstream),
	}
	actualTagsPntr := <-tags
	actualTags := *actualTagsPntr
	for index, tag := range tagList {
		assert.Equal(
			t,
			tag,
			actualTags[index],
			"should generate default tag values as expected",
		)
	}
}
