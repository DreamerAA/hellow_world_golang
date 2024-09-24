package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func readDataFromUrl(worker UrlWorker) string {
	resp, err := http.Get(worker.url)
	if err != nil {
		fmt.Println("Error reading data from url", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body of data", err)
	}
	return worker.operation(string(body))
}

func readFromUrlandMoveToChanel(worker UrlWorker) {
	for i := 0; i < 10; i++ {
		result := readDataFromUrl(worker)
		worker.chanel <- result
	}
	close(worker.chanel)
}

type UrlWorker struct {
	description string
	url         string
	chanel      chan string
	operation   func(url string) string
}

func stringToJson(text string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		fmt.Println("Error reading data from url", err)
		return nil
	}
	return result
}

func main() {
	workers := []UrlWorker{
		{
			"Its currency from USD to RUB",
			"https://api.exchangerate-api.com/v4/latest/USD",
			make(chan string),
			func(text string) string {
				result := stringToJson(text)
				if result == nil {
					fmt.Println("Error reading data from url")
				}
				return result["rates"].(map[string]interface{})["RUB"].(string)
			},
		},
		{
			"Programming joke",
			"https://official-joke-api.appspot.com/jokes/programming/random",
			make(chan string),
			func(text string) string {
				result := stringToJson(text)
				if result == nil {
					fmt.Println("Error reading data from url")
				}
				setup := result["setup"].(string)
				punchline := result["punchline"].(string)
				return setup + " " + punchline
			},
		},
	}
	for _, worker := range workers {
		go readFromUrlandMoveToChanel(worker)

	}
	select {
	case result, ok := <-workers[0].chanel:
		if !ok {
			fmt.Println("channel was closed")
		} else {
			fmt.Println(workers[0].description + " " + result)
		}

	case result, ok := <-workers[1].chanel:
		if !ok {
			fmt.Println("channel was closed")
		} else {
			fmt.Println(workers[1].description + " " + result)
		}
	}
}
