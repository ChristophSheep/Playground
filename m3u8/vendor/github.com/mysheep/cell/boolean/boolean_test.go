package boolean

import (
	"testing"
)

func TestDistributor(t *testing.T) {

	// given
	in := make(chan bool)
	out1 := make(chan bool)
	out2 := make(chan bool)

	xs := []struct {
		given bool
		want1 bool
		want2 bool
	}{
		{true, true, true},
		{false, false, false},
	}

	go Distributor(in, out1, out2)

	for _, x := range xs {
		// when
		in <- x.given
		got1 := <-out1
		got2 := <-out2

		// then
		if got1 != x.want1 {
			t.Errorf("Distrubution was incorrect, got1:%v, want1:%v", got1, x.want1)
		}
		if got2 != x.want2 {
			t.Errorf("Distrubution was incorrect, got2:%v, want2:%v", got2, x.want2)
		}
	}
}

func TestForwarder(t *testing.T) {
	// given
	in := make(chan bool)
	out := make(chan bool)

	xs := []struct {
		given bool
		want  bool
	}{
		{true, true},
		{false, false},
	}

	go Forwarder(in, out)

	for _, x := range xs {
		// when
		in <- x.given
		got := <-out

		// then
		if got != x.want {
			t.Errorf("Forwarder was incorrect, got:%v, want:%v", got, x.want)
		}
	}
}
