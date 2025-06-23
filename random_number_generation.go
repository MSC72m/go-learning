package main

import (
	"fmt"
	"math/rand"
	"time"
)

type randomNumbers struct {
	random  rand.Rand
	numbers []int
}

func (r randomNumbers) getSeed() int {
	return int(time.Now().UnixNano())
}

func (r *randomNumbers) getRand(seed int) {
	r.random = *rand.New(rand.NewSource(int64(seed)))
}

func (r *randomNumbers) appendRandNumber(n int) {
	num := r.random.Intn(n)
	r.numbers = append(r.numbers, num)
	fmt.Printf("Generated %d \n", num)
}
