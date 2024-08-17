package delivery

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/internal/entity"
	"forum/internal/helpers/tmpl"
	"forum/internal/service"
)

func (a *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := service.GetHomePageData(a.storage, r.Context())
	if err != nil {
		a.log.Error(err.Error())
		tmpl.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = tmpl.RenderTemplate(w, a.tmplcache, "./web/html/home.html", data)
	if err != nil {
		a.log.Error(err.Error())
		tmpl.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (a *application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:

		categories, err := service.GetCategories(a.storage)
		if err != nil {
			a.log.Error(err.Error())
			tmpl.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		err = tmpl.RenderTemplate(w, a.tmplcache, "./web/html/post_create.html", categories)
		if err != nil {
			a.log.Error(err.Error())
			tmpl.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		return

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			a.log.Error(err.Error())
			tmpl.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		var post entity.Post

		title := r.Form.Get("title")
		content := r.Form.Get("content")
		categoryIDStr := r.Form.Get("category")

		if title == "" || content == "" || categoryIDStr == "" {
			http.Error(w, "Title, content, and category are required", http.StatusBadRequest)
			return
		}

		post.Title = title
		post.Content = content
		post.CreationDate = time.Now()

		post.AuthorID = r.Context().Value("userID").(int)

		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		post.CategoryID = categoryID

		lastID, err := a.storage.CreatePost(post)
		if err != nil {
			a.log.ErrorLog.Println(err)
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
		app.log.ErrorLog.Println(err)
		return
	}

	tmpl.RenderTemplate(w, app.tmplcache, "./web/html/post_view.html", postData)
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
