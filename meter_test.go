package metrics

import (
	"testing"
	"time"
)

func TestMeterZero(t *testing.T) {
	m := NewMeter()
	if count := m.Count(); 0 != count {
		t.Errorf("m.Count(): 0 != %v\n", count)
	}
}

func TestMeterNonzero(t *testing.T) {
	m := NewMeter()
	m.Mark(3)
	if count := m.Count(); 3 != count {
		t.Errorf("m.Count(): 3 != %v\n", count)
	}
}

func TestMeterRate1(t *testing.T) {
	m := NewMeter()
	m.Mark(3)
	time.Sleep(6 * time.Second)
	if r1 := m.Rate1(); r1 == 0 {
		t.Fatal("m.Rate1() was 0")
	}
}
