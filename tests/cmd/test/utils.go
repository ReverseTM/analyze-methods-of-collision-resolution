package test

import (
	"analyze/internal/hash_table"
	"analyze/internal/hash_table/chain"
	"analyze/internal/hash_table/cuckoo"
	double "analyze/internal/hash_table/double_hash"
	"analyze/internal/hash_table/hopscotch"
	robinhood "analyze/internal/hash_table/robin_hood"
	"iter"
	"math/rand"
)

// CONSTANT
type reserveStrategy int8

const (
	reserveNone reserveStrategy = iota
	reserveExact
)

type lookupStrategy int8

const (
	lookupSuccess lookupStrategy = iota
	lookupMiss
)

var (
	Random = rand.New(rand.NewSource(25))

	Sizes = []int{1, 10, 100, 1000, 10_000, 100_000, 1_000_000, 10_000_000}

	LoadFactors = []float64{0.4, 0.6, 0.8}

	Factories = map[string]func(cap int) hash_table.HashTable{
		"Chain":     func(c int) hash_table.HashTable { return chain.New(c) },
		"Cuckoo":    func(c int) hash_table.HashTable { return cuckoo.New(c) },
		"Double":    func(c int) hash_table.HashTable { return double.New(c) },
		"Hopscotch": func(c int) hash_table.HashTable { return hopscotch.New(c) },
		"RobinHood": func(c int) hash_table.HashTable { return robinhood.New(c) },
	}

	KeyGens = map[string]func(int) iter.Seq[int]{
		"RandomKey":     func(count int) iter.Seq[int] { return genRandomKeys(count) },
		"SequentialKey": func(count int) iter.Seq[int] { return genSequentialKeys(count) },
	}
)

// KEY_GENERATORS

func genRandomKeys(count int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for range count {
			if !yield(Random.Int()) {
				return
			}
		}
	}
}

func genSequentialKeys(count int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range count {
			if !yield(i) {
				return
			}
		}
	}
}
