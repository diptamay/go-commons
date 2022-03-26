package mocks

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockMetrics struct {
	mock.Mock
}

func (m *MockMetrics) Increment(metric string, tags map[string]string, loggerInfo map[string]interface{}) {
	m.Called(metric, tags, loggerInfo)
}
func (m *MockMetrics) IncrNoLogTags(metric string, tags map[string]string) {
	m.Called(metric, tags)
}
func (m *MockMetrics) IncrNoDDTags(metric string, loggerInfo map[string]interface{}) {
	m.Called(metric, loggerInfo)
}
func (m *MockMetrics) IncrNoTags(metric string) {
	m.Called(metric)
}

func (m *MockMetrics) IncrementNoSampling(metric string, tags map[string]string, loggerInfo map[string]interface{}) {
	m.Called(metric, tags, loggerInfo)
}
func (m *MockMetrics) IncrNoSamplingOrLogTags(metric string, tags map[string]string) {
	m.Called(metric, tags)
}
func (m *MockMetrics) IncrNoSamplingOrDDTags(metric string, loggerInfo map[string]interface{}) {
	m.Called(metric, loggerInfo)
}
func (m *MockMetrics) IncrNoSamplingOrTags(metric string) {
	m.Called(metric)
}

func (m *MockMetrics) Histogram(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{}) {
	m.Called(metric, value, tags, loggerInfo)
}
func (m *MockMetrics) HistogramNoLogTags(metric string, value float64, tags map[string]string) {
	m.Called(metric, value, tags)
}
func (m *MockMetrics) HistogramNoDDTags(metric string, value float64, loggerInfo map[string]interface{}) {
	m.Called(metric, value, loggerInfo)
}
func (m *MockMetrics) HistogramNoTags(metric string, value float64) {
	m.Called(metric, value)
}

func (m *MockMetrics) HistogramNoSampling(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{}) {
	m.Called(metric, value, tags, loggerInfo)
}
func (m *MockMetrics) HistogramNoSamplingOrLogTags(metric string, value float64, tags map[string]string) {
	m.Called(metric, value, tags)
}
func (m *MockMetrics) HistogramNoSamplingOrDDTags(metric string, value float64, loggerInfo map[string]interface{}) {
	m.Called(metric, value, loggerInfo)
}
func (m *MockMetrics) HistogramNoSamplingOrTags(metric string, value float64) {
	m.Called(metric, value)
}

func (m *MockMetrics) Gauge(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{}) {
	m.Called(metric, value, tags, loggerInfo)
}
func (m *MockMetrics) GaugeNoLogTags(metric string, value float64, tags map[string]string) {
	m.Called(metric, value, tags)
}
func (m *MockMetrics) GaugeNoDDTags(metric string, value float64, loggerInfo map[string]interface{}) {
	m.Called(metric, value, loggerInfo)
}
func (m *MockMetrics) GaugeNoTags(metric string, value float64) {
	m.Called(metric, value)
}

func (m *MockMetrics) Timing(metric string, value time.Duration, tags map[string]string, loggerInfo map[string]interface{}) {
	m.Called(metric, value, tags, loggerInfo)
}
func (m *MockMetrics) TimingNoLogTags(metric string, value time.Duration, tags map[string]string) {
	m.Called(metric, value, tags)
}
func (m *MockMetrics) TimingNoDDTags(metric string, value time.Duration, loggerInfo map[string]interface{}) {
	m.Called(metric, value, loggerInfo)
}
func (m *MockMetrics) TimingNoTags(metric string, value time.Duration) {
	m.Called(metric, value)
}

func (m *MockMetrics) TimingNoSampling(metric string, value time.Duration, tags map[string]string, loggerInfo map[string]interface{}) {
	m.Called(metric, value, tags, loggerInfo)
}
func (m *MockMetrics) TimingNoSamplingOrLogTags(metric string, value time.Duration, tags map[string]string) {
	m.Called(metric, value, tags)
}
func (m *MockMetrics) TimingNoSamplingOrDDTags(metric string, value time.Duration, loggerInfo map[string]interface{}) {
	m.Called(metric, value, loggerInfo)
}
func (m *MockMetrics) TimingNoSamplingOrTags(metric string, value time.Duration) {
	m.Called(metric, value)
}
