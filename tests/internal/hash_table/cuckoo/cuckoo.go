package cuckoo

import (
	"math/bits"
	"math/rand"
	"time"
)

type entry struct {
	key      int
	value    any
	occupied bool
}

type HashTable struct {
	table1       []entry
	table2       []entry
	capMask      uint32
	size         int
	cap          int
	probes       int
	collisions   int
	maxKicks     int
	loadFactor   float64
	maxRehashes  int
	rehashCount  int
	salt1, salt2 uint64
	rng          *rand.Rand
}

func New(initialCapacity int) *HashTable {
	if initialCapacity < 1 {
		initialCapacity = 8
	}

	lf := 0.5
	minPerTable := int(float64(initialCapacity)/(2*lf)) + 1
	capacity := nextPowerOfTwo(minPerTable)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &HashTable{
		table1:      make([]entry, capacity),
		table2:      make([]entry, capacity),
		capMask:     uint32(capacity - 1),
		cap:         capacity,
		maxKicks:    500,
		loadFactor:  lf,
		maxRehashes: 5,
		salt1:       rng.Uint64(),
		salt2:       rng.Uint64(),
		rng:         rng,
	}
}

func (ht *HashTable) Insert(key int, value any) {
	newEntry := entry{key: key, value: value, occupied: true}

	firstAttempt := true

	for {
		if float64(ht.size+1) > ht.loadFactor*float64(2*len(ht.table1)) {
			ht.resizeDouble()
			ht.rehashCount = 0
		}

		if ht.insertOnce(newEntry, firstAttempt) {
			ht.rehashCount = 0
			return
		}

		firstAttempt = false

		if ht.rehashCount < ht.maxRehashes {
			ht.rehashCount++

			all := make([]entry, 0, ht.size+1)
			for _, e := range ht.table1 {
				if e.occupied {
					all = append(all, e)
				}
			}
			for _, e := range ht.table2 {
				if e.occupied {
					all = append(all, e)
				}
			}
			all = append(all, newEntry)
			if ht.rehash(all) {
				ht.rehashCount = 0
				return
			}

			continue
		}

		ht.resizeDouble()
		ht.rehashCount = 0
	}
}

func (ht *HashTable) Get(key int) (any, bool) {
	idx1 := ht.hash1(key)
	ht.probes++

	if e := ht.table1[idx1]; e.occupied && e.key == key {
		return e.value, true
	}

	idx2 := ht.hash2(key)
	ht.probes++

	if e := ht.table2[idx2]; e.occupied && e.key == key {
		return e.value, true
	}

	return nil, false
}

func (ht *HashTable) Delete(key int) {
	idx1 := ht.hash1(key)
	ht.probes++

	if e := &ht.table1[idx1]; e.occupied && e.key == key {
		e.occupied = false
		ht.size--
		return
	}

	idx2 := ht.hash2(key)
	ht.probes++

	if e := &ht.table2[idx2]; e.occupied && e.key == key {
		e.occupied = false
		ht.size--
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

func (ht *HashTable) insertOnce(e entry, withCollision bool) bool {
	curKey, curVal := e.key, e.value
	table := 0
	for kick := 0; kick < ht.maxKicks; kick++ {
		ht.probes++

		if table == 0 {
			idx := ht.hash1(curKey)
			slot := &ht.table1[idx]
			if !slot.occupied {
				slot.key, slot.value, slot.occupied = curKey, curVal, true
				ht.size++
				return true
			}
			if slot.key == curKey {
				slot.value = curVal
				return true
			}

			if kick == 0 {
				idx2 := ht.hash2(curKey)
				alt := &ht.table2[idx2]
				if alt.occupied && alt.key != curKey && withCollision {
					ht.collisions++
				}
			}

			slot.key, curKey = curKey, slot.key
			slot.value, curVal = curVal, slot.value
		} else {
			idx := ht.hash2(curKey)
			slot := &ht.table2[idx]
			if !slot.occupied {
				slot.key, slot.value, slot.occupied = curKey, curVal, true
				ht.size++
				return true
			}
			if slot.key == curKey {
				slot.value = curVal
				return true
			}
			slot.key, curKey = curKey, slot.key
			slot.value, curVal = curVal, slot.value
		}

		table ^= 1
	}
	return false
}

func (ht *HashTable) rehash(all []entry) bool {
	newSalt1 := ht.rng.Uint64()
	newSalt2 := ht.rng.Uint64()
	n := len(ht.table1)

	t1 := make([]entry, n)
	t2 := make([]entry, n)

	for _, e := range all {
		curKey, curVal := e.key, e.value
		table := 0
		placed := false
		for kick := 0; kick < ht.maxKicks; kick++ {
			ht.probes++

			if table == 0 {
				idx := uint32(splitmix(uint64(curKey)^newSalt1)) & uint32(n-1)
				slot := &t1[idx]
				if !slot.occupied {
					slot.key, slot.value, slot.occupied = curKey, curVal, true
					placed = true
					break
				}
				if slot.key == curKey {
					slot.value = curVal
					placed = true
					break
				}
				slot.key, curKey = curKey, slot.key
				slot.value, curVal = curVal, slot.value
			} else {
				idx := uint32(splitmix(uint64(curKey)^newSalt2)) & uint32(n-1)
				slot := &t2[idx]
				if !slot.occupied {
					slot.key, slot.value, slot.occupied = curKey, curVal, true
					placed = true
					break
				}
				if slot.key == curKey {
					slot.value = curVal
					placed = true
					break
				}
				slot.key, curKey = curKey, slot.key
				slot.value, curVal = curVal, slot.value
			}
			table ^= 1
		}
		if !placed {
			return false
		}
	}

	ht.table1 = t1
	ht.table2 = t2
	ht.salt1 = newSalt1
	ht.salt2 = newSalt2
	ht.size = len(all)
	return true
}

func (ht *HashTable) resizeDouble() {
	old := make([]entry, 0, ht.size)
	for _, e := range ht.table1 {
		if e.occupied {
			old = append(old, e)
		}
	}
	for _, e := range ht.table2 {
		if e.occupied {
			old = append(old, e)
		}
	}

	newCap := len(ht.table1) * 4

	ht.table1 = make([]entry, newCap)
	ht.table2 = make([]entry, newCap)
	ht.capMask = uint32(newCap - 1)
	ht.size = 0
	for _, e := range old {
		ht.insertOnce(e, false)
	}
}

func (ht *HashTable) hash1(key int) uint32 {
	return uint32(uint64(key)*ht.salt1) & ht.capMask
}
func (ht *HashTable) hash2(key int) uint32 {
	return uint32(uint64(key)*ht.salt2) & ht.capMask
}

func splitmix(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
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
