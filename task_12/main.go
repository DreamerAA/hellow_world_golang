package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"psql"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	router.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		filters := map[string]interface{}{"ID": id}

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
		case http.MethodDelete:
			err := psql.CreateDeleteQuery(db, "items", filters)
			if err != nil {
				fmt.Println("Ошибка выполнения запроса для удаления:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT", "DELETE")

	router.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
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
		case http.MethodPost:
			jreq := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&jreq)
			if err != nil {
				fmt.Println("Ошибка декодирования JSON:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			name, ok := jreq["Name"].(string)
			if !ok {
				fmt.Println("Ошибка получения имени элемента")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			details, ok := jreq["Details"].(string)
			if !ok {
				fmt.Println("Ошибка получения деталей элемента")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var id int
			id, err = psql.CreateInsertQuery(db, "items", []string{"Name", "Details"}, []interface{}{name, details})
			if err != nil {
				fmt.Println("Ошибка выполнения запроса для создания:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			jresp := map[string]string{"status": "OK", "id": fmt.Sprintf("%d", id), "message": "Item created"}
			jrespJson, err := json.Marshal(jresp)
			if err != nil {
				fmt.Println("Ошибка декодирования JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			putResponseToJson(jrespJson, w)
		}
	}).Methods("GET", "POST")

	http.Handle("/", router)
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
