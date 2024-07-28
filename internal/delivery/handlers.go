package delivery

import (
	"forum/internal/helpers/template"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/internal/service/session"
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

	posts, err := service.FetchPosts(app.storage)
	if err != nil {
		app.logger.ErrorLog.Println(err)
		return
	}

	// Использование шаблона для рендеринга HTML
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

		var post repository.Post

		title := r.Form.Get("title")
		content := r.Form.Get("content")

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
	// Проверяем, что метод запроса GET
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

	postData, err := service.GetPostData(app.storage, id)
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

	var comment repository.Comment

	// Получаем ID поста из формы
	postIDStr := r.Form.Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comment.PostID = postID
	comment.Content = r.Form.Get("content")

	// TODO author comment id

	err = app.storage.CreateComment(comment)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

func (app *application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		template.RenderTemplate(w, app.templateCache, "./web/html/register.html", nil)
		return
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			return
		}

		// Получаем данные из формы регистрации
		username := r.Form.Get("username")
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		// Создаем нового пользователя
		newUser := repository.User{
			Username: username,
			Email:    email,
			Password: password,
			Role:     "user", // Устанавливаем роль по умолчанию
		}

		service.Register(app.storage, newUser)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		app.logger.InfoLog.Printf("New user detected: %v", newUser)
	}
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		template.RenderTemplate(w, app.templateCache, "./web/html/login.html", nil)
		return
	case r.Method != http.MethodPost:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		// Обработка ошибки
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Получаем данные из формы авторизации
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Проверяем существование пользователя с заданными данными
	user, err := app.storage.GetUserByUsername(username)
	if err != nil {
		// Обработка ошибки
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	app.logger.ErrorLog.Println("we are here")

	// Проверяем соответствие пароля
	if user.Password != password {
		// Обработка ошибки
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	app.logger.ErrorLog.Println("we are here")

	// Создаем сессию и устанавливаем cookie
	sessionToken, err := session.CreateSession(app.storage, user)
	if err != nil {
		app.logger.InfoLog.Println(err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(20 * time.Minute),
		Path:    "/",
	})

	// Перенаправляем пользователя на домашнюю страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
