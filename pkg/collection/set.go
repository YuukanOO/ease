package collection

import "sync"

// Represents a unique set of objects.
type Set[T any] struct {
	mu      sync.Mutex
	items   []T
	indexes map[string]int
}

// Builds an empty set.
func NewSet[T any]() *Set[T] {
	return &Set[T]{
		indexes: make(map[string]int),
	}
}

// Set the given item for the given key if it does not exist. Returns the item.
func (s *Set[T]) Set(key string, item T) T {
	if idx, found := s.indexes[key]; found {
		return s.items[idx]
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.indexes[key] = len(s.items)
	s.items = append(s.items, item)
	return item
}

// Same as Set but build the item only if not already found in the set to prevent
// unneeded allocations.
func (s *Set[T]) SetLazy(key string, item func() T) T {
	if idx, found := s.indexes[key]; found {
		return s.items[idx]
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	created := item()
	s.indexes[key] = len(s.items)
	s.items = append(s.items, created)
	return created
}

// Retrieve all items inside the set.
func (s *Set[T]) Items() []T { return s.items }
