package hash_table

type HashTable interface {
	Insert(key int, value any)
	Get(key int) (any, bool)
	Delete(key int)
	SetLoadFactor(loadFactor float64)
	Probes() int
	ResetProbes()
	Collisions() int
	ResetCollisions()
	Size() int
	Capacity() int
}
