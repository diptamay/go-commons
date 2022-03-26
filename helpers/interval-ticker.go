package helpers

import (
	"time"
)

// ClearInterval stops an existing ticker
func ClearInterval(ticker *IntervalTicker) {
	ticker.Stop()
	ticker.quit <- true
}

// IntervalTicker defines a ticker with a specified interval
type IntervalTicker struct {
	*time.Ticker
	quit chan bool
}

// SetInterval creates a ticker that executes a given function
func SetInterval(action func(), interval int) *IntervalTicker {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				action()
			case <-quit:
				break
			}
		}
	}()
	return &IntervalTicker{ticker, quit}
}
