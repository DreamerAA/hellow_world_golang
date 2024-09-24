package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Skill struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type Person struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Skills []Skill `json:"skills"`
}

func main() {
	json_file, err := os.Open(".\\task_3\\data.json")
	if err != nil {
		fmt.Println("Error opening file", err)
	}
	defer json_file.Close()

	byte_value, _ := io.ReadAll(json_file)
	var person Person
	json.Unmarshal(byte_value, &person)

	fmt.Println(person)
}
