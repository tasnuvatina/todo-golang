package handler

import (
	"text/template"

	"github.com/jmoiron/sqlx"
)

type Todo struct {
	ID          int    `db:"id" json:"id"`
	Task        string `db:"title" json:"task"`
	Iscompleted bool   `db:"is_completed" json:"is_complete"`
}

type Handler struct {
	templates *template.Template
	db        *sqlx.DB
}

func New(db *sqlx.DB) *Handler {
	h := &Handler{
		db: db,
	}
	h.parseTemplates()
	return h

}

func (h *Handler) parseTemplates() {
	h.templates = template.Must(template.ParseFiles(
		"templates/create-todo.html",
		"templates/home.html",
		"templates/edit-todo.html",
	))
}
