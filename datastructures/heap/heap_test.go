package heap

import (
	"math/rand"
	"sort"
	"testing"
)

func TestMinHeap(t *testing.T) {
	h := New[int](func(a, b int) bool { return a < b })
	for _, v := range []int{5, 1, 9, 3, 7, 2} {
		h.Push(v)
	}
	want := []int{1, 2, 3, 5, 7, 9}
	for i, w := range want {
		v, ok := h.Pop()
		if !ok || v != w {
			t.Fatalf("iter %d: Pop = (%d, %v), want (%d, true)", i, v, ok, w)
		}
	}
}

func TestMaxHeap(t *testing.T) {
	h := New[int](func(a, b int) bool { return a > b })
	for _, v := range []int{5, 1, 9, 3, 7, 2} {
		h.Push(v)
	}
	want := []int{9, 7, 5, 3, 2, 1}
	for i, w := range want {
		v, _ := h.Pop()
		if v != w {
			t.Fatalf("iter %d: Pop = %d, want %d", i, v, w)
		}
	}
}

func TestHeap_RandomMatchesSort(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	values := make([]int, 200)
	for i := range values {
		values[i] = r.Intn(1000)
	}
	h := New[int](func(a, b int) bool { return a < b })
	for _, v := range values {
		h.Push(v)
	}
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Ints(sorted)
	for i, want := range sorted {
		v, _ := h.Pop()
		if v != want {
			t.Fatalf("iter %d: got %d, want %d", i, v, want)
		}
	}
}

func TestHeap_EmptyPop(t *testing.T) {
	h := New[int](func(a, b int) bool { return a < b })
	if _, ok := h.Pop(); ok {
		t.Fatal("Pop on empty returned ok=true")
	}
	if _, ok := h.Peek(); ok {
		t.Fatal("Peek on empty returned ok=true")
	}
}
