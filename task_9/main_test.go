package main

import (
	"fileutils"
	"fmt"
	"math"
	"math/rand"
	"testing"
	// "./main"
)

func GenerateRandomFloats(count int, seed int64) []float64 {
	rnd := rand.New(rand.NewSource(seed))
	max_int := int(math.MaxUint >> 1)
	result := make([]float64, count)
	for i := 0; i < count; i++ {
		result[i] = rnd.Float64() * float64(max_int)
	}
	return result
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

func defaultSpaceInserter(t *testing.T, count_spaces int, text_size int) {
	var seed int64 = 10
	text, err := fileutils.GenerateRandomText(text_size, seed)
	if err != nil {
		t.Fatal(err)
	}
	new_text, err := RandomSpaceInserter(text, count_spaces)
	if err != nil {
		t.Fatal(err)
	}
	err = spaceChecker(new_text, text, count_spaces)
	if err != nil {
		t.Fatal(err)
	}
}

func errorSpaceInserter(t *testing.T, count_spaces int, text_size int) {
	var seed int64 = 10
	text, err := fileutils.GenerateRandomText(text_size, seed)
	if err != nil {
		fmt.Println(err)
		return
	}
	new_text, err := RandomSpaceInserter(text, count_spaces)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Fatal(fmt.Errorf("Expected error, got %s", new_text))
}

func TestRandomSpaceInserter(t *testing.T) {
	defaultSpaceInserter(t, 10, 100)
	defaultSpaceInserter(t, 10, 0)
	errorSpaceInserter(t, 10, int(math.MaxUint>>1))
	errorSpaceInserter(t, int(math.MaxUint>>1), 100)
}

func TestSorting(t *testing.T) {
	var seed int64 = 10
	values := GenerateRandomFloats(10000, seed)
	SortArray(values)
	err := checkSorting(values)
	if err != nil {
		t.Fatal(err)
	}

	values = GenerateRandomFloats(0, seed)
	SortArray(values)
	err = checkSorting(values)
	if err != nil {
		t.Fatal(err)
	}
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

func BenchmarkSorting(b *testing.B) {
	var seed int64 = 10
	for i := 0; i < b.N; i++ {
		values := GenerateRandomFloats(100000, seed)
		SortArray(values)
		checkSorting(values)
	}
}

func BenchmarkInsertSpaces(b *testing.B) {
	var seed int64 = 10
	for i := 0; i < b.N; i++ {
		count_spaces := 1000
		text_size := 1000
		text, err := fileutils.GenerateRandomText(text_size, seed)
		if err != nil {
			b.Fatal(err)
		}
		RandomSpaceInserter(text, count_spaces)
	}
}

func BenchmarkSumAllValues(b *testing.B) {
	var seed int64 = 10
	for i := 0; i < b.N; i++ {
		values := GenerateRandomFloats(10000, seed)
		SumAllValues(values)
	}
}
