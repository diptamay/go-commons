package instrumentation

import (
	"runtime"
	"time"

	. "github.com/diptamay/go-commons/helpers"
	dd "github.com/diptamay/go-commons/metrics"
	"github.com/shirou/gopsutil/load"
)

func InitializeMemoryLogger(metrics dd.Metrics, interval int) (*IntervalTicker, bool) {
	return SetInterval(func() {
		memstats := new(runtime.MemStats)
		runtime.ReadMemStats(memstats)
		// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		metrics.GaugeNoTags("service.core.memory.sys", float64(memstats.Sys))
		metrics.GaugeNoTags("service.core.memory.heap_total", float64(memstats.HeapSys))
		metrics.GaugeNoTags("service.core.memory.heap_used", float64(memstats.HeapAlloc))

	}, interval), true
}

func InitializeLoadLogger(metrics dd.Metrics, interval int) (*IntervalTicker, bool) {
	return SetInterval(func() {
		avg, err := load.Avg()
		if err == nil {
			metrics.GaugeNoTags("service.core.load.1", avg.Load1)
			metrics.GaugeNoTags("service.core.load.5", avg.Load5)
			metrics.GaugeNoTags("service.core.load.15", avg.Load15)
		}
	}, interval), true
}

func InitializeGoroutineLogger(metrics dd.Metrics, interval int) (*IntervalTicker, bool) {
	return SetInterval(func() {
		metrics.GaugeNoTags("service.core.numGoroutines", float64(runtime.NumGoroutine()))
	}, interval), true
}

func InitializeGCLogger(metrics dd.Metrics, interval int) (*IntervalTicker, bool) {
	return SetInterval(func() {
		memstats := new(runtime.MemStats)
		runtime.ReadMemStats(memstats)
		pauseMS := float64(memstats.PauseNs[(memstats.NumGC+255)%256] / uint64(time.Millisecond))
		timeSinceLastGC := float64((uint64(time.Now().UnixNano()) - memstats.LastGC) / uint64(time.Millisecond))

		metrics.HistogramNoTags("service.core.gc.time_spent_in_gc", pauseMS)
		metrics.HistogramNoTags("service.core.gc.time_spent_since_last_gc", timeSinceLastGC)
		metrics.HistogramNoTags("service.core.gc.percentage_time_in_gc", pauseMS/timeSinceLastGC)

	}, interval), true
}
