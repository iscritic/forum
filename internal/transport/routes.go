package transport

import (
	"context"
	"forum/internal/sqlite"
	"forum/pkg/logger"
	"net/http"
	"time"
)

type application struct {
	logger  *logger.Logger
	storage *sqlite.Storage
}

func Routes(l *logger.Logger, db *sqlite.Storage) http.Handler {
	mux := http.NewServeMux()

	app := &application{
		logger:  l,
		storage: db,
	}

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.HomeHandler)

	mux.HandleFunc("/post/view", app.ViewPostHandler)
	mux.HandleFunc("/post/create", app.CreatePostHandler)
	mux.HandleFunc("/post/comment", app.CreateComment)

	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)

	return app.SessionMiddleware(mux)
}

func (app *application) SessionMiddleware(next http.Handler) http.Handler {
	excludedPaths := map[string]bool{
		"/login":    true,
		"/register": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.InfoLog.Printf("SessionMiddleware called for: %s", r.URL.Path)

		// 1. Check if the path is excluded from the middleware
		if _, ok := excludedPaths[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		// 2. Get session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err != http.ErrNoCookie {
				app.logger.ErrorLog.Printf("Error getting session cookie: %v", err)
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		app.logger.InfoLog.Printf("Your session token: %v", cookie.Value)

		// 3. Validate the session and delete if exists
		session, err := app.storage.GetSessionByToken(cookie.Value)
		if err != nil {
			app.logger.ErrorLog.Printf("Error fetching session: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if session != nil {
			app.logger.InfoLog.Println("Deleting existing session")
			err = app.storage.DeleteSession(cookie.Value)
			if err != nil {
				app.logger.ErrorLog.Printf("Error deleting session: %v", err)
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
		}

		// 4. Create a new session
		user, err := app.storage.GetUserByID(session.UserID)
		if err != nil {
			app.logger.ErrorLog.Printf("Error fetching user: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		newSessionToken, err := session.CreateSession(app.storage, user)
		if err != nil {
			app.logger.ErrorLog.Printf("Error creating new session: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// 5. Set new cookie and attach user info to context
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    newSessionToken,
			Expires:  time.Now().Add(20 * time.Minute),
			Path:     "/",
			HttpOnly: true,
			Secure:   app.config.IsProduction,
		})
		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		app.logger.InfoLog.Printf("UserID stored in context: %v", ctx.Value("userID"))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
