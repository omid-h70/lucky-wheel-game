package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type user struct {
	uuid   string
	tryNum int
}

var (
	ErrOutOfBondNumber = errors.New("Rand Num or Seeds are not valid")
	ErrNotValidSeeds   = errors.New("Seeds are not valid")
)

//var seeds []float32 = []float32{A, B, C, D, E}

func GetRandomNumberByTime() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return r1.Intn(100)
}

func GetPrizeV1(rand int, seeds []float32) (int, error) {
	var sum int
	for i := 0; i < len(seeds); i++ {
		sum += int(seeds[i] * 100)
		if rand < sum {
			return i, nil
		}
	}
	return 0, ErrOutOfBondNumber
}

func GetPrizeV2(rand int, vseeds map[string]string) (string, error) {

	fmt.Printf("Rand %d Seeds %v \n", rand, vseeds)

	var sum int
	for i, seed := range vseeds {
		floatNum, err := strconv.ParseFloat(seed, 32)
		if err != nil {
			return "0", ErrNotValidSeeds
		}

		sum += int(floatNum * 100)
		if rand < sum {
			return i, nil
		}
	}
	return "0", ErrOutOfBondNumber
}

func SortMapBasedOnValue(m map[string]string) map[string]string {

	// Extract keys from map
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		//return m[keys[i]] < m[keys[j]]
		f1, _ := strconv.ParseFloat(m[keys[i]], 32)
		f2, _ := strconv.ParseFloat(m[keys[j]], 32)
		return f1 < f2
	})

	// Sort keys
	sortedMap := make(map[string]string, len(m))
	for i := 0; i < len(keys); i++ {
		sortedMap[keys[i]] = m[keys[i]]
	}
	//fmt.Println(sortedMap)
	return sortedMap
}
