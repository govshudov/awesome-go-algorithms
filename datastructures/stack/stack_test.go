package stack

import "testing"

func TestStack_LIFO(t *testing.T) {
	s := New[int]()
	for _, v := range []int{1, 2, 3} {
		s.Push(v)
	}
	want := []int{3, 2, 1}
	for i, w := range want {
		v, ok := s.Pop()
		if !ok || v != w {
			t.Fatalf("iter %d: Pop = (%d, %v), want (%d, true)", i, v, ok, w)
		}
	}
}

func TestStack_Peek(t *testing.T) {
	s := New[string]()
	if _, ok := s.Peek(); ok {
		t.Fatal("Peek on empty stack returned ok=true")
	}
	s.Push("a")
	s.Push("b")
	v, ok := s.Peek()
	if !ok || v != "b" {
		t.Fatalf("Peek = (%q, %v), want (\"b\", true)", v, ok)
	}
	if s.Len() != 2 {
		t.Fatalf("Peek modified Len to %d", s.Len())
	}
}

func TestStack_EmptyPop(t *testing.T) {
	s := New[int]()
	if v, ok := s.Pop(); ok || v != 0 {
		t.Fatalf("Pop on empty = (%d, %v), want (0, false)", v, ok)
	}
}
