package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func getFileNames(path string) []string {
	dir, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}
	var file_names []string
	for _, file := range dir {
		file_names = append(file_names, file.Name())
	}
	return file_names
}

func main() {

	list := flag.String("list", "", "directory path")
	convert := flag.String("convert", "", "convert string")

	flag.CommandLine.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  -list <dir>    : List files in a directory")
		fmt.Println("  -convert <text>: Convert text to uppercase")
		os.Exit(1) // Завершение программы с кодом ошибки
	}
	// flag.Parse()
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error:", err)
		flag.CommandLine.Usage()
		return
	}

	has_list := *list != ""
	has_convert := *convert != ""

	if !has_list && !has_convert {
		flag.CommandLine.Usage()
	}

	var file_names []string
	if has_list {
		file_names = getFileNames(*list)
		if len(file_names) == 0 {
			fmt.Println("No files or text to process")
			os.Exit(1)
		}
		for _, file_name := range file_names {
			fmt.Println(file_name)
		}
	}
	if has_convert {
		fmt.Println(strings.ToUpper(*convert))
	}
}
