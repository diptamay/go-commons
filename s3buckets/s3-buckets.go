package s3buckets

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/diptamay/go-commons/crypt"
	"github.com/diptamay/go-commons/helpers"
	"github.com/pkg/errors"
)

type UploadFile func(ctx context.Context, filekey string, contents []byte, tags *string) (string, error)
type UploadAllFiles func(ctx context.Context, filekeys *map[string]string) ([]string, error)
type DownloadFile func(ctx context.Context, filekey string, crypterOldKey crypt.CryptKeeperInterface) ([]byte, error)
type DownloadAllFiles func(ctx context.Context, filekeys *[]string) (*map[string]string, error)
type GetBucketObjectsTimeInterval func(ctx context.Context, prefix *string, startTime time.Time, endTime time.Time) ([]string, error)
type CopyObjectInS3 func(ctx context.Context, filekey string, targetKey string) error
type InitS3Bucket func(ctx context.Context, bucketCfg *S3BucketConfig) error

var (
	bucketName                    *string
	S3Session                     *s3.S3
	Upload                        UploadFile
	Download                      DownloadFile
	GetKeysPerInterval            GetBucketObjectsTimeInterval
	CopyKeysInBucket              CopyObjectInS3
	InitializeS3Bucket            InitS3Bucket
	S3ClientIsNotInitializedError = errors.New("S3 client is not initialized")
	ErrDecryptFail                = errors.New("Decryption failed")
)

type S3BucketConfig struct {
	Name                *string
	S3LocalstackAddress *string
}

type UploaderInterface interface {
	UploadWithContext(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type DownloaderInterface interface {
	DownloadWithContext(aws.Context, io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
}

func makeUploader(uploader UploaderInterface, crypter crypt.CryptKeeperInterface) UploadFile {
	return func(ctx context.Context, filekey string, contents []byte, tags *string) (string, error) {
		var encrypted []byte
		var err error
		encryptedBytes, err := crypter.Encrypt(contents)
		encrypted = []byte(encryptedBytes)
		if err != nil {
			return "", err
		}
		s3Input := &s3manager.UploadInput{
			Bucket: bucketName,
			Key:    aws.String(filekey),
			Body:   bytes.NewReader(encrypted),
		}

		if tags != nil {
			s3Input.Tagging = tags
		}

		// We give control to a timeout to the http client
		result, err := uploader.UploadWithContext(ctx, s3Input)
		if err != nil {
			return "", err
		}
		return result.Location, nil
	}
}

func getS3BucketSession(bucketConfig *S3BucketConfig) (*session.Session, error) {
	config := &aws.Config{
		Region:     aws.String(helpers.GetAWSRegion()),
		HTTPClient: helpers.NewHTTPClientRecommended(),
	}

	if bucketConfig.S3LocalstackAddress != nil && *bucketConfig.S3LocalstackAddress != "" {
		log.Printf("The AWS credentials are pointing to localstack. S3 endpoint: %s\n", *bucketConfig.S3LocalstackAddress)
		config.Endpoint = bucketConfig.S3LocalstackAddress
		config.S3ForcePathStyle = aws.Bool(true)
		// to bypass x509 cert validation error for localstack
		config.HTTPClient = helpers.NewHTTPClientInsecure()
	}

	// For debugging http issues to s3
	if os.Getenv("S3_HTTP_DEBUG") == "true" {
		config.LogLevel = aws.LogLevel(aws.LogDebugWithRequestErrors)
	}

	return session.NewSession(config)
}

func makeDownloader(downloader DownloaderInterface, crypter crypt.CryptKeeperInterface) DownloadFile {
	return func(ctx context.Context, filekey string, crypterOldKey crypt.CryptKeeperInterface) ([]byte, error) {
		writer := &aws.WriteAtBuffer{}
		_, err := downloader.DownloadWithContext(ctx, writer, &s3.GetObjectInput{
			Bucket: bucketName,
			Key:    aws.String(filekey),
		})

		if err != nil {
			return []byte{}, err
		}

		// For Debug documents that end with .enc.json - we do not decrypt as they are currently no encrypted with client-side encryption
		if match, matcherr := regexp.MatchString(".+\\.enc\\.json$", filekey); matcherr == nil && match {
			return writer.Bytes(), nil
		}

		if crypterOldKey != nil {
			content, err := crypterOldKey.Decrypt(string(writer.Bytes()))
			if err != nil {
				return []byte{}, ErrDecryptFail
			}
			return content, nil
		}
		content, err := crypter.Decrypt(string(writer.Bytes()))
		if err != nil {
			return []byte{}, ErrDecryptFail
		}
		return content, nil
	}
}

func makeCopyObjectInS3(session s3iface.S3API) CopyObjectInS3 {
	return func(ctx context.Context, sourceKey string, targetKey string) error {
		// The name of the source bucket and key name of the source object, separated by a slash (/)
		source := fmt.Sprint(*bucketName, "/", sourceKey)
		_, err := session.CopyObjectWithContext(ctx,
			&s3.CopyObjectInput{
				Bucket:     bucketName,
				Key:        aws.String(targetKey),
				CopySource: aws.String(source),
			})
		if err != nil {
			log.Println("Something went wrong with copying ", err)
			return err
		}
		return nil
	}
}

func DeleteObjFromS3(ctx context.Context, obj string) error {
	_, err := S3Session.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{Bucket: bucketName, Key: aws.String(obj)})
	if err != nil {
		log.Println("Something went wrong with deletion", err)
		return err
	}

	err = S3Session.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: bucketName,
		Key:    aws.String(obj),
	})

	if err != nil {
		return err
	}

	return nil
}

func GetBucketObjects(ctx context.Context, prefix *string) ([]string, error) {
	query := &s3.ListObjectsV2Input{
		Bucket: bucketName,
		Prefix: prefix,
	}

	truncatedListing := true
	var result []string

	for truncatedListing {
		response, err := S3Session.ListObjectsV2WithContext(ctx, query)

		if err != nil {
			log.Println("error fetching list of objects in s3bucket", err)
			return result, err
		}

		for _, file := range response.Contents {
			result = append(result, *file.Key)
		}

		// Set continuation token
		query.ContinuationToken = response.NextContinuationToken
		truncatedListing = *response.IsTruncated
	}

	return result, nil
}

// This function gets metadata of all the elements from a bucket and only returns the ones
// that fall inside the passed in time values
func makeGetBucketObjectsTimeInterval(session s3iface.S3API) GetBucketObjectsTimeInterval {
	return func(ctx context.Context, prefix *string, startTime time.Time, endTime time.Time) ([]string, error) {
		query := &s3.ListObjectsV2Input{
			Bucket: bucketName,
			Prefix: prefix,
		}

		truncatedListing := true
		var result []string

		for truncatedListing {
			response, err := session.ListObjectsV2WithContext(ctx, query)

			if err != nil {
				log.Println("error fetching list of objects in s3bucket", err)
				return result, err
			}

			for _, file := range response.Contents {
				fileTime := *file.LastModified
				if fileTime.Equal(startTime) || fileTime.Equal(endTime) ||
					(fileTime.After(startTime) && fileTime.Before(endTime)) {
					result = append(result, *file.Key)
				}
			}

			// Set continuation token
			query.ContinuationToken = response.NextContinuationToken
			truncatedListing = *response.IsTruncated
		}

		return result, nil
	}
}

func InitializeS3Handlers(ctx context.Context, bucketCfg *S3BucketConfig, crypter *crypt.CryptKeeper) error {
	bucketName = bucketCfg.Name
	enabledSession, err := getS3BucketSession(bucketCfg)
	if err != nil {
		return err
	}
	awsSession := session.Must(enabledSession, nil)
	S3Session = s3.New(awsSession)

	InitializeS3Bucket = makeInitializeS3Bucket(S3Session)
	if err := InitializeS3Bucket(ctx, bucketCfg); err != nil {
		return err
	}

	Upload = makeUploader(s3manager.NewUploader(awsSession), crypter)
	Download = makeDownloader(s3manager.NewDownloader(awsSession), crypter)
	GetKeysPerInterval = makeGetBucketObjectsTimeInterval(S3Session)
	CopyKeysInBucket = makeCopyObjectInS3(S3Session)

	return nil
}

func makeS3Bucket(ctx context.Context, session s3iface.S3API, bucket *string) (*s3.CreateBucketOutput, error) {
	return session.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: bucket,
	})
}

func makeInitializeS3Bucket(session s3iface.S3API) InitS3Bucket {
	return func(ctx context.Context, bucketCfg *S3BucketConfig) error {
		log.Println("Checking if  bucket ", *bucketCfg.Name, " exists")

		_, headBucketErr := session.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
			Bucket: bucketCfg.Name,
		})

		if headBucketErr != nil && bucketCfg.S3LocalstackAddress != nil && *bucketCfg.S3LocalstackAddress != "" {
			log.Println("Using localstack, and bucket does not exist, so creating new bucket ", *bucketCfg.Name)
			if _, createErr := makeS3Bucket(ctx, session, bucketCfg.Name); createErr != nil {
				log.Println("creating bucket failed,", *bucketCfg.Name, createErr.Error())
				return createErr
			}
			log.Println("created bucket,", *bucketCfg.Name)
		} else if headBucketErr != nil {
			return headBucketErr
		}
		return nil
	}
}

func DoesExist(ctx context.Context, prefix string) (bool, error) {
	if S3Session == nil {
		return false, S3ClientIsNotInitializedError
	}
	resp, err := S3Session.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: bucketName,
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return false, err
	}
	return len(resp.Contents) > 0, nil
}
