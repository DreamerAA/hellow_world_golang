package main

import (
	"crypto/rand"
	"fmt"
	"sort"

	"fileutils"
)

type UnicodeRange struct {
	start, end int
}

func (urange *UnicodeRange) diff() int {
	return urange.end - urange.start
}

// Выбор случайного числа в диапазоне
func randomInt(max int) (int, error) {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	return int(b[0]) % max, nil
}

func GeneratePassword(length int, useNumbers bool, useSymbols bool, useUppercase bool, useLowercase bool) (string, error) {
	result, err := fileutils.TryAllocateSlice(length)
	if err != nil {
		return "", err
	}

	ranges := []UnicodeRange{}
	if useSymbols {
		ranges = append(ranges, UnicodeRange{33, 48})
		ranges = append(ranges, UnicodeRange{58, 65})
		ranges = append(ranges, UnicodeRange{91, 97})
		ranges = append(ranges, UnicodeRange{123, 127})
	}
	if useNumbers {
		ranges = append(ranges, UnicodeRange{48, 58})
	}
	if useUppercase {
		ranges = append(ranges, UnicodeRange{65, 91})
	}
	if useLowercase {
		ranges = append(ranges, UnicodeRange{97, 123})
	}
	if len(ranges) == 0 {
		return "", fmt.Errorf("no ranges specified")
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	for i := 0; i < length; i++ {
		val, err := randomInt(len(ranges))
		if err != nil {
			return "Cant generate random value", err
		}
		urange := ranges[val]
		val, err = randomInt(urange.diff())
		if err != nil {
			return "Cant generate random value", err
		}
		result[i] = rune(urange.start + val)
	}
	return string(result), nil
}
