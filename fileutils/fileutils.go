package fileutils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func TryOpenFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка при открытии фаила", err)
		return nil
	}
	return file
}

func GenerateRandomText(count int) string {
	result := make([]rune, count)
	for i := 0; i < count; i++ {
		random_int := 32 + rand.Intn(126-32)
		result[i] = rune(random_int)
	}
	return string(result)
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
	for i := 0; i < 10; i++ {
		result := GenerateRandomText(100)
		_, err := writer.WriteString(result + "\n")
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
