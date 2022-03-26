# Go Commons

Repository of all common go modules used by go services.

### Go modules within repository:

- crypt
- doggie
- glogger
- instrumentation
- logSchema
- metrics
- secrets
- errors

### Go Module Details

#### crypt

Used to encrypt/decrypt sensitive information within documents being persisted to ElasticSearch and files being upload/downloaded to/from S3 buckets.
Includes Encrypt/Decrypt functions which encrypts/decrypts an entire string, and EncryptPayload/DecryptPayload functions which
encrypts/decrypts fields called ENCRYPTED_PAYLOAD within json files.

#### doggie

Datadog client used for sending metrics to Datadog dashboard on go services. Includes basic datadog functionalities such as
Histogram and Gauge, along with config and tag initialization within statsdconfig.go and statsdtags.go. Used along with send-stats.go in helpers to
push metrics up. Currently set to have a sampling rate of 0.1.

#### glogger

Standard logger used by go services. Includes a default schema and various log methods such as Info and Warn.

#### instrumentation

Used to produce memory usage, garbage collection information, and load data on go services. Uses datadog client to push metrics and therefore can only be initialized if
datadog client is initialized. Must be put into main.go after datadog client is created for it to run.

#### logSchema

Single struct used to create glogger in the correct format. The initialization of a glogger requires a logSchema, and therefore it
must be present for logger factories.

#### metrics

Convenience helper to send metrics and handle logging of errors. The initialization of a metrics optionally requires a doggie client and a glogger instance.
metrics allows `nil` DD tags and `nil` logging fields to be passed in all function calls.

#### secrets

Used to fetch secrets from directories on different environments as well as from AWS secrets manager. These secrets are required
for establishing sessions with AWS and therefore mandatory for any go service using AWS services (including ElasticSearch).

#### s3buckets

Used to initialize S3 and AWS sessions in go services. Contains methods such as downloading, uploading, deletion of S3 objects, as well as creation of S3 buckets.

#### Contributing

* Ensure githooks are installed by running `make init-git-hooks`.

#### Testing

unit tests: `make unit`
unit tests with verbosity: `make unit-verbose`

### Steps to update dependencies using go mod

1. ```go mod tidy``` to pull newest dependencies or new dependencies
2. Double check go.mod looks correct
3. ```go mod vendor``` to bring dependencies into vendor folder
holl

### Steps to update module to artifactory

#### For first time set-up reference:

First time setup:     ```jfrog rt config```      (It sets up the basic config for the local jfrog cli on your machine)
It goes through following prompts

    ```
    Artifactory server ID [Default-Server]:   <press enter>

    Artifactory URL: enter <YOUR_ARTIFACTOR_URL> . (with out tags)

    Access token (Leave blank for username and password/API key):  <press enter>

    User: <your username>

    Password/API key:  <your password>
    ```

#### Publishing steps

After your changes have been merged from branch to master, do the following steps.

1. cd to the root directory of local repository.

2. ```git status``` to see that it's on master branch and up to date with remote

3. ```git tag <tag_name> .``` to create the tag

4. ```git push origin <tag_name>```

5. run script ```sh publish.sh```
