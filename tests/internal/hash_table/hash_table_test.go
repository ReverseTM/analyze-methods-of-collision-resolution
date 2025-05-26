package hash_table

import (
	"analyze/internal/hash_table/chain"
	"analyze/internal/hash_table/cuckoo"
	double "analyze/internal/hash_table/double_hash"
	"analyze/internal/hash_table/hopscotch"
	robinhood "analyze/internal/hash_table/robin_hood"
	"fmt"
	"testing"
)

func factoryMap() map[string]func(capacity int) HashTable {
	return map[string]func(int) HashTable{
		"Cuckoo":    func(c int) HashTable { return cuckoo.New(c) },
		"Chain":     func(c int) HashTable { return chain.New(c) },
		"Double":    func(c int) HashTable { return double.New(c) },
		"Hopscotch": func(c int) HashTable { return hopscotch.New(c) },
		"RobinHood": func(c int) HashTable { return robinhood.New(c) },
	}
}

func TestBasicOperations(t *testing.T) {
	for name, newTable := range factoryMap() {
		t.Run(name, func(t *testing.T) {
			ht := newTable(8)

			for i := 0; i < 10; i++ {
				val := fmt.Sprintf("val-%d", i)
				ht.Insert(i, val)
				v, found := ht.Get(i)

				if !found {
					t.Errorf("Get: expected key %d to be found", i)
				} else if v != val {
					t.Errorf("Get: value mismatch for key %d: got %v, want %v", i, v, val)
				}
			}
			for i := 0; i < 10; i += 2 {
				ht.Delete(i)
			}
			for i := 0; i < 10; i++ {
				v, found := ht.Get(i)

				if i%2 == 0 {
					if found {
						t.Errorf("Get: expected key %d to be deleted", i)
					}
				} else {
					if !found || v != fmt.Sprintf("val-%d", i) {
						t.Errorf("Get: expected key %d present with correct value, got %v, %v", i, v, found)
					}
				}
			}
		})
	}
}

func TestLargeInsertGet(t *testing.T) {
	for name, newTable := range factoryMap() {
		t.Run(name, func(t *testing.T) {
			count := 10000000
			ht := newTable(count)
			for i := 0; i < count; i++ {
				ht.Insert(i, i*2)
			}

			for i := 0; i < count; i++ {
				v, found := ht.Get(i)
				if !found {
					t.Errorf("Missing key %d", i)
				} else if v.(int) != i*2 {
					t.Errorf("Value mismatch for key %d: got %v, want %v", i, v, i*2)
				}
			}

			for i := 0; i < count; i += 3 {
				ht.Delete(i)
			}

			for i := 0; i < count; i++ {
				v, found := ht.Get(i)

				if i%3 == 0 {
					if found {
						t.Errorf("Expected key %d to be deleted", i)
					}
				} else {
					if !found || v.(int) != i*2 {
						t.Errorf("Key %d should exist with value %d", i, i*2)
					}
				}
			}
		})
	}
}

func TestUpdateValue(t *testing.T) {
	for name, newTable := range factoryMap() {
		t.Run(name, func(t *testing.T) {
			ht := newTable(16)
			ht.Insert(42, "first")
			ht.Insert(42, "second")
			v, found := ht.Get(42)
			if !found || v != "second" {
				t.Errorf("Update failed: got %v, want %v", v, "second")
			}
		})
	}
}

func TestSize(t *testing.T) {
	for name, newTable := range factoryMap() {
		t.Run(name, func(t *testing.T) {
			count := 10000000
			ht := newTable(count)

			for i := 0; i < count; i++ {
				ht.Insert(i, i*2)
			}

			if !(ht.Size() == count) {
				t.Errorf("Size failed: got %d, want %d", ht.Size(), count)
			}
		})
	}
}
