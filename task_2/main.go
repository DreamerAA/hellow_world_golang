package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

// generateRandomText generates a random string of length count consisting of
// characters from the ASCII range of 32 to 126 (inclusive). The string is
// returned as a string. The function uses the rand package to generate the
// random numbers.
func generateRandomText(count int) string {
	result := ""
	for i := 0; i < count; i++ {
		random_int := 32 + rand.Intn(126-32)
		char_as_str := string(random_int)
		result += char_as_str
	}
	return result
}

func writeTofileRandomText(file *os.File) {
	writer := bufio.NewWriter(file)
	for i := 0; i < 10; i++ {
		result := generateRandomText(100)
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

func createNewRandomFile(pathToFile string) {

	file, err := os.Create(pathToFile)
	if err != nil {
		fmt.Println("Ошибка при создании фаила", err)
		return
	}
	defer file.Close()

	writeTofileRandomText(file)
}

func tryOpenFile(pathToFile string) *os.File {
	file, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println("Ошибка при открытии фаила", err)
		return nil
	}
	return file
}

func main() {
	file_name := "input.txt"
	file := tryOpenFile(file_name)

	for file == nil {
		createNewRandomFile(file_name)
		file = tryOpenFile(file_name)
	}
	defer file.Close()

	file_data := make([]string, 0)
	scanner := bufio.NewScanner(file)
	fmt.Println("Выводим на экран")
	for scanner.Scan() {
		text := scanner.Text()
		file_data = append(file_data, text)
		fmt.Println(text)
		fmt.Println("Продолжаем")
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении фаила", err)
		return
	}
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка при создании фаила", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	bytesWritten, err := writer.WriteString("Это строка записана в файл.\n")
	if err != nil {
		fmt.Println("Ошибка при записи в файл", err)
		return
	}
	for i := 0; i < len(file_data); i++ {
		new_bytesWritten, err := writer.WriteString(file_data[i] + "\n")
		bytesWritten += new_bytesWritten
		if err != nil {
			fmt.Println("Error writing to file", err)
			return
		}
	}

	fmt.Printf("Записано байт: %d\n", bytesWritten)

	err = writer.Flush()
	if err != nil {
		fmt.Println("Ошибка при сбросе буфера в файл", err)
		return
	}

}
