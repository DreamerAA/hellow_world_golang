package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"psql"
	"strings"
)

type Item struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Details string `json:"details"`
}

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func getAllItems(db *sql.DB, filters map[string]interface{}) []Item {
	var items []Item
	rows, err := psql.CreateSelectQuery(db, "items", []string{"ID", "Name", "Details"}, filters, "ID DESC")
	if err != nil {
		log.Fatal(err)
		return items
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Details)
		if err == nil {
			items = append(items, item)
		} else {
			log.Fatal(err)
			return items
		}
	}
	return items
}

func putResponseToJson(jresp []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(jresp)
}

func registerRequests(db *sql.DB) {
	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		fmt.Println("Length:", len(parts))
		if len(parts) != 3 || parts[2] == "" {
			fmt.Println("Length:", len(parts))
			http.NotFound(w, r)
			return
		}
		fmt.Println("ID:", parts[2])
		filters := map[string]interface{}{"ID": parts[2]}
		switch r.Method {
		case http.MethodGet:
			items := getAllItems(db, filters)
			jresp, err := json.Marshal(items)
			if err != nil {
				fmt.Println("Ошибка кодирования JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			putResponseToJson(jresp, w)
			w.WriteHeader(http.StatusOK)
			return

		case http.MethodPut:
			jreq := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&jreq)
			if err != nil {
				fmt.Println("Ошибка декодирования JSON:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			values := make(map[string]interface{}, len(jreq))
			for key := range jreq {
				values[key] = jreq[key]
			}

			err = psql.CreateUpdateQuery(db, "items", values, filters)
			if err != nil {
				fmt.Println("Ошибка выполнения запроса для обновления:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodDelete:
			err := psql.CreateDeleteQuery(db, "items", filters)
			if err != nil {
				fmt.Println("Ошибка выполнения запроса для удаления:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusNoContent)
			return
		case http.MethodPost:
			fmt.Println("Обработка такого запроса не предуссмотрена API")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		default:
			fmt.Println("Обработка такого запроса не предуссмотрена API")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		fmt.Println("Length:", len(parts))
		if len(parts) != 2 {
			fmt.Println("Length:", len(parts))
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			items := getAllItems(db, make(map[string]interface{}))

			jresp, err := json.Marshal(items)
			if err != nil {
				fmt.Println("Ошибка декодирования JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			putResponseToJson(jresp, w)
			return
		case http.MethodPost:
			jreq := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&jreq)
			if err != nil {
				fmt.Println("Ошибка декодирования JSON:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var ok bool
			var item Item

			item.Name, ok = jreq["Name"].(string)
			if !ok {
				fmt.Println("Ошибка получения имени элемента")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			item.Details, ok = jreq["Details"].(string)
			if !ok {
				fmt.Println("Ошибка получения деталей элемента")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var id int
			id, err = psql.CreateInsertQuery(db, "items", []string{"Name", "Details"}, []interface{}{item.Name, item.Details})
			if err != nil {
				fmt.Println("Ошибка выполнения запроса для создания:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			print("Element ID:", id)
			jresp := make(map[string]string)
			jresp["status"] = "OK"
			jresp["id"] = fmt.Sprintf("%d", id)
			jresp["message"] = "Item created"
			jrespJson, err := json.Marshal(jresp)
			if err != nil {
				fmt.Println("Ошибка декодирования JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			putResponseToJson(jrespJson, w)
			return
		case http.MethodPut:
			fmt.Println("Обработка такого запроса не предуссмотрена API")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		case http.MethodDelete:
			fmt.Println("Обработка такого запроса не предуссмотрена API")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		default:
			fmt.Println("Обработка такого запроса не предуссмотрена API")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}

func main() {
	fmt.Println("Server started on port 8080")
	db, err := psql.OpenDataBase("usual")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	psql.CreateTable(db, "items", []string{"ID SERIAL PRIMARY KEY", "Name TEXT NOT NULL", "Details TEXT"})

	registerRequests(db)

	psql.CreateEnum(db, "status", []string{"active", "inactive", "undefined"})
	psql.RemoveEnum(db, "status")

	http.ListenAndServe(":8080", nil)
}
