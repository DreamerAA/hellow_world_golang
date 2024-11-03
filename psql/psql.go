package psql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

func addParams(header string, new_params map[string]interface{}, splitter string, istart int) (string, []interface{}) {
	query := header
	ind := istart + 1
	query_args := []interface{}{}
	for k, v := range new_params {
		query += fmt.Sprintf("%s ", k) + splitter + fmt.Sprintf("$%d", ind)
		query_args = append(query_args, v)
		ind++
		if ind <= len(new_params) {
			query += ", "
		}
	}
	return query, query_args
}

func OpenDataBase(name string) (*sql.DB, error) {

	connStr := "postgres://admin:admin@localhost:5432/" + name + "?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Connection to PostgreSQL was established successfully")
	return db, nil
}

func CreateEnum(db *sql.DB, enum_name string, statuses []string) {
	query_enum := "CREATE TYPE " + enum_name + " AS ENUM ("
	for i, status := range statuses {
		query_enum += "'" + status + "'"
		if i != len(statuses)-1 {
			query_enum += ", "
		}
	}
	query_enum += ");"

	fmt.Println("Query create enum:", query_enum)
	_, err := db.Exec(query_enum)
	if err != nil {
		if pg_err, ok := err.(*pq.Error); ok {
			if pg_err.Code == pq.ErrorCode("42710") {
				fmt.Println("enum уже существует.")
			} else {
				fmt.Println("Ошибка создания enum:", err, pg_err.Message, pg_err.Code)
			}
		}
	} else {
		fmt.Println("enum успешно создан!")
	}
}
func RemoveEnum(db *sql.DB, enum_name string) {
	query := "DROP TYPE IF EXISTS " + enum_name + ";"

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Ошибка удаления enum:", err)
	} else {
		fmt.Println("enum успешно удален!")
	}
}

func CreateTable(db *sql.DB, table_name string, fields []string) {
	query := "CREATE TABLE IF NOT EXISTS " + table_name + "("
	for i, field := range fields {
		query += field
		if i < len(fields)-1 {
			query += ", "
		} else {
			query += ");"
		}
	}

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
	} else {
		fmt.Println("Таблица успешно создана или уже существует!")
	}
}

func CreateSelectQuery(db *sql.DB, table_name string, fields []string, filters map[string]interface{}, otext string) (*sql.Rows, error) {
	queryArgs := []interface{}{}
	query := "SELECT " + strings.Join(fields, ", ") + " FROM " + table_name

	if len(filters) > 0 {
		// Создание запроса с фильтрами и сортировкой
		query, queryArgs = addParams(query+" WHERE ", filters, " = ", 0)
	}
	if otext != "" {
		query += " ORDER BY " + otext
	}
	query += ";"
	fmt.Println(query)
	fmt.Println(queryArgs)
	return db.Query(query, queryArgs...)
}

func CreateInsertQuery(db *sql.DB, table_name string, fields []string, values []interface{}) (int, error) {
	query := "INSERT INTO " + table_name + "(" + strings.Join(fields, ", ") + ") VALUES ("
	for i := 0; i < len(fields); i++ {
		query += "$" + fmt.Sprintf("%d", i+1)
		if i < len(fields)-1 {
			query += ", "
		}
	}

	query += ") RETURNING id;"
	var id int
	err := db.QueryRow(query, values...).Scan(&id)

	if err != nil {
		return -1, err
	}
	return id, nil
}

func CreateDeleteQuery(db *sql.DB, table_name string, filters map[string]interface{}) error {
	query, args := addParams("DELETE FROM "+table_name+" WHERE ", filters, " = ", 0)
	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func CreateUpdateQuery(db *sql.DB, table_name string, values map[string]interface{}, filters map[string]interface{}) error {
	query, args := addParams("UPDATE "+table_name+" SET ", values, " = ", 0)
	query, args_where := addParams(query+" WHERE ", filters, " = ", len(values))
	args = append(args, args_where...)
	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
