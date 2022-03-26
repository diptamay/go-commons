package glogger

import (
	"github.com/fatih/structs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"reflect"
	"time"
)

var (
	pid         = os.Getpid()
	hostname, _ = os.Hostname()
	loggerv     = "logilogger@0.0.0-alpha.0"
)
var logLevelSeverity = map[string]int{
	"fatal": 1,
	"error": 3,
	"warn":  4,
	"info":  6,
	"debug": 7,
	"trace": 7,
}
var debugLevel = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"fatal": zap.FatalLevel,
}

type internalLogger interface {
	Error(string, ...zap.Field)
	Warn(string, ...zap.Field)
	Info(string, ...zap.Field)
	Debug(string, ...zap.Field)
	Fatal(string, ...zap.Field)
	Sync() error
}

type Logger interface {
	Error(message string, indexes map[string]interface{})
	Warn(message string, indexes map[string]interface{})
	Info(message string, indexes map[string]interface{})
	Debug(message string, indexes map[string]interface{})
	Trace(message string, indexes map[string]interface{})
	Fatal(message string, indexes map[string]interface{})
	Config() zap.Config
}

type loggerImpl struct {
	internal internalLogger
	cfg      zap.Config
	defaults interface{}
}

func getZapField(key string, value interface{}) (zap.Field, error) {
	var zapField zap.Field
	switch v := value.(type) {
	case int:
		zapField = zap.Int(key, v)
		break
	case int64:
		zapField = zap.Int64(key, v)
		break
	case float32:
		zapField = zap.Float32(key, v)
		break
	case float64:
		zapField = zap.Float64(key, v)
		break
	case string:
		zapField = zap.String(key, v)
		break
	case time.Time:
		zapField = zap.Time(key, v)
		break
	}
	return zapField, nil
}

func getStructFieldNames(fields []*structs.Field) []string {
	keys := []string{}

	for _, field := range fields {
		if field.IsExported() {
			nestedFields := getStructFieldNames(field.Fields())
			keys = append(keys, nestedFields...)
		} else {
			keys = append(keys, field.Name())
		}
	}

	return keys
}

func makeZapFields(logger *loggerImpl, indexes map[string]interface{}) []zap.Field {
	fields := []zap.Field{}
	loggerDefaults := structs.New(logger.defaults)
	keys := getStructFieldNames(loggerDefaults.Fields())

	for i := 0; i < len(keys); i++ {
		if value, ok := indexes[keys[i]]; ok {
			if ok = fieldMatchesSchema(keys[i], reflect.TypeOf(value), logger.defaults); ok {
				field, _ := getZapField(keys[i], value)
				fields = append(fields, field)
			}
		}
	}

	return fields
}

func handleFields(logger *loggerImpl, indexes map[string]interface{}, level string) ([]zap.Field, error) {
	fields := makeZapFields(logger, indexes)
	fields = append(fields, zap.Int("severity", logLevelSeverity[level]))
	return fields, nil
}

func (logger *loggerImpl) Error(message string, indexes map[string]interface{}) {
	fields, _ := handleFields(logger, indexes, "error")
	logger.internal.Error(message, fields...)
}

func (logger *loggerImpl) Warn(message string, indexes map[string]interface{}) {
	fields, _ := handleFields(logger, indexes, "warn")
	logger.internal.Warn(message, fields...)
}

func (logger *loggerImpl) Info(message string, indexes map[string]interface{}) {
	fields, _ := handleFields(logger, indexes, "info")
	logger.internal.Info(message, fields...)
}

func (logger *loggerImpl) Debug(message string, indexes map[string]interface{}) {
	fields, _ := handleFields(logger, indexes, "debug")
	logger.internal.Debug(message, fields...)
}

func (logger *loggerImpl) Trace(message string, indexes map[string]interface{}) {
	fields, _ := handleFields(logger, indexes, "trace")
	logger.internal.Debug(message, fields...)
}

func (logger *loggerImpl) Fatal(message string, indexes map[string]interface{}) {
	fields, _ := handleFields(logger, indexes, "fatal")
	logger.internal.Fatal(message, fields...)
}

func (logger *loggerImpl) Sync() {
	logger.internal.Sync()
}

func (logger *loggerImpl) Config() zap.Config {
	return logger.cfg
}

func CreateLogger(constants map[string]interface{}, schema interface{}, options map[string]interface{}) Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var level zapcore.Level
	ok := false
	if reflect.TypeOf(options["logLevel"]) == reflect.TypeOf("") {
		level, ok = debugLevel[reflect.ValueOf(options["logLevel"]).String()]
	}
	if !ok {
		level = debugLevel["info"]
	}
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    constants,
		EncoderConfig:    encoderConfig,
	}
	if cfg.InitialFields["v"] == nil {
		cfg.InitialFields["v"] = loggerv
	}
	if cfg.InitialFields["pid"] == nil {
		cfg.InitialFields["pid"] = pid
	}
	if cfg.InitialFields["containerId"] == nil {
		cfg.InitialFields["containerId"] = hostname
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return &loggerImpl{
		internal: logger,
		cfg:      cfg,
		defaults: schema,
	}
}
