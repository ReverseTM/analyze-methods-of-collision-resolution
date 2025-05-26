package double

import "math/bits"

const hashConst uint64 = 0xbf58476d1ce4e5b9

type entry struct {
	key   int
	value any
	state uint8
}

type HashTable struct {
	table      []entry
	size       int
	cap        int
	loadFactor float64
	probes     int
	collisions int
}

func New(initialCapacity int) *HashTable {
	capacity := nextPowerOfTwo(initialCapacity)

	return &HashTable{
		table:      make([]entry, capacity),
		size:       0,
		cap:        capacity,
		loadFactor: 0.7,
	}
}

func (ht *HashTable) Insert(key int, value any) {
	if ht.shouldResize() {
		ht.resize()
	}

	ok := ht.insertNoResize(key, value, true)
	if !ok {
		if !ht.resize() {
			return
		}

		ok = ht.insertNoResize(key, value, false)
		if !ok {
			return
		}
	}
}

func (ht *HashTable) insertNoResize(key int, value any, withCollision bool) bool {
	h1 := ht.hash1(key)
	h2 := ht.hash2(key)
	firstTombstone := -1

	for i := 0; i < ht.cap; i++ {
		ht.probes++

		idx := (h1 + i*h2) & (ht.cap - 1)
		state := ht.table[idx].state

		if i == 0 && state == 1 && ht.table[idx].key != key && withCollision {
			ht.collisions++
		}

		if state == 0 {
			target := idx

			if firstTombstone != -1 {
				target = firstTombstone
			}

			htSlot := &ht.table[target]
			htSlot.key = key
			htSlot.value = value
			htSlot.state = 1
			ht.size++

			return true
		}

		if state == 2 {
			if firstTombstone == -1 {
				firstTombstone = idx
			}

			continue
		}

		if ht.table[idx].key == key {
			ht.table[idx].value = value
			return true
		}
	}

	return false
}

func (ht *HashTable) Get(key int) (any, bool) {
	h1 := ht.hash1(key)
	h2 := ht.hash2(key)

	for i := 0; i < ht.cap; i++ {
		ht.probes++

		idx := (h1 + i*h2) & (ht.cap - 1)
		ent := ht.table[idx]

		if ent.state == 0 {
			break
		}

		if ent.state == 1 && ent.key == key {
			return ent.value, true
		}
	}

	return nil, false
}

func (ht *HashTable) Delete(key int) {
	h1 := ht.hash1(key)
	h2 := ht.hash2(key)

	for i := 0; i < ht.cap; i++ {
		ht.probes++

		idx := (h1 + i*h2) & (ht.cap - 1)
		ent := &ht.table[idx]

		if ent.state == 0 {
			return
		}

		if ent.state == 1 && ent.key == key {
			ent.state = 2
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

func (ht *HashTable) resize() bool {
	old := ht.table
	capacity := ht.cap * 2

	ht.table = make([]entry, capacity)
	ht.size = 0
	ht.cap = capacity

	for _, e := range old {
		if e.state == 1 {
			ok := ht.insertNoResize(e.key, e.value, false)
			if !ok {
				return false
			}
		}
	}

	return true
}

func (ht *HashTable) shouldResize() bool {
	return float64(ht.size)/float64(ht.cap) >= ht.loadFactor
}

func (ht *HashTable) hash1(key int) int {
	return int((uint64(key) * hashConst) & uint64(ht.cap-1))
}

func (ht *HashTable) hash2(key int) int {
	return int(1 + uint64(key)&uint64(ht.cap-1))
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
