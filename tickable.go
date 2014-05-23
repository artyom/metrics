package metrics

import (
	"time"
)

// TickDuration defines the rate at which Tick() should get called for EWMA and
// other Tickable things that rely on EWMA like Meter & Timer.
//
// It is caller's responsibility to call Tick() method; the usual case is to
// spawn a goroutine to do this:
//
//	e := metrics.NewEWMA1()
//	go func(f func()){
//		for _ = range time.Tick(metrics.TickDuration) {
//			f()
//		}
//	}(e.Tick)
const TickDuration = 5 * time.Second

// Tickable defines the interface implemented by metrics that need to Tick.
type Tickable interface {
	// Tick the clock to update the moving average.
	Tick()
}
