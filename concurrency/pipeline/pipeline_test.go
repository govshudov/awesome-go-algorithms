package pipeline

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestPipeline_MapFilterCollect(t *testing.T) {
	ctx := context.Background()
	nums := Source(ctx, []int{1, 2, 3, 4, 5})
	doubled := Map(ctx, nums, func(n int) int { return n * 2 })
	evens := Filter(ctx, doubled, func(n int) bool { return n%4 == 0 })

	got := Collect(ctx, evens)
	want := []int{4, 8}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestPipeline_CancellationStopsSource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	// Build a long source; cancel before draining.
	big := make([]int, 10_000)
	for i := range big {
		big[i] = i
	}
	src := Source(ctx, big)

	first, ok := <-src
	if !ok || first != 0 {
		t.Fatalf("first value = (%d, %v), want (0, true)", first, ok)
	}
	cancel()

	// Eventually the source should close.
	timeout := time.After(time.Second)
	for {
		select {
		case _, ok := <-src:
			if !ok {
				return
			}
		case <-timeout:
			t.Fatal("source did not close after cancel")
		}
	}
}

func TestPipeline_TypeTransformation(t *testing.T) {
	ctx := context.Background()
	words := Source(ctx, []string{"a", "bb", "ccc"})
	lengths := Map(ctx, words, func(s string) int { return len(s) })
	got := Collect(ctx, lengths)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
