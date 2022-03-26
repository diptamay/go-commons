package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockLogger struct {
	mock.Mock
}

func (logger *MockLogger) Warn(message string, indexes map[string]interface{}) {
	logger.Called(message, indexes)
}

func (logger *MockLogger) Error(message string, indexes map[string]interface{}) {
	logger.Called(message, indexes)
}

func (logger *MockLogger) Info(message string, indexes map[string]interface{}) {
	logger.Called(message, indexes)
}

func (logger *MockLogger) Debug(message string, indexes map[string]interface{}) {
	logger.Called(message, indexes)
}

func (logger *MockLogger) Trace(message string, indexes map[string]interface{}) {
	logger.Called(message, indexes)
}

func (logger *MockLogger) Fatal(message string, indexes map[string]interface{}) {
	logger.Called(message, indexes)
}

func (logger *MockLogger) Config() zap.Config {
	logger.Called()
	return zap.Config{}
}
