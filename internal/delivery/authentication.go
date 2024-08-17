package delivery

import (
	"forum/internal/entity"
	"forum/internal/helpers/tmpl"
	"forum/internal/service"
	"forum/internal/service/session"
	"net/http"
	"time"
)

func (app *application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl.RenderTemplate(w, app.tmplcache, "./web/html/register.html", nil)
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
		newUser := entity.User{
			Username: username,
			Email:    email,
			Password: password,
			Role:     "user", // Устанавливаем роль по умолчанию
		}

		service.Register(app.storage, newUser)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		app.log.InfoLog.Printf("New user detected: %v", newUser)
	}
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		tmpl.RenderTemplate(w, app.tmplcache, "./web/html/login.html", nil)
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

	// Проверяем соответствие пароля
	if user.Password != password {
		// Обработка ошибки
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	app.log.ErrorLog.Println("we are here")

	// Создаем сессию и устанавливаем cookie
	sessionToken, err := session.CreateSession(app.storage, user)
	if err != nil {
		app.log.InfoLog.Println(err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(20 * time.Minute),
		Path:     "/",
		HttpOnly: true,
	})

	// Перенаправляем пользователя на домашнюю страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the userID from the context
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Delete all sessions for the user
	err := app.storage.DeleteAllSessionsForUser(userId)
	if err != nil {
		app.log.ErrorLog.Printf("Failed to delete existing sessions: %v", err)
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Expire the cookie
		HttpOnly: true,
	})

	// Redirect to the home page after successful logout
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
