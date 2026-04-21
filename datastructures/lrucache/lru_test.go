package lrucache

import "testing"

func TestCache_GetPut(t *testing.T) {
	c := New[string, int](3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = (%d, %v), want (1, true)", v, ok)
	}
	if v, ok := c.Get("missing"); ok || v != 0 {
		t.Fatalf("Get(missing) = (%d, %v), want (0, false)", v, ok)
	}
}

func TestCache_EvictsLeastRecentlyUsed(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a")    // a becomes most recent
	c.Put("c", 3) // b should evict

	if _, ok := c.Get("b"); ok {
		t.Error("b should have been evicted")
	}
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Errorf("a should still be present, got (%d, %v)", v, ok)
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Errorf("c missing, got (%d, %v)", v, ok)
	}
}

func TestCache_UpdateExisting(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("a", 100)
	if v, _ := c.Get("a"); v != 100 {
		t.Errorf("updated value = %d, want 100", v)
	}
	if c.Len() != 1 {
		t.Errorf("Len = %d, want 1", c.Len())
	}
}

func TestCache_ZeroCapacityCoercedToOne(t *testing.T) {
	c := New[int, int](0)
	c.Put(1, 1)
	c.Put(2, 2)
	if _, ok := c.Get(1); ok {
		t.Error("key 1 should have been evicted when key 2 was inserted")
	}
}
