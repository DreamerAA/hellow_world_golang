package main

import (
	"math/rand"
	"sort"
	"strings"
)

func SumAllValues(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum
}

func RandomSpaceInserter(text string, count int) string {
	cur_text := strings.Clone(text)
	for i := 0; i < count; i++ {
		index := rand.Intn(len(cur_text))
		cur_text = cur_text[:index] + " " + cur_text[index:]
	}
	return cur_text
}

func SortArray(array []float64) []float64 {
	sort.Float64s(array)
	return array
}
