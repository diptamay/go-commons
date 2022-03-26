package secrets

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/diptamay/go-commons/helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var SecretsDir = "/run/secrets"
var SecretsInitialized = false
var Secrets *map[string]*Secret

func init() {
	secretDir := os.Getenv("SECRETS_DIRECTORY")
	if secretDir != "" {
		SecretsDir = secretDir
	}
}

type Secret struct {
	Version interface{} `json:"version"`
	Name    string      `json:"name"`
	Value   string      `json:"value"`
}

func fetchSecretFile(filename string, fileChan chan *Secret) {
	fp := filepath.Join(SecretsDir, filename)
	filedata, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Printf("error in ioutil.readFile: %#v", err.Error())
		fileChan <- nil
		return
	}
	secret := new(Secret)
	err = json.Unmarshal(filedata, secret)
	if err != nil {
		fmt.Printf("error in ioutil.readFile: %#v", err.Error())
		fileChan <- nil
	} else {
		fileChan <- secret
	}
}

func FetchSecrets() (*map[string]*Secret, error) {
	files, err := ioutil.ReadDir(SecretsDir)
	if err != nil {
		return nil, fmt.Errorf("Secrets directory is not mounted: %#v", err.Error())
	}
	filedatas := make(chan *Secret, len(files))
	for _, file := range files {
		go fetchSecretFile(file.Name(), filedatas)
	}
	Secrets = &map[string]*Secret{}
	SecretsInitialized = true
	for s := 0; s < len(files); s++ {
		secret := <-filedatas
		if secret != nil && secret.Name != "" {
			(*Secrets)[secret.Name] = secret
		}
	}
	return Secrets, nil
}

func GetSecret(secretName string) (*Secret, error) {
	if !SecretsInitialized {
		FetchSecrets()
	}
	if Secrets != nil {
		if secret, ok := (*Secrets)[secretName]; ok {
			return secret, nil
		}
	}
	return nil, fmt.Errorf("Secret for secretName %#v does not exist", secretName)
}

/**
The following code is for AWS secret manager.
For now, separate the 2 ways of get secrets now. Cz the env above appbuild hasn't moved to k8s yet.
TODO: When every env has k8s, remove the old above getSecret() and fetchSecret() APIs of from runtime mounted disk.
*/

// this secretsMap is for the AWS secret manager version.
var secretsMap = map[string]*Secret{}

func GetSecretsFromAWSSecretManager(secretName string) (*Secret, error) {
	//cache the secret, avoid duplicate call to AWS secret manager to get a same secret
	if secret, ok := secretsMap[secretName]; ok {
		return secret, nil
	} else {
		if secret, err := fetchSecretsFromAWSSecretManager(secretName); err != nil {
			return nil, err
		} else {
			//cache
			secretsMap[secretName] = secret
			return secret, nil
		}
	}
}

//determine the secret name the runtime needed based on different logiflows service running which DL is attaching to.
//sample: the secret of logiflows-service: "service.ENCRYPTION_KEY"
func DetermineSecretName() string {
	secretName := os.Getenv("SECRET_NAME")
	if secretName != "" {
		return secretName
	}
	microSvcNameRunningDL := os.Getenv("SERVICE_8080_NAME")
	if microSvcNameRunningDL == "" {
		microSvcNameRunningDL = "service"
	}
	secretName = microSvcNameRunningDL + ".ENCRYPTION_KEY"
	return secretName
}

func fetchSecretsFromAWSSecretManager(secretName string) (*Secret, error) {
	//Create a Secrets Manager client
	sess, err := session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	})
	if err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("Going to fetch the secret of the secretName %#v, from AWS secrets manager", secretName), new(map[string]interface{}))

	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(helpers.GetAWSRegion()))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		var decodedBinarySecret string
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Println(fmt.Sprintf("Base64 Decode Error when decode the binary result of secrets getting back from AWS secret manager: %#v", err), new(map[string]interface{}))
			return nil, err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		secretString = decodedBinarySecret
	}

	// Currently secrets on AWS secret manager is in string format. Unmarshal the string to Secret struct.
	secretObj := Secret{}
	err = json.Unmarshal([]byte(secretString), &secretObj)
	if err != nil {
		log.Println(fmt.Sprintf("error unmarshaling secret: %#v", err), new(map[string]interface{}))
		return nil, err
	}
	return &secretObj, nil
}
