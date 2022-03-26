package mocks

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockStatsDClient struct {
	mock.Mock
}

func (client *MockStatsDClient) Incr(metric string, tags []string, sampleRate float64) error {
	client.Called(metric, tags, sampleRate)
	return nil
}

func (client *MockStatsDClient) Histogram(metric string, value float64, tags []string, sampleRate float64) error {
	client.Called(metric, value, tags, sampleRate)
	return nil
}

func (client *MockStatsDClient) Gauge(metric string, value float64, tags []string, sampleRate float64) error {
	client.Called(metric, value, tags, sampleRate)
	return nil
}

func (client *MockStatsDClient) Timing(metric string, value time.Duration, tags []string, sampleRate float64) error {
	client.Called(metric, value, tags, sampleRate)
	return nil
}
