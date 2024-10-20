package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lib/pq"
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

func splitParams(text string) (map[string]string, int, error) {

	var new_params = make(map[string]string)
	id := -1
	if text == "" {
		return new_params, id, nil
	}
	args := strings.Split(text, " ")
	for _, arg := range args {
		data := strings.Split(arg, "=")
		if len(data) != 2 {
			return nil, -1, fmt.Errorf("Invalid input")
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
	return new_params, id, nil
}

func createTable(db *sql.DB, table_name string) {
	query_enum := "CREATE TYPE taskstatus AS ENUM ('Opened', 'InProgress', 'Completed');"
	_, err := db.Exec(query_enum)
	if err != nil {
		if pg_err, ok := err.(*pq.Error); ok {
			if pg_err.Code != pq.ErrorCode(42710) {
				fmt.Println("enum уже существует:", pg_err.Message)
			} else {
				fmt.Println("Ошибка создания enum:", err, pg_err.Message, pg_err.Code)
			}
		}
	} else {
		fmt.Println("enum успешно создан!")
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

func createSelectQuery(header string, filters map[string]string, order_params string) (string, []interface{}) {
	query := header
	args := []interface{}{}
	argCount := 1

	// Добавление фильтров с плейсхолдерами
	if len(filters) > 0 {
		query += " WHERE "
		for key, value := range filters {
			query += fmt.Sprintf("%s = $%d", key, argCount)
			if argCount > 1 {
				query += " AND "
			}
			args = append(args, value)
			argCount++
		}
	}
	// Добавление параметров сортировки
	if order_params != "" {
		query += " ORDER BY " + order_params
	}

	return query, args
}

func createQuery(header string, new_params map[string]string, splitter string) (string, []interface{}) {
	query := header
	ind := 1
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

func insertArgsToQuery(header string, new_params map[string]string, splitter string) string {
	query := header //"INSERT INTO tasks (title, description) VALUES "
	ind := 1
	for k, v := range new_params {
		query += k + splitter + v
		if ind < len(new_params) {
			query += ", "
		}
		ind += 1
	}
	return query
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
				ftext := ""
				otext := ""
				if strings.ContainsAny(text, "|") {
					ar_str := strings.Split(text, "|")
					if len(ar_str) != 2 {
						return "", fmt.Errorf("Invalid input")
					}
					ftext = ar_str[0]
					otext = ar_str[1]
				} else {
					ftext = text
					otext = "status ASC, title DESC"
				}
				filters, _, err := splitParams(ftext)
				if err != nil {
					return "", fmt.Errorf("Invalid input")
				}
				var query string
				queryArgs := []interface{}{}
				if len(filters) > 0 {
					// Создание запроса с фильтрами и сортировкой
					query, queryArgs = createSelectQuery("SELECT id, title, description, status FROM tasks", filters, otext)
				} else {
					query = "SELECT id, title, description, status FROM tasks ORDER BY " + otext + ";"
				}
				rows, err := db.Query(query, queryArgs...)
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
				new_params, id, err := splitParams(text)
				if err != nil || len(new_params) == 0 || id == -1 {
					return "", fmt.Errorf("Invalid input")
				}

				query, query_args := createQuery("UPDATE tasks SET ", new_params, "= ")
				query += fmt.Sprintf(" WHERE id = $%d;", len(new_params)+1)
				query_args = append(query_args, id)
				fmt.Println("query=", query)
				fmt.Println("args=", query_args)
				_, err = db.Exec(query, query_args...)
				if err != nil {
					return "", err
				}
				return fmt.Sprintf("Task with id = %d updated", id), nil
			},
		},
	}

	// 	1 Поиск задач по определённым критериям (например, статус):
	// Команда search status=Opened может позволить искать только открытые задачи.
	// 2 Команда для завершения задачи:
	// Создай команду, которая изменяет статус задачи на Completed.

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter query: ")
		ok := scanner.Scan()
		if !ok {
			fmt.Println("Error reading input")
			continue
		}
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
	// value, _ := commands["list"]
	// _, err = value.operation("")
}
