package s3buckets

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"testing"
	"time"

	Chance "github.com/ZeFort/chance"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/diptamay/go-commons/crypt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	MockUpload          UploadFile
	MockDownload        DownloadFile
	MockGenericDownload DownloadFile
	Crypter             *crypt.CryptKeeper
)

type S3BucketsTestSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestS3Buckets(t *testing.T) {
	suite.Run(t, new(S3BucketsTestSuite))
}

func (suite *S3BucketsTestSuite) SetupTest() {
	hexbytes, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	testSecretKey := base64.StdEncoding.EncodeToString(hexbytes)
	crypter, err := crypt.MakeCryptKeeper(string(testSecretKey))
	if err != nil {
		panic(err)
	}
	Crypter = crypter
	bucketName = aws.String("logiflows-test")
}

type MockCrypter struct {
	mock.Mock
	crypt.CryptKeeperInterface
}

func (m *MockCrypter) Encrypt(encStr []byte) (string, error) {
	m.Called(encStr)
	return string(encStr), nil
}

type MockUploaderS3API struct {
	mock.Mock
}

func (m *MockUploaderS3API) UploadWithContext(ctx aws.Context, config *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	m.Called(ctx, config)
	return &s3manager.UploadOutput{
		Location: *config.Key,
	}, nil
}

type MockGetObjectS3API struct {
	s3iface.S3API
}

type MockGetObjectErrorS3API struct {
	s3iface.S3API
}

type MockGetObjectEmptyS3API struct {
	s3iface.S3API
}

type MockCopyObjectS3BucketAPI struct {
	s3iface.S3API
}

type MockCopyObjectErrS3API struct {
	s3iface.S3API
}

func (m *MockCopyObjectErrS3API) CopyObjectWithContext(ctx aws.Context, config *s3.CopyObjectInput, opts ...request.Option) (*s3.CopyObjectOutput, error) {
	output := &s3.CopyObjectOutput{}
	return output, errors.New("error in copyObject")
}

func (m *MockCopyObjectS3BucketAPI) CopyObjectWithContext(ctx aws.Context, config *s3.CopyObjectInput, opts ...request.Option) (*s3.CopyObjectOutput, error) {
	key := "new_key"
	var result s3.CopyObjectResult
	result.SetLastModified(time.Now())
	result.SetETag(key)
	return &s3.CopyObjectOutput{
		CopyObjectResult: &result,
	}, nil
}

func (m *MockGetObjectS3API) ListObjectsV2WithContext(ctx aws.Context, input *s3.ListObjectsV2Input, opts ...request.Option) (*s3.ListObjectsV2Output, error) {
	key := "key"
	truncated := false
	lastModifiedTime, _ := time.Parse(time.RFC3339, "2006-01-01T15:05:05Z")
	content := s3.Object{
		Key:          &key,
		LastModified: &lastModifiedTime,
	}
	var contents []*s3.Object
	contents = append(contents, &content)
	output := &s3.ListObjectsV2Output{
		Contents:    contents,
		IsTruncated: &truncated,
	}
	return output, nil
}

func (m *MockGetObjectErrorS3API) ListObjectsV2WithContext(ctx aws.Context, input *s3.ListObjectsV2Input, opts ...request.Option) (*s3.ListObjectsV2Output, error) {
	output := &s3.ListObjectsV2Output{}
	return output, errors.New("error in listObject")
}

func (m *MockGetObjectEmptyS3API) ListObjectsV2WithContext(ctx aws.Context, input *s3.ListObjectsV2Input, opts ...request.Option) (*s3.ListObjectsV2Output, error) {
	key := "key"
	truncated := false
	lastModifiedTime, _ := time.Parse(time.RFC3339, "2007-01-01T15:05:05Z")
	content := s3.Object{
		Key:          &key,
		LastModified: &lastModifiedTime,
	}
	var contents []*s3.Object
	contents = append(contents, &content)
	output := &s3.ListObjectsV2Output{
		Contents:    contents,
		IsTruncated: &truncated,
	}
	return output, nil
}

type MockInitBucketS3API struct {
	mock.Mock
	s3iface.S3API
}

func (m *MockInitBucketS3API) HeadBucketWithContext(ctx aws.Context, input *s3.HeadBucketInput, opts ...request.Option) (*s3.HeadBucketOutput, error) {
	m.Called(ctx, input)
	return &s3.HeadBucketOutput{}, nil
}

func (m *MockInitBucketS3API) CreateBucketWithContext(ctx aws.Context, input *s3.CreateBucketInput, opts ...request.Option) (*s3.CreateBucketOutput, error) {
	m.Called(ctx, input)
	return &s3.CreateBucketOutput{}, nil
}

type MockInitBucketHeadErrorS3API struct {
	mock.Mock
	s3iface.S3API
}

func (m *MockInitBucketHeadErrorS3API) HeadBucketWithContext(ctx aws.Context, input *s3.HeadBucketInput, opts ...request.Option) (*s3.HeadBucketOutput, error) {
	m.Called(ctx, input)
	return &s3.HeadBucketOutput{}, errors.New("Head bucket error")
}

func (m *MockInitBucketHeadErrorS3API) CreateBucketWithContext(ctx aws.Context, input *s3.CreateBucketInput, opts ...request.Option) (*s3.CreateBucketOutput, error) {
	m.Called(ctx, input)
	return &s3.CreateBucketOutput{}, nil
}

type MockInitBucketCreateErrorS3API struct {
	mock.Mock
	s3iface.S3API
}

func (m *MockInitBucketCreateErrorS3API) HeadBucketWithContext(ctx aws.Context, input *s3.HeadBucketInput, opts ...request.Option) (*s3.HeadBucketOutput, error) {
	m.Called(ctx, input)
	return &s3.HeadBucketOutput{}, errors.New("Head bucket error")
}

func (m *MockInitBucketCreateErrorS3API) CreateBucketWithContext(ctx aws.Context, input *s3.CreateBucketInput, opts ...request.Option) (*s3.CreateBucketOutput, error) {
	m.Called(ctx, input)
	return &s3.CreateBucketOutput{}, errors.New("Create bucket error")
}

type MockUploaderWithError struct {
	mock.Mock
}

func (m *MockUploaderWithError) UploadWithContext(ctx aws.Context, config *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	m.Called(ctx, config)
	return &s3manager.UploadOutput{}, errors.New("error in upload")
}

type MockDownloader struct {
	mock.Mock
}

func (m *MockDownloader) DownloadWithContext(ctx aws.Context, writer io.WriterAt, config *s3.GetObjectInput, opts ...func(*s3manager.Downloader)) (int64, error) {
	m.Called(ctx, writer, config)
	encrypted, err := Crypter.Encrypt([]byte(Chance.New().String()))
	if err != nil {
		panic(err)
	}
	writer.WriteAt([]byte(encrypted), int64(0))
	return int64(100), nil
}

type MockDownloaderWithError struct {
	mock.Mock
}

func (m *MockDownloaderWithError) DownloadWithContext(ctx aws.Context, writer io.WriterAt, config *s3.GetObjectInput, opts ...func(*s3manager.Downloader)) (int64, error) {
	m.Called(ctx, writer, config)
	return int64(0), errors.New("error in download")
}

func (suite *S3BucketsTestSuite) TestUpload() {
	chance := Chance.New()
	value := []byte("UNENCRYPTED_CONTENTS")

	mockCrypter := new(MockCrypter)
	mockCrypter.
		On("Encrypt", value).
		Return("ENCRYPTED_CONTENTS", nil)

	key := chance.Word()
	expected := &s3manager.UploadInput{
		Bucket: aws.String("logiflows-test"),
		Key:    aws.String(key),
		Body:   bytes.NewReader(value),
	}
	mockUploader := new(MockUploaderS3API)
	MockUpload := makeUploader(mockUploader, mockCrypter)
	mockUploader.
		On("UploadWithContext", context.Background(), expected).
		Return(mock.AnythingOfType("*s3manager.UploadOutput"), nil)
	location, err := MockUpload(context.Background(), key, value, nil)
	if err != nil {
		panic(err)
	}
	assert.Equal(suite.T(), key, location, "should use key for file location")
	mockUploader.AssertExpectations(suite.T())
}

func (suite *S3BucketsTestSuite) TestUploadWithError() {
	chance := Chance.New()
	key := chance.Word()
	value := []byte("UNENCRYPTED_CONTENTS")

	mockCrypter := new(MockCrypter)
	mockCrypter.
		On("Encrypt", value).
		Return("ENCRYPTED_CONTENTS", nil)

	mockUploadWithErr := new(MockUploaderWithError)
	MockUpload := makeUploader(mockUploadWithErr, Crypter)
	mockUploadWithErr.
		On("UploadWithContext", mock.Anything, mock.AnythingOfType("*s3manager.UploadInput")).
		Return(mock.AnythingOfType("*s3manager.UploadOutput"), mock.AnythingOfType("error"))
	location, err := MockUpload(context.Background(), key, value, nil)

	assert.Equal(suite.T(), "", location, "location should be nil")
	assert.Equal(suite.T(), "error in upload", err.Error(), "should surface an error in upload go routine")
}

func (suite *S3BucketsTestSuite) TestUploadWithTag() {
	chance := Chance.New()
	value := []byte("UNENCRYPTED_CONTENTS")

	mockCrypter := new(MockCrypter)
	mockCrypter.
		On("Encrypt", value).
		Return("ENCRYPTED_CONTENTS", nil)

	key := chance.Word()
	tag := "key=value"
	expected := &s3manager.UploadInput{
		Bucket:  aws.String("logiflows-test"),
		Key:     aws.String(key),
		Body:    bytes.NewReader(value),
		Tagging: aws.String(tag),
	}
	mockUploader := new(MockUploaderS3API)
	MockUpload := makeUploader(mockUploader, mockCrypter)
	mockUploader.
		On("UploadWithContext", context.Background(), expected).
		Return(mock.AnythingOfType("*s3manager.UploadOutput"), nil)
	location, err := MockUpload(context.Background(), key, value, &tag)
	if err != nil {
		panic(err)
	}
	assert.Equal(suite.T(), key, location, "should use key for file location")
	mockUploader.AssertExpectations(suite.T())
}

func (suite *S3BucketsTestSuite) TestDownload() {
	chance := Chance.New()
	mockDownload := new(MockDownloader)
	MockDownload := makeDownloader(mockDownload, Crypter)
	file := chance.Word()
	mockDownload.
		On("DownloadWithContext", mock.Anything, mock.AnythingOfType("*aws.WriteAtBuffer"), mock.AnythingOfType("*s3.GetObjectInput")).
		Return(mock.AnythingOfType("int64"), nil)

	result, _ := MockDownload(context.Background(), file, nil)

	assert.IsType(suite.T(), []byte{}, result, "should return a decrypted byte array")
}

func (suite *S3BucketsTestSuite) TestDownloadWithErr() {
	chance := Chance.New()
	file := chance.Word()
	mockDownload := new(MockDownloaderWithError)
	MockDownload := makeDownloader(mockDownload, Crypter)
	mockDownload.
		On("DownloadWithContext", mock.Anything, mock.AnythingOfType("*aws.WriteAtBuffer"), mock.AnythingOfType("*s3.GetObjectInput")).
		Return(mock.AnythingOfType("int64"), mock.AnythingOfType("error"))

	result, err := MockDownload(context.Background(), file, nil)

	assert.Equal(suite.T(), []byte{}, result, "should be an empty byte array")
	assert.Equal(suite.T(), "error in download", err.Error(), "should surface an error in download go routine")
}

func (suite *S3BucketsTestSuite) TestGetBucketObjectsTimeInterval() {
	prefix := "prefix"
	startTime, _ := time.Parse(time.RFC3339, "2006-01-01T15:04:05Z")
	endTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	mockGetObject := new(MockGetObjectS3API)
	mockObjectGetter := makeGetBucketObjectsTimeInterval(mockGetObject)
	response, _ := mockObjectGetter(context.Background(), &prefix, startTime, endTime)
	assert.Equal(suite.T(), []string{"key"}, response, "should return single key")
}

func (suite *S3BucketsTestSuite) TestGetBucketObjectsTimeIntervalError() {
	prefix := "prefix"
	startTime, _ := time.Parse(time.RFC3339, "2006-01-01T15:04:05Z")
	endTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	mockGetObjectErr := new(MockGetObjectErrorS3API)
	mockObjectGetter := makeGetBucketObjectsTimeInterval(mockGetObjectErr)
	_, err := mockObjectGetter(context.Background(), &prefix, startTime, endTime)
	assert.Equal(suite.T(), "error in listObject", err.Error(), "should surface an error in download go routine")
}

func (suite *S3BucketsTestSuite) TestGetBucketObjectsTimeIntervalEmpty() {
	prefix := "prefix"
	startTime, _ := time.Parse(time.RFC3339, "2006-01-01T15:04:05Z")
	endTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	mockGetObjectEmpty := new(MockGetObjectEmptyS3API)
	mockObjectGetter := makeGetBucketObjectsTimeInterval(mockGetObjectEmpty)
	response, _ := mockObjectGetter(context.Background(), &prefix, startTime, endTime)
	assert.Equal(suite.T(), []string(nil), response, "should return single key")
}

func (suite *S3BucketsTestSuite) TestCopyObjectInS3() {
	mockCopyObject := new(MockCopyObjectS3BucketAPI)
	prefix := "prefix"
	target := "target"
	mockCopyKeyInS3 := makeCopyObjectInS3(mockCopyObject)
	response := mockCopyKeyInS3(context.Background(), prefix, target)
	assert.Equal(suite.T(), nil, response, "should return nothing if successful")
}

func (suite *S3BucketsTestSuite) TestCopyObjectInS3WithError() {
	mockCopyObject := new(MockCopyObjectErrS3API)
	prefix := "prefix"
	target := "target"
	mockCopyKeyInS3 := makeCopyObjectInS3(mockCopyObject)
	err := mockCopyKeyInS3(context.Background(), prefix, target)
	assert.Equal(suite.T(), "error in copyObject", err.Error(), "should return an error")
}

func (suite *S3BucketsTestSuite) TestGetS3BucketSessionWithLocalstack() {
	expected := "https://localstack.service.consul:4572"

	result, _ := getS3BucketSession(&S3BucketConfig{
		S3LocalstackAddress: aws.String(expected),
	})

	assert.Equal(suite.T(), expected, *result.Config.Endpoint, "should have localstack as endpoint")
	assert.Equal(suite.T(), "us-east-1", *result.Config.Region, "should have correct region")
	assert.Equal(suite.T(), true, *result.Config.S3ForcePathStyle, "should forceS3PathStyle")
}

func (suite *S3BucketsTestSuite) TestGetS3BucketSessionNoLocalstack() {
	result, _ := getS3BucketSession(&S3BucketConfig{})

	assert.Nil(suite.T(), result.Config.Endpoint, "should be nil")
	assert.Equal(suite.T(), "us-east-1", *result.Config.Region, "should have correct region")
}

func (suite *S3BucketsTestSuite) TestInitializeS3BucketCreatedWhenLocalstackIsEnabled() {
	mockInitBucket := new(MockInitBucketHeadErrorS3API)
	bucketName := aws.String("logiflows-test-1")
	localStackAddr := aws.String("localstack")
	ctx := context.Background()
	MockInitBucket := makeInitializeS3Bucket(mockInitBucket)

	mockInitBucket.
		On("HeadBucketWithContext", ctx, &s3.HeadBucketInput{
			Bucket: bucketName,
		})

	mockInitBucket.
		On("CreateBucketWithContext", ctx, &s3.CreateBucketInput{
			Bucket: bucketName,
		})

	err := MockInitBucket(ctx, &S3BucketConfig{Name: bucketName, S3LocalstackAddress: localStackAddr})

	mockInitBucket.AssertCalled(suite.T(), "CreateBucketWithContext", ctx, &s3.CreateBucketInput{
		Bucket: bucketName,
	})

	assert.Nil(suite.T(), err, "Error should be nil")
}

func (suite *S3BucketsTestSuite) TestInitializeS3BucketNotCreatedWhenLocalstackDisabled() {
	mockInitBucket := new(MockInitBucketS3API)
	bucketName := aws.String("logiflows-test-1")
	ctx := context.Background()
	MockInitBucket := makeInitializeS3Bucket(mockInitBucket)

	mockInitBucket.
		On("HeadBucketWithContext", ctx, &s3.HeadBucketInput{
			Bucket: bucketName,
		})

	err := MockInitBucket(ctx, &S3BucketConfig{Name: bucketName})

	mockInitBucket.AssertNotCalled(suite.T(), "CreateBucketWithContext")

	assert.Nil(suite.T(), err, "Error should be nil")
}

func (suite *S3BucketsTestSuite) TestInitializeS3BucketErrorReturnedIfBucketDoesNotExist() {
	mockInitBucket := new(MockInitBucketHeadErrorS3API)
	bucketName := aws.String("logiflows-test-1")
	ctx := context.Background()
	MockInitBucket := makeInitializeS3Bucket(mockInitBucket)

	mockInitBucket.
		On("HeadBucketWithContext", ctx, &s3.HeadBucketInput{
			Bucket: bucketName,
		})

	err := MockInitBucket(ctx, &S3BucketConfig{Name: bucketName})

	mockInitBucket.AssertNotCalled(suite.T(), "CreateBucketWithContext")
	assert.Error(suite.T(), err, "Error should not be nil")
}

func (suite *S3BucketsTestSuite) TestInitializeS3BucketCreateFails() {
	mockInitBucket := new(MockInitBucketCreateErrorS3API)
	bucketName := aws.String("logiflows-test-1")
	localStackAddr := aws.String("localstack")
	ctx := context.Background()
	MockInitBucket := makeInitializeS3Bucket(mockInitBucket)

	mockInitBucket.
		On("HeadBucketWithContext", ctx, &s3.HeadBucketInput{
			Bucket: bucketName,
		})

	mockInitBucket.
		On("CreateBucketWithContext", ctx, &s3.CreateBucketInput{
			Bucket: bucketName,
		})

	err := MockInitBucket(ctx, &S3BucketConfig{Name: bucketName, S3LocalstackAddress: localStackAddr})

	assert.Error(suite.T(), err, "Error should be nil")
}

func (suite *S3BucketsTestSuite) TestInitializeS3BucketExists() {
	mockInitBucket := new(MockInitBucketS3API)
	bucketName := aws.String("logiflows-test-1")
	ctx := context.Background()
	MockInitBucket := makeInitializeS3Bucket(mockInitBucket)

	mockInitBucket.
		On("HeadBucketWithContext", ctx, &s3.HeadBucketInput{
			Bucket: bucketName,
		})

	err := MockInitBucket(ctx, &S3BucketConfig{Name: bucketName})

	mockInitBucket.AssertNotCalled(suite.T(), "CreateBucketWithContext")
	assert.Nil(suite.T(), err, "Error should be nil")
}
