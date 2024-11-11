package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"psql"

	"github.com/gorilla/mux"
)

type SalesRecord struct {
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Date      string `json:"date"`
	OrderID   int    `json:"order_id"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func extractSalesFromJson(data []byte) ([]SalesRecord, error) {
	sales := make([]SalesRecord, 0)
	err := json.Unmarshal(data, &sales)
	return sales, err
}

func putResponseToJson(jresp []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(jresp)
}

func readDataAndInsertToDB(db *sql.DB, data []byte) ([]int, error) {
	sales, err := extractSalesFromJson(data)
	if err != nil {
		return nil, err
	}
	ids := make([]int, len(sales))
	for i, sale := range sales {
		fmt.Println(sale)
		id, err := psql.CreateInsertQuery(db, "sales", []string{"ID", "ProductID", "Quantity", "Date"}, []interface{}{sale.OrderID, sale.ProductID, sale.Quantity, sale.Date})
		if err != nil {
			fmt.Println("Ошибка для вставки записи:", err, sale)
			continue
		}
		ids[i] = id
	}
	return ids, err
}

func registerRequests(db *sql.DB) {
	router := mux.NewRouter()

	router.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Ошибка чтения тела запроса:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ids, err := readDataAndInsertToDB(db, body)
			if err != nil {
				fmt.Println("Ошибка:", err)
				return
			}

			jresp := map[string]interface{}{"status": "OK", "ids": ids, "message": "Item created"}
			jrespJson, err := json.Marshal(jresp)
			if err != nil {
				fmt.Println("Ошибка декодирования JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			putResponseToJson(jrespJson, w)
		}
	}).Methods("POST")

	http.Handle("/", router)
}

func runScanner(db *sql.DB) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		ok := scanner.Scan()
		if !ok {
			fmt.Println("Error reading input")
			continue
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error during input scanning:", err)

		}
		query := scanner.Text()
		if query == "\\q" {
			break
		}
		bytes, _ := os.ReadFile(query)
		ids, err := readDataAndInsertToDB(db, bytes)
		if err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			fmt.Println("Sales Records were created with ids:", ids)
		}
	}
}

func getAllSales(db *sql.DB) []SalesRecord {
	var sales []SalesRecord
	rows, err := psql.CreateSelectQuery(db, "sales", "*", map[string]interface{}{}, "ID DESC")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return sales
	}
	for rows.Next() {
		var sale SalesRecord
		err := rows.Scan(&sale.OrderID, &sale.ProductID, &sale.Quantity, &sale.Date)
		if err == nil {
			sales = append(sales, sale)
		} else {
			log.Fatal(err)
			return sales
		}
	}
	return sales
}

func getAllProducts(db *sql.DB) []Product {
	var products []Product
	rows, err := psql.CreateSelectQuery(db, "products", "*", map[string]interface{}{}, "ID DESC")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return products
	}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err == nil {
			products = append(products, product)
		} else {
			log.Fatal(err)
			return products
		}
	}
	return products
}

func prepareReport(db *sql.DB) string {
	products := getAllProducts(db)
	sales := getAllSales(db)
	report := "Отчет о продажах\n================\n"
	report += "Товар       Сумма\n------------------\n"
	total := 0.0
	for _, sale := range sales {
		product := products[sale.ProductID-1]
		sum := product.Price * float64(sale.Quantity)
		report += fmt.Sprintf("%d %s %s $%.2f\n", sale.OrderID, sale.Date, product.Name, sum)
		total += sum
	}
	report += "------------------\n"
	report += fmt.Sprintf("Итого:     $%.2f\n", total)
	return report

}

func main() {
	mode := flag.String("mode", "http", "use `json`, `http` or `stat` mods")
	flag.Parse()

	db, err := psql.OpenDataBase("usual")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	psql.CreateTable(db, "sales", []string{"ID INT UNIQUE NOT NULL", "ProductID INT NOT NULL", "Quantity TEXT NOT NULL", "Date DATE NOT NULL"})

	if *mode != "json" && *mode != "http" && *mode != "stat" {
		log.Fatal("mode must be `json`, `http` or `stat`")
		return
	}
	if *mode == "json" {
		fmt.Println("Enter path to json file:")
		runScanner(db)
	} else if *mode == "http" {
		registerRequests(db)
		fmt.Println("Server started on port 8080")
		// read http
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else if *mode == "stat" {
		fmt.Println("Print Stats starts...")
		report := prepareReport(db)
		fmt.Println(report)
	}

}
