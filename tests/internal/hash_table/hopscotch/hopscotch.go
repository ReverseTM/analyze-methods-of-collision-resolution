package hopscotch

import (
	"math/bits"
)

const (
	hashConst         uint64 = 0xbf58476d1ce4e5b9
	neighbourhoodSize        = 64
	maxDistance              = 256
)

type entry struct {
	key   int
	value any
	inUse bool
}

type HashTable struct {
	buckets       []entry
	hopInfo       []uint32
	size          int
	cap           int
	loadFactor    float64
	probes        int
	collisions    int
	withCollision bool
}

func New(initialCapacity int) *HashTable {
	capacity := nextPowerOfTwo(initialCapacity)

	return &HashTable{
		buckets:       make([]entry, capacity),
		hopInfo:       make([]uint32, capacity),
		size:          0,
		cap:           capacity,
		loadFactor:    1,
		withCollision: true,
	}
}

func (ht *HashTable) Insert(key int, value any) {
	if ht.shouldResize() {
		ht.resize()
	}

	base := ht.hash(key)
	ht.probes++

	hop := ht.hopInfo[base]

	for hop != 0 {
		offset := bits.TrailingZeros32(hop)
		idx := (base + offset) & (ht.cap - 1)
		ht.probes++

		if ht.buckets[idx].inUse && ht.buckets[idx].key == key {
			ht.buckets[idx].value = value

			return
		}

		hop &= hop - 1
	}

	if ht.buckets[base].inUse && ht.withCollision {
		ht.collisions++
		ht.withCollision = false
	}

	free := base
	dist := 0

	for ; dist < maxDistance; dist++ {
		idx := (base + dist) & (ht.cap - 1)
		ht.probes++

		if !ht.buckets[idx].inUse {
			free = idx
			break
		}
	}

	if dist == maxDistance {
		ht.resize()
		ht.Insert(key, value)
		ht.withCollision = true

		return
	}

	for dist >= neighbourhoodSize {
		moved := false
		for hopDist := neighbourhoodSize - 1; hopDist > 0; hopDist-- {
			idx := (free - hopDist) & (ht.cap - 1)
			if ht.hopInfo[idx]&(1<<hopDist) != 0 {
				from := (idx + hopDist) & (ht.cap - 1)
				ht.buckets[free] = ht.buckets[from]
				ht.buckets[from].inUse = false

				ht.hopInfo[idx] &^= 1 << hopDist
				newOff := (free - idx) & (ht.cap - 1)
				ht.hopInfo[idx] |= 1 << newOff

				free = from
				dist = (free - base) & (ht.cap - 1)
				moved = true
				break
			}
		}

		if !moved {
			ht.resize()
			ht.Insert(key, value)
			ht.withCollision = true

			return
		}
	}

	ht.buckets[free] = entry{key: key, value: value, inUse: true}
	ht.hopInfo[base] |= 1 << dist
	ht.size++
}

func (ht *HashTable) Get(key int) (any, bool) {
	base := ht.hash(key)
	hop := ht.hopInfo[base]

	for hop != 0 {
		ht.probes++

		offset := bits.TrailingZeros32(hop)
		idx := (base + offset) & (ht.cap - 1)

		if ht.buckets[idx].inUse && ht.buckets[idx].key == key {
			return ht.buckets[idx].value, true
		}

		hop &= hop - 1
	}

	return nil, false
}

func (ht *HashTable) Delete(key int) {
	if len(ht.buckets) == 0 {
		return
	}

	base := ht.hash(key)
	hop := ht.hopInfo[base]

	for hop != 0 {
		ht.probes++

		offset := bits.TrailingZeros32(hop)
		idx := (base + offset) & (ht.cap - 1)

		if ht.buckets[idx].inUse && ht.buckets[idx].key == key {
			ht.buckets[idx].inUse = false
			ht.hopInfo[base] &^= 1 << offset
			ht.size--

			return
		}

		hop &= hop - 1
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

	ht.buckets = make([]entry, capacity)
	ht.hopInfo = make([]uint32, capacity)
	ht.size = 0
	ht.cap = capacity

	for _, e := range old {
		if e.inUse {
			ht.Insert(e.key, e.value)
		}
	}

	ht.collisions = oldCollision
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
