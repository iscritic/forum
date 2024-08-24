package delivery

import (
	"net/http"
	"time"

	"forum/internal/service"
	"forum/internal/service/session"

	tmpl2 "forum/pkg/tmpl"
)

type TemplateData struct {
	Error string
}

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
			data := TemplateData{
				Error: err.Error(),
			}
			tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/register.html", data)
			return
		}

		err = service.Register(a.storage, user)
		if err != nil {
			a.log.Error(err.Error())
			data := TemplateData{
				Error: err.Error(),
			}
			tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/register.html", data)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		a.log.Info("New user detected: %v", user)
	}
}

func (a *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/login.html", nil)
		return
	case r.Method != http.MethodPost:
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	err := r.ParseForm()
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	user, err := a.storage.GetUserByUsername(username)
	if err != nil || user.Password != password {
		data := struct {
			Error string
		}{
			Error: "Password or username are not correct, try again",
		}
		tmpl2.RenderTemplate(w, a.tmplcache, "./web/html/login.html", data)
		return
	}

	sessionToken, err := session.CreateSession(a.storage, user)
	if err != nil {
		a.log.Error(err.Error())
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(20 * time.Minute),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		a.log.Error("User ID not found in context")
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err := a.storage.DeleteAllSessionsForUser(userId)
	if err != nil {
		a.log.Error("Failed to delete existing sessions: %v", err)
		tmpl2.RenderErrorPage(w, a.tmplcache, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
