package delivery

import (
	"context"
	"net/http"
	"time"
)

func (a *application) sessionManager(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Try to get session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				ctx = context.WithValue(ctx, "role", "guest")
				ctx = context.WithValue(ctx, "IsLogin", false)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			a.log.Error("Error getting session cookie: %v", err)
			ctx = context.WithValue(ctx, "role", "guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Fetch session by token
		session, err := a.storage.GetSessionByToken(cookie.Value)
		if err != nil {
			a.log.Error("Error fetching session: %v", err)
			ctx = context.WithValue(ctx, "role", "guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Check if session is valid and not expired
		if session == nil || session.ExpiredAt.Before(time.Now()) {
			a.log.Info("Session not found or expired")
			ctx = context.WithValue(ctx, "role", "guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Fetch the user by userID
		user, err := a.storage.GetUserByID(session.UserID)
		if err != nil {
			a.log.Error("Error fetching user by ID: %v", err)
			ctx = context.WithValue(ctx, "role", "guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Attach userID and role to context
		ctx = context.WithValue(ctx, "userID", user.ID)
		ctx = context.WithValue(ctx, "role", user.Role)
		ctx = context.WithValue(ctx, "IsLogin", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *application) requireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("role").(string)
			if !ok || role == "guest" {
				a.log.Warn("No role found in context, redirecting to login")
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Role not allowed, 403 - forbidden
			a.log.Info("Role '%s' not authorized for this route", role)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		})
	}
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		//w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (a *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.log.Info("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (a *application) recoverPanic(next http.Handler) http.Handler {
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
