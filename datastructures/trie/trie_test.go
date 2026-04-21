package trie

import "testing"

func TestTrie_InsertSearch(t *testing.T) {
	tr := New()
	words := []string{"apple", "app", "apply", "banana"}
	for _, w := range words {
		tr.Insert(w)
	}

	tests := []struct {
		word string
		want bool
	}{
		{"apple", true},
		{"app", true},
		{"apply", true},
		{"banana", true},
		{"ap", false},
		{"appl", false},
		{"bananas", false},
		{"", false},
	}
	for _, tc := range tests {
		if got := tr.Search(tc.word); got != tc.want {
			t.Errorf("Search(%q) = %v, want %v", tc.word, got, tc.want)
		}
	}
}

func TestTrie_HasPrefix(t *testing.T) {
	tr := New()
	tr.Insert("hello")
	tr.Insert("help")

	tests := []struct {
		prefix string
		want   bool
	}{
		{"hel", true},
		{"he", true},
		{"help", true},
		{"helper", false},
		{"world", false},
	}
	for _, tc := range tests {
		if got := tr.HasPrefix(tc.prefix); got != tc.want {
			t.Errorf("HasPrefix(%q) = %v, want %v", tc.prefix, got, tc.want)
		}
	}
}

func TestTrie_Delete(t *testing.T) {
	tr := New()
	tr.Insert("apple")
	tr.Insert("app")

	if !tr.Delete("apple") {
		t.Fatal("Delete(apple) returned false")
	}
	if tr.Search("apple") {
		t.Error("apple still found after Delete")
	}
	if !tr.Search("app") {
		t.Error("app lost after deleting apple")
	}
	if tr.Delete("unknown") {
		t.Error("Delete(unknown) returned true")
	}
}

func TestTrie_Unicode(t *testing.T) {
	tr := New()
	tr.Insert("İstanbul")
	tr.Insert("çay")
	if !tr.Search("İstanbul") {
		t.Error("failed to find Unicode word")
	}
	if !tr.HasPrefix("ça") {
		t.Error("failed Unicode prefix match")
	}
}
