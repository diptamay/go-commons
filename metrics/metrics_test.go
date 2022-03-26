package metrics

import (
	"testing"
	"time"

	Chance "github.com/ZeFort/chance"
	"github.com/diptamay/go-commons/mocks"
	"github.com/stretchr/testify/assert"
)

func setup(nilTags bool, nilLogs bool) (*mocks.MockDDClient, *mocks.MockLogger, *Chance.Chance, string, map[string]string, map[string]interface{}) {
	client := &mocks.MockDDClient{}
	logger := &mocks.MockLogger{}

	chance := Chance.New()
	metric := chance.String()
	var tagsMap map[string]string
	if !nilTags {
		tagsMap = tagsData(chance)
	} else {
		tagsMap = emptyDDTagsMap
	}
	var logIndexes map[string]interface{}
	if !nilLogs {
		logIndexes = logIndexData(chance)
	} else {
		logIndexes = emptyLogTagsMap
	}
	return client, logger, chance, metric, tagsMap, logIndexes
}

func tagsData(chance *Chance.Chance) map[string]string {
	return map[string]string{
		"a": chance.Word(),
		"b": chance.Word(),
	}
}

func logIndexData(chance *Chance.Chance) map[string]interface{} {
	return map[string]interface{}{
		"m": chance.Word(),
		"n": chance.Word(),
	}
}

func setupForCounters(statsFuncName string, logFuncName string, nilTags bool, nilLogs bool) (*mocks.MockDDClient, *mocks.MockLogger, string, map[string]string, map[string]interface{}) {
	client, logger, _, metric, tagsMap, logIndexes := setup(nilTags, nilLogs)

	client.On(statsFuncName, metric, tagsMap).Return(nil)
	logger.On(logFuncName, logIndexes).Return(nil)
	return client, logger, metric, tagsMap, logIndexes
}

func setupForValues(statsFuncName string, logFuncName string, nilTags bool, nilLogs bool) (*mocks.MockDDClient, *mocks.MockLogger, string, float64, map[string]string, map[string]interface{}) {
	client, logger, chance, metric, tagsMap, logIndexes := setup(nilTags, nilLogs)

	value := chance.Float()
	client.On(statsFuncName, metric, value, tagsMap).Return(nil)
	logger.On(logFuncName, logIndexes).Return(nil)
	return client, logger, metric, value, tagsMap, logIndexes
}

func setupForTimings(statsFuncName string, logFuncName string, nilTags bool, nilLogs bool) (*mocks.MockDDClient, *mocks.MockLogger, string, time.Duration, map[string]string, map[string]interface{}) {
	client, logger, chance, metric, tagsMap, logIndexes := setup(nilTags, nilLogs)

	value := time.Duration(chance.Int())
	client.On(statsFuncName, metric, value, tagsMap).Return(nil)
	logger.On(logFuncName, logIndexes).Return(nil)
	return client, logger, metric, value, tagsMap, logIndexes
}

func TestHandleNils(t *testing.T) {
	retTags, retLoggerInfo := HandleNils(nil, nil)
	assert.NotNil(t, retTags)
	assert.NotNil(t, retLoggerInfo)
}

func TestHandleNilLogs(t *testing.T) {
	tagsMap := tagsData(Chance.New())
	retTags, retLoggerInfo := HandleNils(tagsMap, nil)
	assert.NotNil(t, retTags)
	assert.NotNil(t, retLoggerInfo)
}

func TestHandleNilTags(t *testing.T) {
	logIndexes := logIndexData(Chance.New())
	retTags, retLoggerInfo := HandleNils(nil, logIndexes)
	assert.NotNil(t, retTags)
	assert.NotNil(t, retLoggerInfo)
}

func TestNewEmptyMetrics(t *testing.T) {
	m := NewEmptyMetrics()
	assert.NotNil(t, m)
}

func TestIncrement(t *testing.T) {
	client, logger, metric, tagsMap, logIndexes := setupForCounters("Increment", "Warn", false, false)
	m := NewMetrics(client, logger)

	m.Increment(metric, tagsMap, logIndexes)
	client.AssertExpectations(t)
}

func TestIncrNoLogTags(t *testing.T) {
	client, logger, metric, tagsMap, _ := setupForCounters("Increment", "Warn", false, true)
	m := NewMetrics(client, logger)

	m.IncrNoLogTags(metric, tagsMap)
	client.AssertExpectations(t)
}

func TestIncrNoDDTags(t *testing.T) {
	client, logger, metric, _, logIndexes := setupForCounters("Increment", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.IncrNoDDTags(metric, logIndexes)
	client.AssertExpectations(t)
}

func TestIncrNoTags(t *testing.T) {
	client, logger, metric, _, _ := setupForCounters("Increment", "Warn", true, true)
	m := NewMetrics(client, logger)

	m.IncrNoTags(metric)
	client.AssertExpectations(t)
}

func TestIncrementNoSampling(t *testing.T) {
	client, logger, metric, tagsMap, logIndexes := setupForCounters("IncrementNoSampling", "Warn", false, false)
	m := NewMetrics(client, logger)

	m.IncrementNoSampling(metric, tagsMap, logIndexes)
	client.AssertExpectations(t)
}

func TestIncrNoSamplingOrLogTags(t *testing.T) {
	client, logger, metric, tagsMap, _ := setupForCounters("IncrementNoSampling", "Warn", false, true)
	m := NewMetrics(client, logger)

	m.IncrNoSamplingOrLogTags(metric, tagsMap)
	client.AssertExpectations(t)
}

func TestIncrNoSamplingOrDDTags(t *testing.T) {
	client, logger, metric, _, logIndexes := setupForCounters("IncrementNoSampling", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.IncrNoSamplingOrDDTags(metric, logIndexes)
	client.AssertExpectations(t)
}

func TestIncrNoSamplingOrTags(t *testing.T) {
	client, logger, metric, _, _ := setupForCounters("IncrementNoSampling", "Warn", true, true)
	m := NewMetrics(client, logger)

	m.IncrNoSamplingOrTags(metric)
	client.AssertExpectations(t)
}

func TestHistogram(t *testing.T) {
	client, logger, metric, value, tagsMap, logIndexes := setupForValues("Histogram", "Warn", false, false)
	m := NewMetrics(client, logger)

	m.Histogram(metric, value, tagsMap, logIndexes)
	client.AssertExpectations(t)
}

func TestHistogramNoLogTags(t *testing.T) {
	client, logger, metric, value, tagsMap, _ := setupForValues("Histogram", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.HistogramNoLogTags(metric, value, tagsMap)
	client.AssertExpectations(t)
}

func TestHistogramNoDDTags(t *testing.T) {
	client, logger, metric, value, _, logIndexes := setupForValues("Histogram", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.HistogramNoDDTags(metric, value, logIndexes)
	client.AssertExpectations(t)
}

func TestHistogramNoTags(t *testing.T) {
	client, logger, metric, value, _, _ := setupForValues("Histogram", "Warn", true, true)
	m := NewMetrics(client, logger)

	m.HistogramNoTags(metric, value)
	client.AssertExpectations(t)
}

func TestHistogramNoSampling(t *testing.T) {
	client, logger, metric, value, tagsMap, logIndexes := setupForValues("HistogramNoSampling", "Warn", false, false)
	m := NewMetrics(client, logger)

	m.HistogramNoSampling(metric, value, tagsMap, logIndexes)
	client.AssertExpectations(t)
}

func TestHistogramNoSamplingOrLogTags(t *testing.T) {
	client, logger, metric, value, tagsMap, _ := setupForValues("HistogramNoSampling", "Warn", false, true)
	m := NewMetrics(client, logger)

	m.HistogramNoSamplingOrLogTags(metric, value, tagsMap)
	client.AssertExpectations(t)
}

func TestHistogramNoSamplingOrDDTags(t *testing.T) {
	client, logger, metric, value, _, logIndexes := setupForValues("HistogramNoSampling", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.HistogramNoSamplingOrDDTags(metric, value, logIndexes)
	client.AssertExpectations(t)
}

func TestHistogramNoSamplingOrTags(t *testing.T) {
	client, logger, metric, value, _, _ := setupForValues("HistogramNoSampling", "Warn", true, true)
	m := NewMetrics(client, logger)

	m.HistogramNoSamplingOrTags(metric, value)
	client.AssertExpectations(t)
}

func TestGauge(t *testing.T) {
	client, logger, metric, value, tagsMap, logIndexes := setupForValues("Gauge", "Warn", false, false)
	m := NewMetrics(client, logger)

	m.Gauge(metric, value, tagsMap, logIndexes)
	client.AssertExpectations(t)
}

func TestGaugeNoLogTags(t *testing.T) {
	client, logger, metric, value, tagsMap, _ := setupForValues("Gauge", "Warn", false, true)
	m := NewMetrics(client, logger)

	m.GaugeNoLogTags(metric, value, tagsMap)
	client.AssertExpectations(t)
}

func TestGaugeNoDDTags(t *testing.T) {
	client, logger, metric, value, _, logIndexes := setupForValues("Gauge", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.GaugeNoDDTags(metric, value, logIndexes)
	client.AssertExpectations(t)
}

func TestGaugeNoTags(t *testing.T) {
	client, logger, metric, value, _, _ := setupForValues("Gauge", "Warn", true, true)
	m := NewMetrics(client, logger)

	m.GaugeNoTags(metric, value)
	client.AssertExpectations(t)
}

func TestTiming(t *testing.T) {
	client, logger, metric, value, tagsMap, logIndexes := setupForTimings("Timing", "Warn", false, false)
	m := NewMetrics(client, logger)

	m.Timing(metric, value, tagsMap, logIndexes)
	client.AssertExpectations(t)
}

func TestTimingNoLogTags(t *testing.T) {
	client, logger, metric, value, tagsMap, _ := setupForTimings("Timing", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.TimingNoLogTags(metric, value, tagsMap)
	client.AssertExpectations(t)
}

func TestTimingNoDDTags(t *testing.T) {
	client, logger, metric, value, _, logIndexes := setupForTimings("Timing", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.TimingNoDDTags(metric, value, logIndexes)
	client.AssertExpectations(t)
}

func TestTimingNoTags(t *testing.T) {
	client, logger, metric, value, _, _ := setupForTimings("Timing", "Warn", true, true)
	m := NewMetrics(client, logger)

	m.TimingNoTags(metric, value)
	client.AssertExpectations(t)
}

func TestTimingNoSampling(t *testing.T) {
	client, logger, metric, value, tagsMap, logIndexes := setupForTimings("TimingNoSampling", "Warn", false, false)
	m := NewMetrics(client, logger)

	m.TimingNoSampling(metric, value, tagsMap, logIndexes)
	client.AssertExpectations(t)
}

func TestTimingNoSamplingOrLogTags(t *testing.T) {
	client, logger, metric, value, tagsMap, _ := setupForTimings("TimingNoSampling", "Warn", false, true)
	m := NewMetrics(client, logger)

	m.TimingNoSamplingOrLogTags(metric, value, tagsMap)
	client.AssertExpectations(t)
}

func TestTimingNoSamplingOrDDTags(t *testing.T) {
	client, logger, metric, value, _, logIndexes := setupForTimings("TimingNoSampling", "Warn", true, false)
	m := NewMetrics(client, logger)

	m.TimingNoSamplingOrDDTags(metric, value, logIndexes)
	client.AssertExpectations(t)
}

func TestTimingNoSamplingOrTags(t *testing.T) {
	client, logger, metric, value, _, _ := setupForTimings("TimingNoSampling", "Warn", true, true)
	m := NewMetrics(client, logger)

	m.TimingNoSamplingOrTags(metric, value)
	client.AssertExpectations(t)
}
