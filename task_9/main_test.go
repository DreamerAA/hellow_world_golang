package main

import (
	"fileutils"
	"fmt"
	"math"
	"math/rand"
	"testing"
	// "./main"
)

func GenerateRandomFloats(count int) []float64 {
	max_int := int(math.MaxUint >> 1)
	result := make([]float64, count)
	for i := 0; i < count; i++ {
		result[i] = rand.Float64() * float64(max_int)
	}
	return result
}

func TestSumofAllValues(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0}
	expected := 6.0
	result := SumAllValues(values)
	if result != expected {
		t.Fatalf("Expected %f, got %f", expected, result)
	}

	values = []float64{}
	expected = 0.
	result = SumAllValues(values)
	if result != expected {
		t.Fatalf("Expected %f, got %f", expected, result)
	}
}

func checkSorting(values []float64) error {
	for i := 1; i < len(values); i++ {
		if values[i-1] > values[i] {
			return fmt.Errorf("Expected sorted values: %f ,  %f", values[i-1], values[i])
		}
	}
	return nil
}

func spaceChecker(new_text string, old_text string, count_spaces int) error {
	if len(new_text) != len(old_text)+count_spaces {
		return fmt.Errorf("Expected %d, got %d", len(old_text)+count_spaces, len(new_text))
	}
	return nil
}

func TestRandomSpaceInserter(t *testing.T) {
	count_spaces := 10
	text_size := 100
	text := fileutils.GenerateRandomText(text_size)
	new_text := RandomSpaceInserter(text, count_spaces)
	err := spaceChecker(new_text, text, count_spaces)
	if err != nil {
		t.Fatal(err)
	}

}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func insertToEmptyWithPanic() {
	count_spaces := 10
	var text2 = ""
	new_text := RandomSpaceInserter(text2, count_spaces)
	err := spaceChecker(new_text, text2, count_spaces)
	if err != nil {
		fmt.Println(err)
	}
}

func insertMaxIntWithPanic() {
	count_spaces := 10
	text_size := int(math.MaxUint >> 1)
	text2 := fileutils.GenerateRandomText(text_size)
	new_text := RandomSpaceInserter(text2, count_spaces)
	err := spaceChecker(new_text, text2, count_spaces)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRandomSpaceInserterWithPanic(t *testing.T) {
	assertPanic(t, insertToEmptyWithPanic)
	assertPanic(t, insertMaxIntWithPanic)
}

func TestSorting(t *testing.T) {
	values := GenerateRandomFloats(10000)
	SortArray(values)
	err := checkSorting(values)
	if err != nil {
		t.Fatal(err)
	}

	values = GenerateRandomFloats(0)
	SortArray(values)
	err = checkSorting(values)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSorting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		values := GenerateRandomFloats(100000)
		SortArray(values)
	}
}

func BenchmarkInsertSpaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		count_spaces := 1000
		text_size := 1000
		text2 := fileutils.GenerateRandomText(text_size)
		RandomSpaceInserter(text2, count_spaces)
	}
}

func BenchmarkSumAllValues(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandomFloats(10000)
	}
}
