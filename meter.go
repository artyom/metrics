package metrics

import "time"

// Meters count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type Meter interface {
	// Return the count of events seen.
	Count() int64

	// Mark the occurance of n events.
	Mark(n int64)

	// Return the meter's one-minute moving average rate of events.
	Rate1() float64

	// Return the meter's five-minute moving average rate of events.
	Rate5() float64

	// Return the meter's fifteen-minute moving average rate of events.
	Rate15() float64

	// Return the meter's mean rate of events.
	RateMean() float64
}

// The standard implementation of a Meter uses a goroutine to synchronize
// its calculations and another goroutine (via time.Ticker) to produce
// clock ticks.
type meter struct {
	in     chan int64
	out    chan meterV
	ticker *time.Ticker
}

// A meterV contains all the values that would need to be passed back
// from the synchronizing goroutine.
type meterV struct {
	count                          int64
	rate1, rate5, rate15, rateMean float64
}

// Create a new meter.  Create the communication channels and start the
// synchronizing goroutine.
func NewMeter() Meter {
	m := &meter{
		make(chan int64),
		make(chan meterV),
		time.NewTicker(5 * time.Second),
	}
	go m.arbiter()
	return m
}

func (m *meter) Count() int64 {
	return (<-m.out).count
}

func (m *meter) Mark(n int64) {
	m.in <- n
}

func (m *meter) Rate1() float64 {
	return (<-m.out).rate1
}

func (m *meter) Rate5() float64 {
	return (<-m.out).rate5
}

func (m *meter) Rate15() float64 {
	return (<-m.out).rate15
}

func (m *meter) RateMean() float64 {
	return (<-m.out).rateMean
}

// Receive inputs and send outputs.  Count each input and update the various
// moving averages and the mean rate of events.  Send a copy of the meterV
// as output.
func (m *meter) arbiter() {
	var mv meterV
	a1 := NewEWMA1()
	a5 := NewEWMA5()
	a15 := NewEWMA15()
	t := time.Now()
	set := func() {
		mv.rate1 = a1.Rate()
		mv.rate5 = a5.Rate()
		mv.rate15 = a15.Rate()
		mv.rateMean = float64(1e9*mv.count) / float64(time.Since(t))
	}
	for {
		select {
		case n := <-m.in:
			mv.count += n
			a1.Update(n)
			a5.Update(n)
			a15.Update(n)
			set()
		case m.out <- mv:
		case <-m.ticker.C:
			a1.Tick()
			a5.Tick()
			a15.Tick()
			set()
		}
	}
}
