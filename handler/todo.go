package handler

import (
	"fmt"
	"net/http"
)

type FormData struct {
	Todo  Todo
	Error map[string]string
}

//show create todo form
func (h *Handler) CreateTodo(rw http.ResponseWriter, r *http.Request) {
	todo := Todo{}
	vErrs := map[string]string{}
	h.LoadCreatedTodoForm(rw, todo, vErrs)
}

//take input from the create todo form
func (h *Handler) StoreTodo(rw http.ResponseWriter, r *http.Request) {
	canBeInserted := true
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	task := r.FormValue("Task")
	todo := Todo{
		Task: task,
	}
	if task == "" {
		vErrs := map[string]string{
			"Task": "field is required",
		}
		h.LoadCreatedTodoForm(rw, todo, vErrs)
		return
	}

	if len(task) < 3 {
		vErrs := map[string]string{
			"Task": "field must be greater than or equals to 3",
		}
		h.LoadCreatedTodoForm(rw, todo, vErrs)
		return
	}
	todos := []Todo{}
	h.db.Select(&todos, "SELECT * FROM tasks")
	for _, t := range todos {
		if t.Task == task {
			http.Redirect(rw, r, "/todos/create", http.StatusTemporaryRedirect)
			canBeInserted = false
			break
		}
	}
	if canBeInserted {
		// h.todos =append(h.todos, Todo{Task:task})
		const insertIntoTodo = "INSERT INTO tasks (title,is_completed) VALUES ($1,$2)"
		res := h.db.MustExec(insertIntoTodo, task, false)
		if ok, err := res.RowsAffected(); err != nil || ok == 0 {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
	}

}

//mark complete the uncompleted todos
func (h *Handler) CompleteTodo(rw http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/todos/complete/"):]

	if taskID == "" {
		http.Error(rw, "task is not given", http.StatusInternalServerError)
		return
	}
	const markCompleteTodo = "UPDATE tasks SET is_completed = true WHERE id=$1"
	res := h.db.MustExec(markCompleteTodo, taskID)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

//delete an uncompleted todo from the list

func (h *Handler) DeleteTodo(rw http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/todos/delete/"):]

	if taskID == "" {
		http.Error(rw, "unable to delete empty task", http.StatusInternalServerError)
		return
	}

	const deleteTodoFromTable = "DELETE FROM tasks WHERE id=$1"
	res := h.db.MustExec(deleteTodoFromTable, taskID)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

// show edit todo form

func (h *Handler) EditTodo(rw http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/todos/edit/"):]
	if taskID == "" {
		http.Error(rw, "unable to edit empty task", http.StatusInternalServerError)
		return
	}

	const getTodo = "SELECT * FROM tasks WHERE id=$1"
	todo := Todo{}
	h.db.Get(&todo, getTodo, taskID)
	if todo.ID == 0 {
		http.Error(rw, "Invalid url", http.StatusInternalServerError)
		return
	}
	h.LoadEditTodoForm(rw, todo, map[string]string{})
}

//update todo from the form
func (h *Handler) UpdateTodo(rw http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/todos/update/"):]
	if taskID == "" {
		http.Error(rw, "unable to edit empty task", http.StatusInternalServerError)
		return
	}
	const getTodo = "SELECT * FROM tasks WHERE id=$1"
	todo := Todo{}
	h.db.Get(&todo, getTodo, taskID)
	if todo.ID == 0 {
		http.Error(rw, "Invalid url", http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	taskValue := r.FormValue("Task")
	todo.Task = taskValue
	if taskValue == "" {
		vErrs := map[string]string{
			"Task": "field is required",
		}
		h.LoadEditTodoForm(rw, todo, vErrs)
		return
	}

	if len(taskValue) < 3 {
		vErrs := map[string]string{
			"Task": "field must be greater than or equals to 3",
		}
		h.LoadEditTodoForm(rw, todo, vErrs)
		return
	}

	const updateTodoInTable = "UPDATE tasks SET title = $2 WHERE id = $1"
	res := h.db.MustExec(updateTodoInTable, taskID, taskValue)
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) LoadCreatedTodoForm(rw http.ResponseWriter, todo Todo, errs map[string]string) {
	form := FormData{
		Todo:  todo,
		Error: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "create-todo.html", form); err != nil {
		http.Error(rw, "Unable to execute template", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LoadEditTodoForm(rw http.ResponseWriter, todo Todo, errs map[string]string) {
	form := FormData{
		Todo:  todo,
		Error: errs,
	}

	if err := h.templates.ExecuteTemplate(rw, "edit-todo.html", form); err != nil {
		http.Error(rw, "Unable to execute template", http.StatusInternalServerError)
		return
	}
}
