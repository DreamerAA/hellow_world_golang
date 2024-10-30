package fileutils

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
)

func TryAllocateSlice(size int) (data []rune, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Паника при выделении памяти:", r)
			err = fmt.Errorf("Невозможно выделить память %d", size)
		}
	}()
	if size < 0 || size > int(math.MaxInt32) {
		return nil, fmt.Errorf("Невозможно выделить память %d", size)
	}
	// Пытаемся выделить слайс нужного размера
	data = make([]rune, size)

	return data, nil
}

func TryOpenFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка при открытии фаила", err)
		return nil
	}
	return file
}

func GenerateRandomText(count int, seed int64) (string, error) {
	rnd := rand.New(rand.NewSource(seed))
	result, err := TryAllocateSlice(count)
	if err != nil {
		return "", err
	}
	for i := 0; i < count; i++ {
		random_int := 32 + rnd.Intn(126-32)
		result[i] = rune(random_int)
	}
	return string(result), nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriteTofileRandomText(file *os.File) {
	writer := bufio.NewWriter(file)
	var seed int64 = 10
	for i := 0; i < 10; i++ {
		result, err := GenerateRandomText(100, seed)
		if err != nil {
			fmt.Println("Невозможно выделить память")
			return
		}
		_, err = writer.WriteString(result + "\n")
		if err != nil {
			fmt.Println("Error writing to file", err)
			return
		}
	}
	err := writer.Flush()
	if err != nil {
		fmt.Println("Ошибка при сбросе буфера в файл", err)
		return
	}
}
