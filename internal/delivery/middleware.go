package delivery

import (
	"context"
	"net/http"
	"time"
)

func (app *application) requiredAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.InfoLog.Printf("SessionMiddleware called for: %s", r.URL.Path)

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
		defer func() {
			if err := recover(); err != nil {

				w.Header().Set("Connection", "close")

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
