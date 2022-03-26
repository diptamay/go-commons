package metrics

import (
	"fmt"
	"time"

	"github.com/diptamay/go-commons/doggie"
	"github.com/diptamay/go-commons/glogger"
)

var (
	emptyDDTagsMap  = map[string]string{}
	emptyLogTagsMap = map[string]interface{}{}
)

func HandleNils(tags map[string]string, loggerInfo map[string]interface{}) (map[string]string, map[string]interface{}) {
	var retTags = tags
	if retTags == nil {
		retTags = emptyDDTagsMap
	}
	var retLoggerInfo = loggerInfo
	if retLoggerInfo == nil {
		retLoggerInfo = emptyLogTagsMap
	}
	return retTags, retLoggerInfo
}

type Metrics interface {
	Increment(metric string, tags map[string]string, loggerInfo map[string]interface{})
	IncrNoLogTags(metric string, tags map[string]string)
	IncrNoDDTags(metric string, loggerInfo map[string]interface{})
	IncrNoTags(metric string)

	IncrementNoSampling(metric string, tags map[string]string, loggerInfo map[string]interface{})
	IncrNoSamplingOrLogTags(metric string, tags map[string]string)
	IncrNoSamplingOrDDTags(metric string, loggerInfo map[string]interface{})
	IncrNoSamplingOrTags(metric string)

	Histogram(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{})
	HistogramNoLogTags(metric string, value float64, tags map[string]string)
	HistogramNoDDTags(metric string, value float64, loggerInfo map[string]interface{})
	HistogramNoTags(metric string, value float64)

	HistogramNoSampling(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{})
	HistogramNoSamplingOrLogTags(metric string, value float64, tags map[string]string)
	HistogramNoSamplingOrDDTags(metric string, value float64, loggerInfo map[string]interface{})
	HistogramNoSamplingOrTags(metric string, value float64)

	Gauge(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{})
	GaugeNoLogTags(metric string, value float64, tags map[string]string)
	GaugeNoDDTags(metric string, value float64, loggerInfo map[string]interface{})
	GaugeNoTags(metric string, value float64)

	Timing(metric string, value time.Duration, tags map[string]string, loggerInfo map[string]interface{})
	TimingNoLogTags(metric string, value time.Duration, tags map[string]string)
	TimingNoDDTags(metric string, value time.Duration, loggerInfo map[string]interface{})
	TimingNoTags(metric string, value time.Duration)

	TimingNoSampling(metric string, value time.Duration, tags map[string]string, loggerInfo map[string]interface{})
	TimingNoSamplingOrLogTags(metric string, value time.Duration, tags map[string]string)
	TimingNoSamplingOrDDTags(metric string, value time.Duration, loggerInfo map[string]interface{})
	TimingNoSamplingOrTags(metric string, value time.Duration)
}

type metricsImpl struct {
	client doggie.DataDogClient
	logger glogger.Logger
}

func NewMetrics(d doggie.DataDogClient, gl glogger.Logger) Metrics {
	return &metricsImpl{d, gl}
}

func NewEmptyMetrics() Metrics {
	return &metricsImpl{}
}

func (m *metricsImpl) callIncr(metric string, tags map[string]string, loggerInfo map[string]interface{}, fnString string) {
	if m.client != nil {
		var action func(string, map[string]string) error
		if fnString == "Increment" {
			action = m.client.Increment
		} else {
			action = m.client.IncrementNoSampling
		}
		tags, loggerInfo = HandleNils(tags, loggerInfo)
		if err := action(metric, tags); err != nil {
			if m.logger != nil {
				m.logger.Warn(fmt.Sprintf("datadog reporting failed with error: %s", err.Error()), loggerInfo)
			}
		}
	}
}

func (m *metricsImpl) callTiming(metric string, value time.Duration, tags map[string]string, loggerInfo map[string]interface{}, fnString string) {
	if m.client != nil {
		var action func(string, time.Duration, map[string]string) error
		if fnString == "Timing" {
			action = m.client.Timing
		} else {
			action = m.client.TimingNoSampling
		}
		tags, loggerInfo = HandleNils(tags, loggerInfo)
		if err := action(metric, value, tags); err != nil {
			if m.logger != nil {
				m.logger.Warn(fmt.Sprintf("datadog reporting failed with error: %s", err.Error()), loggerInfo)
			}
		}
	}
}

func (m *metricsImpl) callOthers(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{}, fnString string) {
	if m.client != nil {
		var action func(string, float64, map[string]string) error
		if fnString == "Gauge" {
			action = m.client.Gauge
		} else if fnString == "Histogram" {
			action = m.client.Histogram
		} else {
			action = m.client.HistogramNoSampling
		}
		tags, loggerInfo = HandleNils(tags, loggerInfo)
		if err := action(metric, value, tags); err != nil {
			if m.logger != nil {
				m.logger.Warn(fmt.Sprintf("datadog reporting failed with error: %s", err.Error()), loggerInfo)
			}
		}
	}
}

func (m *metricsImpl) Increment(metric string, tags map[string]string, loggerInfo map[string]interface{}) {
	m.callIncr(metric, tags, loggerInfo, "Increment")
}

func (m *metricsImpl) IncrNoLogTags(metric string, tags map[string]string) {
	m.Increment(metric, tags, nil)
}

func (m *metricsImpl) IncrNoDDTags(metric string, loggerInfo map[string]interface{}) {
	m.Increment(metric, nil, loggerInfo)
}

func (m *metricsImpl) IncrNoTags(metric string) {
	m.Increment(metric, nil, nil)
}

func (m *metricsImpl) IncrementNoSampling(metric string, tags map[string]string, loggerInfo map[string]interface{}) {
	m.callIncr(metric, tags, loggerInfo, "IncrementNoSampling")
}

func (m *metricsImpl) IncrNoSamplingOrLogTags(metric string, tags map[string]string) {
	m.IncrementNoSampling(metric, tags, nil)
}

func (m *metricsImpl) IncrNoSamplingOrDDTags(metric string, loggerInfo map[string]interface{}) {
	m.IncrementNoSampling(metric, nil, loggerInfo)
}

func (m *metricsImpl) IncrNoSamplingOrTags(metric string) {
	m.IncrementNoSampling(metric, nil, nil)
}

func (m *metricsImpl) Histogram(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{}) {
	m.callOthers(metric, value, tags, loggerInfo, "Histogram")
}

func (m *metricsImpl) HistogramNoLogTags(metric string, value float64, tags map[string]string) {
	m.Histogram(metric, value, tags, nil)
}

func (m *metricsImpl) HistogramNoDDTags(metric string, value float64, loggerInfo map[string]interface{}) {
	m.Histogram(metric, value, nil, loggerInfo)
}

func (m *metricsImpl) HistogramNoTags(metric string, value float64) {
	m.Histogram(metric, value, nil, nil)
}

func (m *metricsImpl) HistogramNoSampling(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{}) {
	m.callOthers(metric, value, tags, loggerInfo, "HistogramNoSampling")
}

func (m *metricsImpl) HistogramNoSamplingOrLogTags(metric string, value float64, tags map[string]string) {
	m.HistogramNoSampling(metric, value, tags, nil)
}

func (m *metricsImpl) HistogramNoSamplingOrDDTags(metric string, value float64, loggerInfo map[string]interface{}) {
	m.HistogramNoSampling(metric, value, nil, loggerInfo)
}

func (m *metricsImpl) HistogramNoSamplingOrTags(metric string, value float64) {
	m.HistogramNoSampling(metric, value, nil, nil)
}

func (m *metricsImpl) Gauge(metric string, value float64, tags map[string]string, loggerInfo map[string]interface{}) {
	m.callOthers(metric, value, tags, loggerInfo, "Gauge")
}

func (m *metricsImpl) GaugeNoLogTags(metric string, value float64, tags map[string]string) {
	m.Gauge(metric, value, tags, nil)
}

func (m *metricsImpl) GaugeNoDDTags(metric string, value float64, loggerInfo map[string]interface{}) {
	m.Gauge(metric, value, nil, loggerInfo)
}

func (m *metricsImpl) GaugeNoTags(metric string, value float64) {
	m.Gauge(metric, value, nil, nil)
}

func (m *metricsImpl) Timing(metric string, value time.Duration, tags map[string]string, loggerInfo map[string]interface{}) {
	m.callTiming(metric, value, tags, loggerInfo, "Timing")
}

func (m *metricsImpl) TimingNoLogTags(metric string, value time.Duration, tags map[string]string) {
	m.Timing(metric, value, tags, nil)
}

func (m *metricsImpl) TimingNoDDTags(metric string, value time.Duration, loggerInfo map[string]interface{}) {
	m.Timing(metric, value, nil, loggerInfo)
}

func (m *metricsImpl) TimingNoTags(metric string, value time.Duration) {
	m.Timing(metric, value, nil, nil)
}

func (m *metricsImpl) TimingNoSampling(metric string, value time.Duration, tags map[string]string, loggerInfo map[string]interface{}) {
	m.callTiming(metric, value, tags, loggerInfo, "TimingNoSampling")
}

func (m *metricsImpl) TimingNoSamplingOrLogTags(metric string, value time.Duration, tags map[string]string) {
	m.TimingNoSampling(metric, value, tags, nil)
}

func (m *metricsImpl) TimingNoSamplingOrDDTags(metric string, value time.Duration, loggerInfo map[string]interface{}) {
	m.TimingNoSampling(metric, value, nil, loggerInfo)
}

func (m *metricsImpl) TimingNoSamplingOrTags(metric string, value time.Duration) {
	m.TimingNoSampling(metric, value, nil, nil)
}
