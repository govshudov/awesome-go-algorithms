// Package unionfind implements a disjoint-set (union-find) data structure
// with path compression and union by rank. Near-constant amortized time per
// operation (inverse Ackermann).
package unionfind

type UnionFind struct {
	parent []int
	rank   []int
	sets   int
}

func New(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
		sets:   n,
	}
	for i := range uf.parent {
		uf.parent[i] = i
	}
	return uf
}

func (u *UnionFind) Find(x int) int {
	for u.parent[x] != x {
		u.parent[x] = u.parent[u.parent[x]]
		x = u.parent[x]
	}
	return x
}

// Union merges the sets containing x and y. Returns true if they were
// previously in different sets.
func (u *UnionFind) Union(x, y int) bool {
	rx, ry := u.Find(x), u.Find(y)
	if rx == ry {
		return false
	}
	switch {
	case u.rank[rx] < u.rank[ry]:
		u.parent[rx] = ry
	case u.rank[rx] > u.rank[ry]:
		u.parent[ry] = rx
	default:
		u.parent[ry] = rx
		u.rank[rx]++
	}
	u.sets--
	return true
}

func (u *UnionFind) Connected(x, y int) bool { return u.Find(x) == u.Find(y) }

func (u *UnionFind) Sets() int { return u.sets }
