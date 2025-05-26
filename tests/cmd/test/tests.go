package test

import (
	"encoding/csv"
	"log"
	"math/bits"
	"os"
	"path/filepath"
	"strconv"
)

const OutputDir = "results"

func RunCollisionsAndProbesTest() {
	for method := range Factories {
		for keyKind := range KeyGens {
			for _, loadFactor := range LoadFactors {
				CollisionsCountTest(method, keyKind, loadFactor)
			}

			ProbesCountTest(method, keyKind)
		}
	}
}

func CollisionsCountTest(method string, keyKind string, loadFactor float64) {
	var (
		collisionsMetrics [][]string
	)

	for _, size := range Sizes {
		ht := Factories[method](size)
		ht.SetLoadFactor(loadFactor)
		keysGen := KeyGens[keyKind](size)

		ht.ResetCollisions()

		for key := range keysGen {
			ht.Insert(key, key)
		}

		collisionsMetrics = append(collisionsMetrics, getRecord(size, ht.Collisions()))
	}

	lfString := format(loadFactor)
	saveMetrics(filepath.Join(OutputDir, "Collision", method, lfString), keyKind, collisionsMetrics)
}

func ProbesCountTest(method string, keyKind string) {
	var (
		size          = 5000
		samples       = 1_000
		probesMetrics [][]string
	)

	loadFactors := []float64{0.5, 0.65, 0.75, 0.9}

	for _, loadFactor := range loadFactors {
		ht := Factories[method](size)
		ht.SetLoadFactor(1.0)

		desiredInsertions := int(loadFactor * float64(nextPowerOfTwo(size)))
		keysGen := KeyGens[keyKind](desiredInsertions)

		insertedKeys := make([]int, 0, desiredInsertions)
		for key := range keysGen {
			ht.Insert(key, key)
			insertedKeys = append(insertedKeys, key)
		}

		ht.ResetProbes()

		Random.Shuffle(len(insertedKeys), func(i, j int) {
			insertedKeys[i], insertedKeys[j] = insertedKeys[j], insertedKeys[i]
		})

		for i := range samples {
			key := insertedKeys[i%len(insertedKeys)]
			ht.Get(key)
		}

		probesMetrics = append(probesMetrics, getRecord(loadFactor, float64(ht.Probes())/float64(samples)))
	}

	saveMetrics(filepath.Join(OutputDir, "Probes", method), keyKind, probesMetrics)
}

func saveMetrics(dir, keyKind string, metrics [][]string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("failed to create directory %s: %v", dir, err)
	}

	filePath := filepath.Join(dir, keyKind+".csv")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, metric := range metrics {
		if err = writer.Write(metric); err != nil {
			log.Fatalf("failed to write CSV: %v", err)
		}
	}
}

func getRecord(metrics ...any) []string {
	var record []string

	for _, metric := range metrics {
		record = append(record, format(metric))
	}

	return record
}

func format(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	default:
		log.Fatalf("unsupported metric type: %T", v)
		return ""
	}
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
