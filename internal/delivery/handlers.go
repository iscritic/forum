package delivery

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/helpers/template"
	"forum/internal/service"
	"net/http"
	"strconv"
	"strings"
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

	posts, err := service.GetAllPostRelatedData(app.storage)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		return
	}

	template.RenderTemplate(w, app.templateCache, "./web/html/home.html", posts)
}

func (app *application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		// Использование шаблона для рендеринга формы создания поста
		template.RenderTemplate(w, app.templateCache, "./web/html/post_create.html", nil)
		return

	case http.MethodPost:

		err := r.ParseForm()
		if err != nil {
			return
		}

		var post entity.Post

		title := r.Form.Get("title")
		content := r.Form.Get("content")

		if title == "" || content == "" {
			http.Error(w, "Tittle and content are required", http.StatusBadRequest)
			return
		}

		post.Title = title
		post.Content = content
		post.CreationDate = time.Now()

		//  r.Context().Value("userID")
		post.AuthorID = r.Context().Value("userID").(int)
		post.CategoryID = 9

		fmt.Println(post.AuthorID)

		lastID, err := app.storage.CreatePost(post)
		if err != nil {
			app.logger.ErrorLog.Println(err)
			return
		}

		http.Redirect(w, r, "/post/"+strconv.Itoa(lastID), http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	postData, err := service.GetPostRelatedData(app.storage, id)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		return
	}

	template.RenderTemplate(w, app.templateCache, "./web/html/post_view.html", postData)
}

func (app *application) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	var comment entity.Comment

	// Получаем ID поста из формы
	postIDStr := r.Form.Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comment.PostID = postID
	comment.Content = r.Form.Get("content")

	comment.AuthorID = r.Context().Value("userID").(int)

	// TODO author comment id

	err = app.storage.CreateComment(comment)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}
