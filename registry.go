package metrics

import "sync"

// A Registry holds references to a set of metrics by name and can iterate
// over them, calling callback functions provided by the user.
type Registry interface {
	// Call the given function for each registered metric.
	Each(func(string, interface{}))

	// Get the metric by the given name or nil if none is registered.
	Get(name string) interface{}

	// Register the given metric under the given name.
	Register(name string, metric interface{})

	// Run all registered healthchecks.
	RunHealthchecks()

	// Unregister the metric with the given name.
	Unregister(name string)
}

// The standard implementation of a Registry is a mutex-protected map
// of names to metrics.
type registry struct {
	metrics map[string]interface{}
	mutex   sync.Mutex
}

// Create a new registry.
func NewRegistry() Registry {
	return &registry{metrics: make(map[string]interface{})}
}

func (r *registry) Each(f func(string, interface{})) {
	for name, i := range r.registered() {
		f(name, i)
	}
}

func (r *registry) Get(name string) interface{} {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.metrics[name]
}

func (r *registry) Register(name string, i interface{}) {
	switch i.(type) {
	case Counter, Gauge, Healthcheck, Histogram, Meter, Timer:
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.metrics[name] = i
	}
}

func (r *registry) RunHealthchecks() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, i := range r.metrics {
		if h, ok := i.(Healthcheck); ok {
			h.Check()
		}
	}
}

func (r *registry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.metrics, name)
}

func (r *registry) registered() map[string]interface{} {
	metrics := make(map[string]interface{}, len(r.metrics))
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for name, i := range r.metrics {
		metrics[name] = i
	}
	return metrics
}
