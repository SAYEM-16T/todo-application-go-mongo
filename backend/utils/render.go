package utils

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"todo-app/backend/models"
)

func RenderAppPage(w http.ResponseWriter, email string, todos []models.Todo) {
	var b strings.Builder
	b.WriteString(`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">`)
	b.WriteString(`<title>TODO App</title><link rel="stylesheet" href="/static/styles/main.css"></head><body>`)
	b.WriteString(`<div class="container"><nav class="nav"><a class="brand" href="/app">Todo</a><div class="spacer"></div><div class="user">`)
	b.WriteString(html.EscapeString(email))
	b.WriteString(`</div><form class="inline" action="/logout" method="POST"><button class="btn btn-link" type="submit">Logout</button></form></nav>`)
	b.WriteString(`<main class="card"><h1>Your Tasks</h1>
   <form class="add-form" action="/todo" method="POST">
     <input name="title" placeholder="Add a new task..." maxlength="100" required>
     <button class="btn btn-primary" type="submit">Add</button>
   </form>
   <ul class="todo-list">`)
	if len(todos) == 0 {
		b.WriteString(`<li class="empty">No tasks yet. Add your first task!</li>`)
	}
	for _, t := range todos {
		title := html.EscapeString(t.Title)
		cls := ""
		btn := "Done"
		if t.Done {
			cls = " completed"
			btn = "Undo"
		}
		b.WriteString(`<li class="todo-item` + cls + `">`)
		b.WriteString(`<span class="title">` + title + `</span>
       <div class="actions">
         <form class="inline" action="/todo/` + t.ID.Hex() + `/toggle" method="POST"><button class="btn btn-secondary" type="submit">` + btn + `</button></form>
         <form class="inline" action="/todo/` + t.ID.Hex() + `/delete" method="POST"><button class="btn btn-danger" type="submit">Delete</button></form>
       </div>`)
		b.WriteString(`</li>`)
	}
	b.WriteString(`</ul></main><footer class="footer">Built with Go + MongoDB</footer></div></body></html>`)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, b.String())
}
