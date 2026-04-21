package faninout

import (
	"context"
	"sort"
	"testing"
)

func TestFanOutFanIn_RoundTrip(t *testing.T) {
	ctx := context.Background()
	in := make(chan int)
	go func() {
		defer close(in)
		for i := range 20 {
			in <- i
		}
	}()

	outs := FanOut(ctx, in, 4)
	merged := FanIn(ctx, outs...)

	var got []int
	for v := range merged {
		got = append(got, v)
	}
	sort.Ints(got)

	want := make([]int, 20)
	for i := range want {
		want[i] = i
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestFanIn_MultipleSources(t *testing.T) {
	ctx := context.Background()
	make1 := func(start, count int) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := range count {
				ch <- start + i
			}
		}()
		return ch
	}
	merged := FanIn(ctx, make1(0, 5), make1(100, 5), make1(200, 5))
	var got []int
	for v := range merged {
		got = append(got, v)
	}
	if len(got) != 15 {
		t.Fatalf("len = %d, want 15", len(got))
	}
}

func TestFanOut_CancelPropagates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan int)
	go func() {
		defer close(in)
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			case in <- i:
			}
		}
	}()

	outs := FanOut(ctx, in, 2)
	// Read one value, then cancel.
	<-outs[0]
	cancel()

	// Both outputs should eventually close.
	for _, c := range outs {
		for range c {
		}
	}
}
