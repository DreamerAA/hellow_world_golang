package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
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

func readFromUrlandMoveToChanel(worker UrlWorker, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		result := readDataFromUrl(worker)
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
				}
				file_name := "response.json"
				os.WriteFile(file_name, []byte(text), 0644)

				rates := result["rates"].(map[string]interface{})
				rub := rates["RUB"].(float64)
				return fmt.Sprint(rub)
			},
		},
		{
			"Programming joke",
			"https://official-joke-api.appspot.com/jokes/programming/random",
			make(chan string),
			func(text string) string {
				if len(text) == 0 {
					return text
				}
				if text[0] == '[' {
					text = text[1:(len(text) - 1)]
				}
				result := stringToJsonObject(text)
				if result == nil {
					fmt.Println("Error reading data from url")
				}
				setup := result["setup"].(string)
				punchline := result["punchline"].(string)
				return setup + " " + punchline
			},
		},
	}
	var wg sync.WaitGroup
	for _, worker := range workers {
		wg.Add(1)
		go readFromUrlandMoveToChanel(worker, &wg)
	}
	for _, worker := range workers {
		go func(w UrlWorker) {
			for result := range w.chanel {
				fmt.Println(w.description + " " + result)
			}
		}(worker)
	}
	// Ждём завершения всех горутин
	wg.Wait()
	fmt.Println("Все горутины завершены")

}
