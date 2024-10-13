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

func createSuccessResponse(message string) Response {
	return Response{Message: message, Status: 200}
}

func putResponseToJson(response Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func setNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func registerGetRequests() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		response := createSuccessResponse("Hello, this is JSON response!")
		putResponseToJson(response, w)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setNotFoundResponse(w, r)
	})
}
func registerPostRequests() {
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
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
			response := createSuccessResponse(fmt.Sprintf("Hello, %s!", message))
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
