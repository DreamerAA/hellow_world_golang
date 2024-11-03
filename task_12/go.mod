module main.go

go 1.23.1

require (
	github.com/gorilla/mux v1.8.1
	psql v0.0.0
)

require github.com/lib/pq v1.10.9 // indirect

replace psql => ../psql
