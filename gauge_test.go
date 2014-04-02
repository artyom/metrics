package metrics_test

import (
	"testing"

	"github.com/facebookgo/metrics"
)

func TestGauge(t *testing.T) {
	g := metrics.NewGauge()
	g.Update(int64(47))
	if v := g.Value(); 47 != v {
		t.Errorf("g.Value(): 47 != %v\n", v)
	}
}
