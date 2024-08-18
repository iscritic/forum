package delivery

import (
	"forum/internal/service"
	"forum/internal/service/session"
	tmpl2 "forum/pkg/tmpl"
	"net/http"
	"time"
)

func (a *application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		err := tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/register.html", nil)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	case http.MethodPost:

		user, err := service.DecodeUser(r)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		err = service.Register(a.storage, user)
		if err != nil {
			a.log.Error(err.Error())
			tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		a.log.Info("New user detected: %v", user)
	}
}

func (a *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/login.html", nil)
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
	user, err := a.storage.GetUserByUsername(username)
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

	// Создаем сессию и устанавливаем cookie
	sessionToken, err := session.CreateSession(a.storage, user)
	if err != nil {
		a.log.Error("CreateSession: %v", err)
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, "Internal Server Error")
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

func (a *application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the userID from the context
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Delete all sessions for the user
	err := a.storage.DeleteAllSessionsForUser(userId)
	if err != nil {
		a.log.Error("Failed to delete existing sessions: %v", err)
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
