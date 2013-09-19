package metrics

// Healthchecks hold an os.Error value describing an arbitrary up/down status.
//
// This is an interface so as to encourage other structs to implement
// the Healthcheck API as appropriate.
type Healthcheck interface {
	// Update the healthcheck's status.
	Check()

	// Return the healthcheck's status, which will be nil if it is healthy.
	Error() error

	// Mark the healthcheck as healthy.
	Healthy()

	// Mark the healthcheck as unhealthy.  The error should provide details.
	Unhealthy(error)
}

// The standard implementation of a Healthcheck stores the status and a
// function to call to update the status.
type healthcheck struct {
	err error
	f   func(Healthcheck)
}

// Create a new healthcheck, which will use the given function to update
// its status.
func NewHealthcheck(f func(Healthcheck)) Healthcheck {
	return &healthcheck{nil, f}
}

func (h *healthcheck) Check() {
	h.f(h)
}

func (h *healthcheck) Error() error {
	return h.err
}

func (h *healthcheck) Healthy() {
	h.err = nil
}

func (h *healthcheck) Unhealthy(err error) {
	h.err = err
}
