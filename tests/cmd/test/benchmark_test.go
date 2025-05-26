package test

import (
	"fmt"
	"testing"
)

// BENCHMARK_TESTS

func BenchmarkInsertNoReserve(b *testing.B) {
	runInsertBenchmark(b, reserveNone)
}

func BenchmarkInsertReserve(b *testing.B) {
	runInsertBenchmark(b, reserveExact)
}

func BenchmarkSuccessGet(b *testing.B) {
	runGetBenchmark(b, lookupSuccess)
}

func BenchmarkUnsuccessGet(b *testing.B) {
	runGetBenchmark(b, lookupMiss)
}

func BenchmarkDelete(b *testing.B) {
	runDeleteBenchmark(b)
}

// BENCHMARK_FUNCTIONS

func runInsertBenchmark(b *testing.B, strategy reserveStrategy) {
	for method, newHashTable := range Factories {
		for _, size := range Sizes {
			for keyKind, keyGen := range KeyGens {
				for _, loadFactor := range LoadFactors {
					lfString := format(loadFactor)
					testName := fmt.Sprintf("%s-%s-%s-%d", method, keyKind, lfString, size)

					b.Run(testName, func(b *testing.B) {
						b.ReportAllocs()

						initCup := 8
						if strategy == reserveExact {
							initCup = size
						}

						for b.Loop() {
							b.StopTimer()
							ht := newHashTable(initCup)
							ht.SetLoadFactor(loadFactor)
							keysGen := keyGen(size)
							b.StartTimer()

							for key := range keysGen {
								ht.Insert(key, key)
							}
						}

						nsPerOp := float64(b.Elapsed().Nanoseconds()) / float64(b.N) / float64(size)

						b.ReportMetric(nsPerOp, "ns/insert")
					})
				}
			}
		}
	}
}

func runGetBenchmark(b *testing.B, strategy lookupStrategy) {
	for method, newHashTable := range Factories {
		for _, size := range Sizes {
			for keyKind, keyGen := range KeyGens {
				for _, loadFactor := range LoadFactors {
					lfString := format(loadFactor)
					testName := fmt.Sprintf("%s-%s-%s-%d", method, keyKind, lfString, size)

					b.Run(testName, func(b *testing.B) {
						ht := newHashTable(size)
						ht.SetLoadFactor(loadFactor)
						keysGen := keyGen(size)

						insertedKeys := make([]int, 0, size)
						for key := range keysGen {
							ht.Insert(key, key)
							insertedKeys = append(insertedKeys, key)
						}

						if strategy == lookupMiss {
							for i := range insertedKeys {
								insertedKeys[i] = -i
							}
						}

						Random.Shuffle(len(insertedKeys), func(i, j int) {
							insertedKeys[i], insertedKeys[j] = insertedKeys[j], insertedKeys[i]
						})

						var idx int

						for b.Loop() {
							//b.StopTimer()
							key := insertedKeys[idx%len(insertedKeys)]
							idx++
							//b.StartTimer()

							ht.Get(key)
						}
					})
				}
			}
		}
	}
}

func runDeleteBenchmark(b *testing.B) {
	for method, newHashTable := range Factories {
		for _, size := range Sizes {
			for keyKind, keyGen := range KeyGens {
				for _, loadFactor := range LoadFactors {
					lfString := format(loadFactor)
					testName := fmt.Sprintf("%s-%s-%s-%d", method, keyKind, lfString, size)

					b.Run(testName, func(b *testing.B) {
						ht := newHashTable(size)
						ht.SetLoadFactor(loadFactor)
						keysGen := keyGen(size)

						insertedKeys := make([]int, 0, size)
						for key := range keysGen {
							ht.Insert(key, key)
							insertedKeys = append(insertedKeys, key)
						}

						Random.Shuffle(len(insertedKeys), func(i, j int) {
							insertedKeys[i], insertedKeys[j] = insertedKeys[j], insertedKeys[i]
						})

						var idx int
						for b.Loop() {
							b.StopTimer()
							key := insertedKeys[idx%len(insertedKeys)]
							idx++
							b.StartTimer()

							ht.Delete(key)

							b.StopTimer()
							ht.Insert(key, key)
							b.StartTimer()
						}
					})
				}
			}
		}
	}
}
