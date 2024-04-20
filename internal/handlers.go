package internal

import (
	"forum/internal/sqlite"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Использование шаблона для рендеринга HTML
	renderTemplate(w, "./web/html/home.html", nil)
}

func (app *application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		// Использование шаблона для рендеринга формы создания поста
		renderTemplate(w, "./web/html/post_create.html", nil)

	case http.MethodPost:

		var post sqlite.Post

		title := r.FormValue("title")
		content := r.FormValue("content")

		// Обработка ошибок при получении данных из формы
		if title == "" || content == "" {
			http.Error(w, "Tittle and content are required", http.StatusBadRequest)
			return
		}

		post.Title = title
		post.Content = content
		post.CreationDate = time.Now()
		post.AuthorID = 69
		post.ID = 1

		err := app.storage.CreatePost(post)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := app.storage.GetPostByID(id)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		return
	}

	renderTemplate(w, "./web/html/post_view.html", post)
}
