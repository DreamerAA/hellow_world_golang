module main.go

go 1.23.1

require (
	github.com/gorilla/mux v1.8.1
	github.com/sirupsen/logrus v1.9.3
	psql v0.0.0
)

require (
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace psql => ../psql
