package unionfind

import "testing"

func TestUnionFind_Basic(t *testing.T) {
	uf := New(6)
	if uf.Sets() != 6 {
		t.Fatalf("initial Sets = %d, want 6", uf.Sets())
	}

	uf.Union(0, 1)
	uf.Union(1, 2)
	uf.Union(3, 4)

	if !uf.Connected(0, 2) {
		t.Error("0 and 2 should be connected")
	}
	if uf.Connected(0, 3) {
		t.Error("0 and 3 should not be connected")
	}
	if uf.Sets() != 3 {
		t.Errorf("Sets = %d, want 3", uf.Sets())
	}
}

func TestUnionFind_UnionReturnsFalseWhenSameSet(t *testing.T) {
	uf := New(3)
	if !uf.Union(0, 1) {
		t.Error("first Union should return true")
	}
	if uf.Union(0, 1) {
		t.Error("Union of already-joined returned true")
	}
}

func TestUnionFind_AllConnected(t *testing.T) {
	uf := New(5)
	for i := range 4 {
		uf.Union(i, i+1)
	}
	if uf.Sets() != 1 {
		t.Fatalf("Sets = %d, want 1", uf.Sets())
	}
	for i := range 5 {
		for j := range 5 {
			if !uf.Connected(i, j) {
				t.Errorf("%d and %d should be connected", i, j)
			}
		}
	}
}
