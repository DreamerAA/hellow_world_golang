package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func putResponseToJson(response Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func setNotFoundResponse(w http.ResponseWriter) {
	response := Response{Message: "Not Found", Status: 404}
	putResponseToJson(response, w)
}

func registerGetRequests() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		response := Response{Message: "Hello, this is JSON response!", Status: 200}
		putResponseToJson(response, w)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setNotFoundResponse(w)
	})
}
func registerPostRequests() {
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			setNotFoundResponse(w)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if message, ok := data["message"].(string); ok {
			response := Response{Message: fmt.Sprintf("Hello, %s!", message), Status: 200}
			putResponseToJson(response, w)
		} else {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	})
}

func main() {
	registerGetRequests()
	registerPostRequests()
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
