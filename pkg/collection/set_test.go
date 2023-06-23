package collection_test

import (
	"testing"

	"github.com/YuukanOO/ease/pkg/collection"
)

func TestSet(t *testing.T) {
	t.Run("should set an item already if it does not exist", func(t *testing.T) {
		s := collection.NewSet[string]()

		s.Set("foo", "bar")
		s.Set("foo", "bar")

		items := s.Items()

		if len(items) != 1 {
			t.Errorf("expected set to contain 1 item, got %d", len(items))
		}

		if items[0] != "bar" {
			t.Errorf("expected set to contain item 'bar', got '%s'", items[0])
		}
	})

	t.Run("should set an item lazily if it does not exist", func(t *testing.T) {
		s := collection.NewSet[string]()
		callCount := 0

		lazyFn := func() string {
			callCount++
			return "bar"
		}

		s.SetFunc("foo", lazyFn)
		s.SetFunc("foo", lazyFn)

		items := s.Items()

		if len(items) != 1 {
			t.Errorf("expected set to contain 1 item, got %d", len(items))
		}

		if items[0] != "bar" {
			t.Errorf("expected set to contain item 'bar', got '%s'", items[0])
		}

		if callCount != 1 {
			t.Errorf("expected lazy function to be called once, got %d", callCount)
		}
	})
}
