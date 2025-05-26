package robinhood

import (
	"math/bits"
)

type state uint8

const (
	empty state = iota
	occupied
	tomb
)

const (
	hashConst uint64 = 0xbf58476d1ce4e5b9
)

type bucket struct {
	key   int
	value any
	flag  state
}

type HashTable struct {
	table      []bucket
	size       int
	cap        int
	loadFactor float64
	probes     int
	collisions int
}

func New(initialCapacity int) *HashTable {
	capacity := nextPowerOfTwo(initialCapacity)

	return &HashTable{
		table:      make([]bucket, capacity),
		size:       0,
		cap:        capacity,
		loadFactor: 0.7,
	}
}

func (ht *HashTable) Insert(key int, value any) {
	if ht.shouldResize() {
		ht.resize()
	}

	idx := ht.hash(key)
	dist := 0
	collisionCounted := false

	for {
		ht.probes++
		b := &ht.table[idx]

		switch b.flag {
		case empty, tomb:
			*b = bucket{key: key, value: value, flag: occupied}
			ht.size++
			return

		case occupied:
			if b.key == key {
				b.value = value
				return
			}

			if !collisionCounted {
				ht.collisions++
				collisionCounted = true
			}

			home := ht.hash(b.key)
			existingDist := (idx - home) & (ht.cap - 1)
			if existingDist < dist {
				key, b.key = b.key, key
				value, b.value = b.value, value
				dist = existingDist
			}

			dist++
			idx = (idx + 1) & (ht.cap - 1)
		}
	}
}

func (ht *HashTable) Get(key int) (any, bool) {
	idx := ht.hash(key)
	dist := 0

	for {
		ht.probes++

		b := &ht.table[idx]
		switch b.flag {
		case empty:
			return nil, false
		case occupied:
			if b.key == key {
				return b.value, true
			}

			home := ht.hash(b.key)
			existingDist := (idx - home) & (ht.cap - 1)
			if existingDist < dist {
				return nil, false
			}
		}

		dist++
		idx = (idx + 1) & (ht.cap - 1)

		if dist > ht.cap {
			return nil, false
		}
	}
}

func (ht *HashTable) Delete(key int) {
	idx := ht.hash(key)
	dist := 0

	for {
		ht.probes++

		b := &ht.table[idx]
		switch b.flag {
		case empty:
			return
		case occupied:
			if b.key == key {
				b.flag = tomb
				ht.size--

				return
			}

			home := ht.hash(b.key)
			existingDist := (idx - home) & (ht.cap - 1)
			if existingDist < dist {
				return
			}
		}

		dist++
		idx = (idx + 1) & (ht.cap - 1)
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
	old := ht.table
	oldCollisions := ht.collisions
	capacity := ht.cap * 2

	ht.table = make([]bucket, capacity)
	ht.size = 0
	ht.cap = capacity

	for _, e := range old {
		if e.flag == occupied {
			ht.Insert(e.key, e.value)
		}
	}

	ht.collisions = oldCollisions
}

func (ht *HashTable) hash(key int) int {
	return int((uint64(key) * hashConst) & uint64(ht.cap-1))
}

func (ht *HashTable) shouldResize() bool {
	return float64(ht.size)/float64(ht.cap) >= ht.loadFactor
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
