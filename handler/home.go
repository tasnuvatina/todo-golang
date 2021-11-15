package handler

import (
	"net/http"
)
type ListTodo struct{
	Todos []Todo
}

func (h *Handler) Home(rw http.ResponseWriter, r *http.Request) {
	todos := []Todo{}
	h.db.Select(&todos, "SELECT * FROM tasks")
	lt := ListTodo{
		Todos:todos,
	}
	if err := h.templates.ExecuteTemplate(rw,"home.html",lt);err!=nil{
		http.Error(rw,"unable to execute home template",http.StatusInternalServerError)
		return
	}
}