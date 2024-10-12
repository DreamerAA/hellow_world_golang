package main

import (
	"bufio"
	"fileutils"
	"fmt"
	"os"
	"strconv"
	"sync"
)

func createFolderWithRandomTxtFiles(folder_path string, count_files int) {
	exist, err := fileutils.Exists(folder_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	if exist {
		os.RemoveAll(folder_path)
	}
	os.Mkdir(folder_path, 0777)
	for i := 0; i < count_files; i++ {
		file_path := folder_path + strconv.Itoa(i) + ".txt"
		file, err := os.Create(file_path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fileutils.WriteTofileRandomText(file)
	}
}

func clearOutputFolder(output_dir string) {
	exist, err := fileutils.Exists(output_dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	if exist {
		os.RemoveAll(output_dir)
	}
	os.Mkdir(output_dir, 0777)
}

type ReadWorker struct {
	file_name string
	chanel    chan string
}

type WriteWorker struct {
	full_file_name string
	chanel         chan string
}

func readFileUseWorker(worker ReadWorker, wg *sync.WaitGroup) {
	defer wg.Done()
	file := fileutils.TryOpenFile(worker.file_name)
	scanner := bufio.NewScanner(file)
	var full_text string = ""
	for scanner.Scan() {
		text := scanner.Text()
		full_text += text
		worker.chanel <- text
	}
}

func writeTextToFile(text string, file *os.File) {
	if _, err := file.WriteString(text + "\n"); err != nil {
		fmt.Println(err)
	}
}

func workOfWritingText(worker WriteWorker, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.OpenFile(worker.full_file_name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for result := range worker.chanel {
		writeTextToFile(result, file)
	}
}

func closeChannels(write_workers []WriteWorker) {
	for _, worker := range write_workers {
		close(worker.chanel)
	}
}

func readWriteManager(input_dir string, output_dir string, file_names []string, wg_main *sync.WaitGroup) {
	defer wg_main.Done()
	read_workers := []ReadWorker{}
	write_workers := []WriteWorker{}
	// create workers
	for _, file_name := range file_names {
		chanel := make(chan string)
		read_workers = append(read_workers, ReadWorker{input_dir + file_name, chanel})
		write_workers = append(write_workers, WriteWorker{output_dir + file_name, chanel})
	}

	// run workers
	var wg_read sync.WaitGroup
	for _, worker := range read_workers {
		wg_read.Add(1)
		go readFileUseWorker(worker, &wg_read)
	}
	var wg_write sync.WaitGroup
	for _, worker := range write_workers {
		wg_write.Add(1)
		go workOfWritingText(worker, &wg_write)
	}
	go func() {
		wg_read.Wait() // Ожидание завершения всех воркеров чтения
		closeChannels(write_workers)
	}()

	wg_write.Wait() // Ожидание завершения записи
}

func main() {
	input_dir := "./input/"
	output_dir := "./output/"
	count_files := 3
	createFolderWithRandomTxtFiles(input_dir, count_files)
	clearOutputFolder(output_dir)

	var file_names []string
	for i := 0; i < count_files; i++ {
		file_names = append(file_names, strconv.Itoa(i)+".txt")
	}
	var wg_main sync.WaitGroup
	wg_main.Add(1)
	go readWriteManager(input_dir, output_dir, file_names, &wg_main)
	wg_main.Wait()
	fmt.Println("Основная программа завершена")
}
