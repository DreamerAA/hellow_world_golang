package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

type Command struct {
	name      string
	operation func(string) (string, error)
}

func createTable(db *sql.DB, table_name string) {
	query :=
		"CREATE TABLE IF NOT EXISTS " + table_name +
			"(id SERIAL PRIMARY KEY, title TEXT NOT NULL, description TEXT, status BOOLEAN DEFAULT FALSE);"

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
	} else {
		fmt.Println("Таблица успешно создана или уже существует!")
	}
}

func main() {
	// Строка подключения к PostgreSQL
	connStr := "postgres://admin:admin@localhost:5432/todobd?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	fmt.Println("Connection to PostgreSQL was established successfully")

	createTable(db, "tasks")

	commands := map[string]Command{
		"add": {
			"Add task",
			func(text string) (string, error) {
				datas := strings.Split(text, " ")
				if len(datas) != 2 {
					return "", fmt.Errorf("Invalid input")
				}

				title, description := datas[0], datas[1]
				query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, false) RETURNING id`
				var id int
				err := db.QueryRow(query, title, description).Scan(&id)

				if err != nil {
					return "", err
				}
				return "Tast created with id = " + strconv.Itoa(id), nil
			},
		},
		"rm": {
			"Insert into table",
			func(text string) (string, error) {
				id, err := strconv.Atoi(text)
				if err != nil {
					fmt.Printf("%q does not looks like a number.\n", text)
					return "", fmt.Errorf("Invalid input")
				}
				query := `DELETE FROM tasks WHERE id = $1`
				_, err = db.Exec(query, id)
				if err != nil {
					return "", err
				}
				return "", nil
			},
		},
	}

	// 	Таблица tasks:
	// id: уникальный идентификатор задачи (целое число).
	// title: заголовок задачи (строка).
	// description: описание задачи (строка).
	// status: статус выполнения задачи (булево значение, например, true — выполнена, false — невыполнена).

	// Добавление задачи: Указать заголовок и описание задачи.
	// Удаление задачи: Удалить задачу по её id.
	// Изменение задачи: Изменить статус задачи или её описание.
	// Просмотр всех задач: Показать все задачи, отсортированные по статусу выполнения.

	// value, _ := commands["add"]
	// value.operation("task_1 description_1")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter query: ")
		scanner.Scan()
		query := scanner.Text()
		if query == "\\q" {
			break
		}
		splits := strings.Split(query, " ")
		command := splits[0]
		value, exists := commands[command]
		if !exists {
			fmt.Println("Invalid command")
			continue
		}
		message, err := value.operation(strings.Join(splits[1:], " "))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(message)
		}
	}
}
