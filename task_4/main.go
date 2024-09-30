package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func readDataFromUrl(worker UrlWorker) (string, error) {
	resp, err := http.Get(worker.url)
	if err != nil {
		return "", fmt.Errorf("error reading data from url: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body of data: %w", err)
	}
	return worker.operation(string(body)), nil
}

func readFromUrlandMoveToChanel(worker UrlWorker, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		result, err := readDataFromUrl(worker)
		if err != nil {
			fmt.Println(err)
			return
		}
		worker.chanel <- result
		fmt.Println(i)
	}
	close(worker.chanel)
}

type UrlWorker struct {
	description string
	url         string
	chanel      chan string
	operation   func(url string) string
}

func stringToJsonObject(text string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		fmt.Println("Error converting string to json", err)
		return nil
	}
	return result
}

func stringToJsonArray(text string) []map[string]interface{} {
	var result []map[string]interface{}
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		fmt.Println("Error converting string to json", err)
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
				result := stringToJsonObject(text)
				if result == nil {
					fmt.Println("Error reading data from url")
					return ""
				}

				rates, ok := result["rates"].(map[string]interface{})
				if !ok {
					fmt.Println("Error: 'rates' not found or has wrong format")
					return ""
				}
				rub, ok := rates["RUB"].(float64)
				if !ok {
					fmt.Println("Error: 'RUB' not found or has wrong format")
					return ""
				}
				return fmt.Sprint(rub)
			},
		},
		{
			"Programming joke",
			"https://official-joke-api.appspot.com/jokes/programming/random",
			make(chan string),
			func(text string) string {
				result := stringToJsonArray(text)
				if result == nil {
					fmt.Println("Error reading data from url")
					return ""
				}
				if len(result) != 1 {
					fmt.Println("Error: expected 1 joke")
					return ""
				}
				json_joke := result[0]
				setup, ok := json_joke["setup"].(string)
				if !ok {
					fmt.Println("Error: 'setup' not found or has wrong format")
					return ""
				}
				punchline, ok := json_joke["punchline"].(string)
				if !ok {
					fmt.Println("Error: 'punchline' not found or has wrong format")
					return ""
				}
				return setup + " " + punchline
			},
		},
	}
	var wg_get sync.WaitGroup
	for _, worker := range workers {
		wg_get.Add(1)
		go readFromUrlandMoveToChanel(worker, &wg_get)
	}

	for _, worker := range workers {
		go func(w UrlWorker) {
			for result := range w.chanel {
				fmt.Println(w.description + " " + result)
			}
		}(worker)
	}
	// Ждём завершения всех горутин
	wg_get.Wait()
	fmt.Println("Все горутины завершены")

}
