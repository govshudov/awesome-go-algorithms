package linkedlist

import (
	"reflect"
	"testing"
)

func TestList_PushPop(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushFront(0)
	l.PushBack(3)

	if got := l.ToSlice(); !reflect.DeepEqual(got, []int{0, 1, 2, 3}) {
		t.Fatalf("ToSlice = %v, want [0 1 2 3]", got)
	}
	if l.Len() != 4 {
		t.Fatalf("Len = %d, want 4", l.Len())
	}

	v, ok := l.PopFront()
	if !ok || v != 0 {
		t.Fatalf("PopFront = (%d, %v), want (0, true)", v, ok)
	}
	v, ok = l.PopBack()
	if !ok || v != 3 {
		t.Fatalf("PopBack = (%d, %v), want (3, true)", v, ok)
	}
}

func TestList_EmptyPop(t *testing.T) {
	l := New[string]()
	if v, ok := l.PopFront(); ok || v != "" {
		t.Fatalf("PopFront on empty = (%q, %v), want (\"\", false)", v, ok)
	}
	if v, ok := l.PopBack(); ok || v != "" {
		t.Fatalf("PopBack on empty = (%q, %v), want (\"\", false)", v, ok)
	}
}

func TestList_DrainMaintainsInvariants(t *testing.T) {
	l := New[int]()
	for i := range 5 {
		l.PushBack(i)
	}
	for i := range 5 {
		if v, ok := l.PopFront(); !ok || v != i {
			t.Fatalf("iter %d: got (%d, %v), want (%d, true)", i, v, ok, i)
		}
	}
	if l.Len() != 0 {
		t.Fatalf("Len = %d after draining, want 0", l.Len())
	}
	if l.head != nil || l.tail != nil {
		t.Fatalf("head/tail not nil after draining")
	}
}
