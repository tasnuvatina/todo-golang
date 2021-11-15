package main

import (
	"log"
	"net/http"
	//"database/sql"

	"golang-practice/todo/handler"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	//declaring table schema
	var schema = `
		CREATE TABLE IF NOT EXISTS  tasks (
			id serial,
			title text,
			is_completed boolean,

			primary key(id)
		);`

	db, err := sqlx.Connect("postgres", "user=postgres password=password dbname=todos sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	h := handler.New(db)

	http.HandleFunc("/", h.Home)
	http.HandleFunc("/todos/create", h.CreateTodo)
	http.HandleFunc("/todos/store", h.StoreTodo)
	http.HandleFunc("/todos/complete/", h.CompleteTodo)
	http.HandleFunc("/todos/delete/", h.DeleteTodo)
	http.HandleFunc("/todos/edit/", h.EditTodo)
	http.HandleFunc("/todos/update/", h.UpdateTodo)

	if err := http.ListenAndServe("127.0.0.1:3000", nil); err != nil {
		log.Fatal(err)
	}
}
