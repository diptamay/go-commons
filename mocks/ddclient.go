package mocks

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockDDClient struct {
	mock.Mock
}

func (d *MockDDClient) Increment(metric string, tags map[string]string) error {
	d.Called(metric, tags)
	return nil
}

func (d *MockDDClient) IncrementNoSampling(metric string, tags map[string]string) error {
	d.Called(metric, tags)
	return nil
}

func (d *MockDDClient) Histogram(metric string, value float64, tags map[string]string) error {
	d.Called(metric, value, tags)
	return nil
}

func (d *MockDDClient) HistogramNoSampling(metric string, value float64, tags map[string]string) error {
	d.Called(metric, value, tags)
	return nil
}

func (d *MockDDClient) Gauge(metric string, value float64, tags map[string]string) error {
	d.Called(metric, value, tags)
	return nil
}

func (d *MockDDClient) Timing(metric string, value time.Duration, tags map[string]string) error {
	d.Called(metric, value, tags)
	return nil
}

func (d *MockDDClient) TimingNoSampling(metric string, value time.Duration, tags map[string]string) error {
	d.Called(metric, value, tags)
	return nil
}

func (d *MockDDClient) Namespace() string {
	d.Called()
	return ""
}
