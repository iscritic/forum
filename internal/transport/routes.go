package transport

import (
	"context"
	"forum/internal/helpers/template"
	"forum/internal/sqlite"
	"forum/pkg/logger"
	"net/http"
	"time"
)

type application struct {
	logger        *logger.Logger
	storage       *sqlite.Storage
	templateCache *template.TemplateCache
}

func Routes(l *logger.Logger, db *sqlite.Storage, tc *template.TemplateCache) http.Handler {
	mux := http.NewServeMux()

	app := &application{
		logger:        l,
		storage:       db,
		templateCache: tc,
	}

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/post/", app.ViewPostHandler)
	mux.HandleFunc("/post/create", app.CreatePostHandler)
	mux.HandleFunc("/post/comment", app.CreateComment)

	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)

	mux.HandleFunc("/post/like", app.LikePostHandler)
	mux.HandleFunc("/post/dislike", app.DislikePostHandler)
	mux.HandleFunc("/comment/like", app.LikeCommentHandler)
	mux.HandleFunc("/comment/dislike", app.DislikeCommentHandler)

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

		// 3. Validate the session
		session, err := app.storage.GetSessionByToken(cookie.Value)
		if err != nil {
			app.logger.ErrorLog.Printf("Error fetching session: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// 4. Handle session validation and potential renewal
		if session == nil {
			// No session found for the token, redirect to login
			app.logger.ErrorLog.Println("Session not found")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		} else if session.ExpiredAt.Before(time.Now()) {
			// Session expired, delete the old session and redirect to login
			app.logger.InfoLog.Println("Session expired")
			if err := app.storage.DeleteSession(cookie.Value); err != nil {
				app.logger.ErrorLog.Printf("Error deleting expired session: %v", err) // Log the error
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		} else {
			// Session is valid, attach user info to context and proceed
			ctx := context.WithValue(r.Context(), "userID", session.UserID)
			app.logger.InfoLog.Printf("UserID stored in context: %v", ctx.Value("userID"))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	})
}
