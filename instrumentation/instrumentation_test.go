package instrumentation

import (
	"testing"
	"time"

	. "github.com/diptamay/go-commons/helpers"
	"github.com/diptamay/go-commons/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitializeMemoryLogger(t *testing.T) {
	statsMock := &mocks.MockMetrics{}
	statsMock.
		On("GaugeNoTags", "service.core.memory.sys", mock.AnythingOfType("float64")).
		Return().
		Once()
	statsMock.
		On("GaugeNoTags", "service.core.memory.heap_total", mock.AnythingOfType("float64")).
		Return().
		Once()
	statsMock.
		On("GaugeNoTags", "service.core.memory.heap_used", mock.AnythingOfType("float64")).
		Return().
		Once()

	timer := time.NewTimer(time.Duration(1) * time.Second)
	ticker, ok := InitializeMemoryLogger(statsMock, 750)

	assert.Equal(t, true, ok, "should initialize properly")
	<-timer.C
	ClearInterval(ticker)

	statsMock.AssertExpectations(t)
}

func TestInitializeLoadLogger(t *testing.T) {
	statsMock := &mocks.MockMetrics{}
	statsMock.
		On("GaugeNoTags", "service.core.load.1", mock.AnythingOfType("float64")).
		Return().
		Once()
	statsMock.
		On("GaugeNoTags", "service.core.load.5", mock.AnythingOfType("float64")).
		Return().
		Once()
	statsMock.
		On("GaugeNoTags", "service.core.load.15", mock.AnythingOfType("float64")).
		Return().
		Once()

	timer := time.NewTimer(time.Duration(1) * time.Second)
	ticker, ok := InitializeLoadLogger(statsMock, 750)

	assert.Equal(t, true, ok, "should initialize properly")
	<-timer.C
	ClearInterval(ticker)

	statsMock.AssertExpectations(t)
}

func TestInitializeGoroutineLogger(t *testing.T) {
	statsMock := &mocks.MockMetrics{}
	statsMock.
		On("GaugeNoTags", "service.core.numGoroutines", mock.AnythingOfType("float64")).
		Return().
		Once()

	timer := time.NewTimer(time.Duration(1) * time.Second)
	ticker, ok := InitializeGoroutineLogger(statsMock, 750)

	assert.Equal(t, true, ok, "should initialize properly")
	<-timer.C
	ClearInterval(ticker)

	statsMock.AssertExpectations(t)
}

func TestInitializeGCLogger(t *testing.T) {
	statsMock := &mocks.MockMetrics{}
	statsMock.
		On("HistogramNoTags", "service.core.gc.time_spent_in_gc", mock.AnythingOfType("float64")).
		Return().
		Once()
	statsMock.
		On("HistogramNoTags", "service.core.gc.time_spent_since_last_gc", mock.AnythingOfType("float64")).
		Return().
		Once()
	statsMock.
		On("HistogramNoTags", "service.core.gc.percentage_time_in_gc", mock.AnythingOfType("float64")).
		Return().
		Once()

	timer := time.NewTimer(time.Duration(1) * time.Second)
	ticker, ok := InitializeGCLogger(statsMock, 750)

	assert.Equal(t, true, ok, "should initialize properly")
	<-timer.C
	ClearInterval(ticker)

	statsMock.AssertExpectations(t)
}
