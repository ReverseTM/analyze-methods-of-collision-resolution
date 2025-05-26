package chain

import (
	"math/bits"
)

const hashConst uint64 = 0xbf58476d1ce4e5b9

type entry struct {
	key   int
	value any
}

type HashTable struct {
	buckets    [][]entry
	size       int
	cap        int
	loadFactor float64
	probes     int
	collisions int
}

func New(initialCapacity int) *HashTable {
	capacity := nextPowerOfTwo(initialCapacity)

	return &HashTable{
		buckets:    make([][]entry, capacity),
		size:       0,
		cap:        capacity,
		loadFactor: 1.,
	}
}

func (ht *HashTable) Insert(key int, value any) {
	if ht.shouldResize() {
		ht.resize()
	}

	ht.insertNoResize(key, value)
}

func (ht *HashTable) insertNoResize(key int, value any) {
	idx := ht.hash(key)

	for i := range ht.buckets[idx] {
		ht.probes++

		if ht.buckets[idx][i].key == key {
			ht.buckets[idx][i].value = value
			return
		}
	}

	if len(ht.buckets[idx]) > 0 {
		ht.collisions++
	}

	ht.buckets[idx] = append(ht.buckets[idx], entry{key, value})
	ht.size++

	return
}

func (ht *HashTable) Get(key int) (any, bool) {
	idx := ht.hash(key)

	for _, e := range ht.buckets[idx] {
		ht.probes++

		if e.key == key {
			return e.value, true
		}
	}

	return nil, false
}

func (ht *HashTable) Delete(key int) {
	idx := ht.hash(key)
	chain := ht.buckets[idx]

	for i, e := range chain {
		ht.probes++

		if e.key == key {
			last := len(chain) - 1
			chain[i] = chain[last]
			ht.buckets[idx] = chain[:last]
			ht.size--
			return
		}
	}
}

func (ht *HashTable) SetLoadFactor(loadFactor float64) {
	ht.loadFactor = loadFactor
}

func (ht *HashTable) Probes() int {
	return ht.probes
}

func (ht *HashTable) ResetProbes() {
	ht.probes = 0
}

func (ht *HashTable) Collisions() int {
	return ht.collisions
}

func (ht *HashTable) ResetCollisions() {
	ht.collisions = 0
}

func (ht *HashTable) Size() int {
	return ht.size
}

func (ht *HashTable) Capacity() int {
	return ht.cap
}

func (ht *HashTable) resize() {
	old := ht.buckets
	oldCollision := ht.collisions
	capacity := ht.cap * 2

	ht.buckets = make([][]entry, capacity)
	ht.size = 0
	ht.cap = capacity

	for _, chain := range old {
		for _, e := range chain {
			ht.insertNoResize(e.key, e.value)
		}
	}

	ht.collisions = oldCollision
}

func (ht *HashTable) shouldResize() bool {
	return float64(ht.size)/float64(ht.cap) >= ht.loadFactor
}

func (ht *HashTable) hash(key int) int {
	return int((uint64(key) * hashConst) & uint64(ht.cap-1))
}

func nextPowerOfTwo(n int) int {
	if n < 8 {
		return 8
	}
	if (n & (n - 1)) == 0 {
		return n
	}
	return 1 << (bits.Len(uint(n)))
}
