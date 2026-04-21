// Package trie provides a rune-based prefix tree supporting insert, search,
// prefix lookup, and delete. All operations run in O(len(word)).
package trie

type node struct {
	children map[rune]*node
	end      bool
}

type Trie struct {
	root *node
}

func New() *Trie { return &Trie{root: &node{children: map[rune]*node{}}} }

func (t *Trie) Insert(word string) {
	n := t.root
	for _, r := range word {
		child, ok := n.children[r]
		if !ok {
			child = &node{children: map[rune]*node{}}
			n.children[r] = child
		}
		n = child
	}
	n.end = true
}

func (t *Trie) Search(word string) bool {
	n := t.find(word)
	return n != nil && n.end
}

func (t *Trie) HasPrefix(prefix string) bool {
	return t.find(prefix) != nil
}

func (t *Trie) Delete(word string) bool {
	return deleteRec(t.root, word, 0)
}

func (t *Trie) find(s string) *node {
	n := t.root
	for _, r := range s {
		child, ok := n.children[r]
		if !ok {
			return nil
		}
		n = child
	}
	return n
}

func deleteRec(n *node, word string, depth int) bool {
	runes := []rune(word)
	if depth == len(runes) {
		if !n.end {
			return false
		}
		n.end = false
		return true
	}
	r := runes[depth]
	child, ok := n.children[r]
	if !ok {
		return false
	}
	deleted := deleteRec(child, word, depth+1)
	if deleted && !child.end && len(child.children) == 0 {
		delete(n.children, r)
	}
	return deleted
}
