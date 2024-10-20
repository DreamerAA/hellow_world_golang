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

type Task struct {
	id          int
	title       string
	description string
	status      TaskStatus
}

type TaskStatus int

const (
	Opened = iota
	InProgress
	Completed
)

var (
	capabilitiesMap = map[string]TaskStatus{
		"Opened":     Opened,
		"InProgress": InProgress,
		"Completed":  Completed,
	}
)

func ParseString(str string) (TaskStatus, bool) {
	c, ok := capabilitiesMap[str]
	return c, ok
}
func (s TaskStatus) String() string {
	switch s {
	case Opened:
		return "Открыт"
	case InProgress:
		return "В процессе"
	case Completed:
		return "Завершен"
	default:
		return "Неизвестно"
	}
}

func createTable(db *sql.DB, table_name string) {
	query_enum := "CREATE TYPE taskstatus AS ENUM ('Opened', 'InProgress', 'Completed');"
	_, err := db.Exec(query_enum)
	if err != nil {
		fmt.Println("Ошибка создания enum:", err)
	} else {
		fmt.Println("enum успешно создан или уже существует!")
	}

	query :=
		"CREATE TABLE IF NOT EXISTS " + table_name +
			"(id SERIAL PRIMARY KEY, title TEXT NOT NULL, description TEXT, status taskstatus DEFAULT 'Opened');"

	_, err = db.Exec(query)
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
				query := `INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id;`
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
				data := strings.Split(text, "=")
				if len(data) != 2 || data[0] != "id" {
					return "", fmt.Errorf("Invalid input")
				}

				id, err := strconv.Atoi(data[1])

				if err != nil {
					fmt.Printf("%q does not looks like a number.\n", data[1])
					return "", fmt.Errorf("Invalid input")
				}
				query := `DELETE FROM tasks WHERE id = $1`
				_, err = db.Exec(query, id)
				if err != nil {
					return "", err
				}
				return fmt.Sprintf("Tast deleted with id = %d", id), nil
			},
		},
		"list": {
			"List tasks",
			func(text string) (string, error) {
				sortBy, sortOrder := "status", "ASC"
				query := fmt.Sprintf("SELECT id, title, description, status FROM tasks ORDER BY %s %s;", sortBy, sortOrder)
				rows, err := db.Query(query)
				if err != nil {
					return "", err
				}
				defer rows.Close()

				var tasks []string
				var task Task
				for rows.Next() {
					str_status := ""
					err := rows.Scan(&task.id, &task.title, &task.description, &str_status)

					if err != nil {
						return "", err
					}
					var ok bool
					task.status, ok = ParseString(str_status)
					if !ok {
						return "", err
					}
					result := strconv.Itoa(task.id) + " "
					result += task.title + " "
					result += task.status.String() + " "
					result += task.description
					tasks = append(tasks, result)
				}
				return strings.Join(tasks, "\n"), nil
			},
		},
		"update": {
			"Update task",
			func(text string) (string, error) {
				args := strings.Split(text, " ")
				var new_params = make(map[string]string)
				id := -1
				for _, arg := range args {
					data := strings.Split(arg, "=")
					if len(data) != 2 {
						return "", fmt.Errorf("Invalid input")
					}
					name_param := data[0]
					if name_param == "id" {
						if id != -1 {
							fmt.Println("You can use only one id")
						}
						lid, err := strconv.Atoi(data[1])
						if err != nil {
							fmt.Printf("%q does not looks like a number.\n", data[1])
						} else {
							id = lid
						}
					} else {
						new_params[name_param] = data[1]
					}
				}
				if len(new_params) == 0 || id == -1 {
					return "", fmt.Errorf("Invalid input")
				}
				query := "UPDATE tasks SET "
				ind := 1
				for k, v := range new_params {
					query += k + "='" + v + "'"
					ind = ind + 1
					if ind < len(new_params) {
						query += ", "
					}
				}
				query += " WHERE id = $1;"
				fmt.Println("query=", query)
				_, err := db.Exec(query, id)

				if err != nil {
					return "", err
				}
				return fmt.Sprintf("Task with id = %d updated", id), nil
			},
		},
	}

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
