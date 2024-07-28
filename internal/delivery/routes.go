package delivery

import (
	"context"
	"forum/internal/helpers/template"
	"forum/internal/repository"
	"forum/pkg/logger"
	"net/http"
	"time"
)

type application struct {
	logger        *logger.Logger
	storage       *repository.Storage
	templateCache *template.TemplateCache
}

func Routes(l *logger.Logger, db *repository.Storage, tc *template.TemplateCache) http.Handler {
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

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
