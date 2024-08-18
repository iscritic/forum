package tmpl

import (
	"net/http"
)

// RenderErrorPage renders a custom error page using the TemplateCache.
func RenderErrorPage(w http.ResponseWriter, tc *TemplateCache, statusCode int, message string) {
	tmpl, ok := tc.Get("error.html")
	if !ok {
		http.Error(w, "Error page template not found", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	err := tmpl.Execute(w, map[string]interface{}{
		"ErrorCode":    statusCode,
		"ErrorMessage": message,
	})
	if err != nil {
		http.Error(w, "Failed to render error page", http.StatusInternalServerError)
	}
}
