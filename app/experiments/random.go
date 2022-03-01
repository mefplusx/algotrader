package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	const COIN_HEADS = true  // орел
	const COIN_TAILS = false // решка

	road := []bool{}

	const CNT_ATTEMPTS = 100000000

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < CNT_ATTEMPTS; i++ {
		road = append(road, rand.Float32() < 0.5)
	}

	statisticsDirection := make(map[bool]int)
	statisticsDirection[COIN_HEADS] = 0
	statisticsDirection[COIN_TAILS] = 0
	for _, coin := range road {
		if coin == COIN_HEADS {
			statisticsDirection[COIN_HEADS]++
		} else {
			statisticsDirection[COIN_TAILS]++
		}
	}

	fmt.Println("")
	fmt.Println("coins heads :", statisticsDirection[COIN_HEADS])
	fmt.Println("coins tails :", statisticsDirection[COIN_TAILS])

	statisticsCnt := make(map[int]int)
	cnt := 1
	for i, _ := range road {
		if i == 0 {
			continue
		}

		if road[i] != road[i-1] {
			statisticsCnt[cnt]++
			cnt = 1
		} else {
			cnt++
		}
	}

	fmt.Println("")
	for i := 0; i < 100; i++ {
		if _, ok := statisticsCnt[i]; ok {
			fmt.Println(i, ":", statisticsCnt[i])
		}
	}

	// Существуют точки синхронизации для уравнивания статистики???????
}
