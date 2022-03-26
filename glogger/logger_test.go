package glogger

import (
	"fmt"
	Chance "github.com/ZeFort/chance"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
	"time"
)

const IntExample int = 32767
const Int64Example int64 = 9223372036854775807
const Float32Example float32 = 32767.01
const Float64Example float64 = 2147483647.01
const StringExample string = "hello"

var chance = Chance.New()
var (
	ErrorMessage = chance.String()
	WarnMessage  = chance.String()
	InfoMessage  = chance.String()
	DebugMessage = chance.String()
	FatalMessage = chance.String()
)

var TimeExample = time.Now()
var levelKeys = []string{
	"fatal",
	"error",
	"warn",
	"info",
	"debug",
	"trace",
}

type mockZapLogger struct {
	mock.Mock
}

func (logger *mockZapLogger) Error(message string, fields ...zap.Field) {
	logger.Called(message, fields[0], fields[1])
}
func (logger *mockZapLogger) Warn(message string, fields ...zap.Field) {
	logger.Called(message, fields[0], fields[1])
}
func (logger *mockZapLogger) Info(message string, fields ...zap.Field) {
	logger.Called(message, fields[0], fields[1])
}
func (logger *mockZapLogger) Debug(message string, fields ...zap.Field) {
	logger.Called(message, fields[0], fields[1])
}
func (logger *mockZapLogger) Fatal(message string, fields ...zap.Field) {
	logger.Called(message, fields[0], fields[1])
}
func (logger *mockZapLogger) Sync() error {
	logger.Called()
	return nil
}

func TestGetZapField(t *testing.T) {
	testkey := "test"
	tests := []struct {
		expected zap.Field
		actual   interface{}
		result   bool
		message  string
	}{
		{
			expected: zap.Int(testkey, IntExample),
			actual:   IntExample,
			result:   true,
			message:  fmt.Sprintf("value of type %s should return %s derived zap.Field", "int", "zap.Int"),
		},
		{
			expected: zap.Int64(testkey, Int64Example),
			actual:   Int64Example,
			result:   true,
			message:  fmt.Sprintf("value of type %s should return %s derived zap.Field", "int64", "zap.Int64"),
		},
		{
			expected: zap.Float32(testkey, Float32Example),
			actual:   Float32Example,
			result:   true,
			message:  fmt.Sprintf("value of type %s should return %s derived zap.Field", "float32", "zap.Float32"),
		},
		{
			expected: zap.Float64(testkey, Float64Example),
			actual:   Float64Example,
			result:   true,
			message:  fmt.Sprintf("value of type %s should return %s derived zap.Field", "float64", "zap.Float64"),
		},
		{
			expected: zap.String(testkey, StringExample),
			actual:   StringExample,
			result:   true,
			message:  fmt.Sprintf("value of type %s should return %s derived zap.Field", "string", "zap.String"),
		},
		{
			expected: zap.Time(testkey, TimeExample),
			actual:   TimeExample,
			result:   true,
			message:  fmt.Sprintf("value of type %s should return %s derived zap.Field", "time.Time", "zap.Time"),
		},
	}
	for _, test := range tests {
		field, _ := getZapField(testkey, test.actual)
		assert.Equal(t, test.result, test.expected == field, test.message)
	}
}

func generateZapFieldsComparison(key string, fields []zap.Field, shouldExist bool) assert.Comparison {
	return assert.Comparison(func() bool {
		var exists bool
		for _, field := range fields {
			if key == field.Key {
				exists = true
				break
			}
		}
		return exists == shouldExist
	})
}

func TestMakeZapFields(t *testing.T) {
	type NestedTestStruct struct {
		traceId string
		spanId  string
	}
	type TestStruct struct {
		NestedTestStruct
		id        int
		name      string
		createdat time.Time
		updatedat time.Time
	}
	mockLogger := &loggerImpl{
		defaults: TestStruct{},
	}
	values := &map[string]interface{}{
		"id":        chance.Int(),
		"name":      chance.String(),
		"createdat": time.Now(),
		"updatedat": chance.String(),
		"other":     chance.String(),
		"traceId":   chance.String(),
		"spanId":    chance.String(),
	}
	fields := makeZapFields(mockLogger, *values)
	expectedFields := []string{"id", "name", "createdat", "traceId", "spanId"}
	unexpectedFields := []string{"updatedat", "other"}
	for _, key := range expectedFields {
		assert.Condition(t, generateZapFieldsComparison(key, fields, true), fmt.Sprintf("\"%s\" should appear in fields", key))
	}
	for _, key := range unexpectedFields {
		assert.Condition(t, generateZapFieldsComparison(key, fields, false), fmt.Sprintf("\"%s\" should not appear in fields", key))
	}
}

func TestHandleFields(t *testing.T) {
	type TestStruct struct {
		id int
	}
	mockLogger := &loggerImpl{
		defaults: TestStruct{},
	}
	level := levelKeys[chance.IntBtw(0, len(levelKeys))]
	fields, _ := handleFields(mockLogger, map[string]interface{}{
		"id": chance.Int(),
	}, level)
	expectedFields := []string{"id", "severity"}
	for _, key := range expectedFields {
		assert.Condition(t, generateZapFieldsComparison(key, fields, true), fmt.Sprintf("\"%s\" should appear in fields", key))
	}
}

func TestGetStructFieldNames(t *testing.T) {
	type DeeplyNestedTestStruct struct {
		deeplyNestedId int
	}
	type NestedTestStruct struct {
		DeeplyNestedTestStruct
		nestedId int
	}
	type TestStruct struct {
		NestedTestStruct
		id int
	}
	mockLogger := &loggerImpl{
		defaults: TestStruct{},
	}
	mockLoggerDefaults := structs.New(mockLogger.defaults)
	keys := getStructFieldNames(mockLoggerDefaults.Fields())
	expectedFields := []string{"deeplyNestedId", "nestedId", "id"}
	assert.Equal(
		t,
		expectedFields,
		keys,
		"should generate a slice of field names of a struct, including any level of nested structs",
	)
}

func TestError(t *testing.T) {
	type TestStruct struct {
		id int
	}
	mockErrorZapLogger := &mockZapLogger{}
	mockErrorZapLogger.
		On("Error", ErrorMessage, zap.Int("id", 1), zap.Int("severity", logLevelSeverity["error"])).
		Return()
	mockLogger := &loggerImpl{
		internal: mockErrorZapLogger,
		defaults: TestStruct{},
	}
	mockLogger.Error(ErrorMessage, map[string]interface{}{
		"id": 1,
	})
	mockErrorZapLogger.AssertExpectations(t)
}

func TestWarn(t *testing.T) {
	type TestStruct struct {
		id int
	}
	mockWarnZapLogger := &mockZapLogger{}
	mockWarnZapLogger.
		On("Warn", WarnMessage, zap.Int("id", 1), zap.Int("severity", logLevelSeverity["warn"])).
		Return()
	mockLogger := &loggerImpl{
		internal: mockWarnZapLogger,
		defaults: TestStruct{},
	}
	mockLogger.Warn(WarnMessage, map[string]interface{}{
		"id": 1,
	})
	mockWarnZapLogger.AssertExpectations(t)
}

func TestInfo(t *testing.T) {
	type TestStruct struct {
		id int
	}
	mockInfoZapLogger := &mockZapLogger{}
	mockInfoZapLogger.
		On("Info", InfoMessage, zap.Int("id", 1), zap.Int("severity", logLevelSeverity["info"])).
		Return()
	mockLogger := &loggerImpl{
		internal: mockInfoZapLogger,
		defaults: TestStruct{},
	}
	mockLogger.Info(InfoMessage, map[string]interface{}{
		"id": 1,
	})
	mockInfoZapLogger.AssertExpectations(t)
}

func TestDebug(t *testing.T) {
	type TestStruct struct {
		id int
	}
	mockDebugZapLogger := &mockZapLogger{}
	mockDebugZapLogger.
		On("Debug", DebugMessage, zap.Int("id", 1), zap.Int("severity", logLevelSeverity["debug"])).
		Return()
	mockLogger := &loggerImpl{
		internal: mockDebugZapLogger,
		defaults: TestStruct{},
	}
	mockLogger.Debug(DebugMessage, map[string]interface{}{
		"id": 1,
	})
	mockDebugZapLogger.AssertExpectations(t)
}

func TestFatal(t *testing.T) {
	type TestStruct struct {
		id int
	}
	mockFatalZapLogger := &mockZapLogger{}
	mockFatalZapLogger.
		On("Fatal", FatalMessage, zap.Int("id", 1), zap.Int("severity", logLevelSeverity["fatal"])).
		Return()
	mockLogger := &loggerImpl{
		internal: mockFatalZapLogger,
		defaults: TestStruct{},
	}
	mockLogger.Fatal(FatalMessage, map[string]interface{}{
		"id": 1,
	})
	mockFatalZapLogger.AssertExpectations(t)
}

func TestTrace(t *testing.T) {
	type TestStruct struct {
		id int
	}
	mockTraceZapLogger := &mockZapLogger{}
	mockTraceZapLogger.
		On("Debug", DebugMessage, zap.Int("id", 1), zap.Int("severity", logLevelSeverity["trace"])).
		Return()
	mockLogger := &loggerImpl{
		internal: mockTraceZapLogger,
		defaults: TestStruct{},
	}
	mockLogger.Trace(DebugMessage, map[string]interface{}{
		"id": 1,
	})
	mockTraceZapLogger.AssertExpectations(t)
}

func TestSync(t *testing.T) {
	type TestStruct struct {
		id int
	}
	mockSyncZapLogger := &mockZapLogger{}
	mockSyncZapLogger.
		On("Sync").
		Return()
	mockLogger := &loggerImpl{
		internal: mockSyncZapLogger,
		defaults: TestStruct{},
	}
	mockLogger.Sync()
	mockSyncZapLogger.AssertExpectations(t)
}

func TestCreateLogger(t *testing.T) {
	constants := map[string]interface{}{
		"id":   1,
		"name": chance.Word(),
	}
	type Schema struct {
		id   string
		name string
	}
	schema := &Schema{}
	options := map[string]interface{}{
		"logLevel": "info",
	}
	logger := CreateLogger(constants, schema, options)
	values := constants
	assert.Equal(t, map[string]interface{}{
		"id":          1,
		"name":        values["name"],
		"pid":         pid,
		"containerId": hostname,
		"v":           loggerv,
	}, logger.Config().InitialFields, "should create logger with initial fields defined as defaults and provided constant values")
}

func TestCreateLoggerDefaultLogLevel(t *testing.T) {
	constants := map[string]interface{}{
		"id":   1,
		"name": chance.Word(),
	}
	type Schema struct {
		id   string
		name string
	}
	schema := &Schema{}
	options := map[string]interface{}{}
	logger := CreateLogger(constants, schema, options)
	assert.Equal(
		t,
		zap.NewAtomicLevelAt(zap.InfoLevel),
		logger.Config().Level,
		"should use default debug level if not defined in options",
	)
}
