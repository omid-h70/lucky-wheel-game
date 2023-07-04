package domain

import (
	"fmt"
	"testing"
	"time"
)

const (
	AA = 0.1
	BB = 0.3
	CC = 0.2
	DD = 0.15
	EE = 0.25
)

var testSeeds []float32 = []float32{AA, BB, CC, DD, EE}

func Test_should_sort_map_based_on_values(t *testing.T) {
	m := map[string]string{
		"A": "1.0",
		"B": "2.0",
		"C": "0.12",
	}

	mm := SortMapBasedOnValue(m)
	fmt.Println(mm)
}

func Test_transfer_should_return_prize_based_on_rand_num(t *testing.T) {
	rand := GetRandomNumberByTime()
	prize, err := GetPrizeV1(rand, testSeeds)

	if err != nil {
		t.Error("Test Failed")
	}
	fmt.Printf("Rand %d Prize %d", rand, prize)
}

func Test_transfer_should_return_prize_based_on_multiple_rand_num(t *testing.T) {

	for i := 0; i < 100; i++ {
		rand := GetRandomNumberByTime()
		prize, err := GetPrizeV1(rand, testSeeds)

		if err != nil {
			t.Error("Test Failed")
		}
		fmt.Printf("Rand %d Prize %d \n", rand, prize)
		time.Sleep(500 * time.Millisecond)
	}
}
