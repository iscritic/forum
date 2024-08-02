package template

import (
	"net/http"
	"text/template"
)

// RenderTemplate renders a template using the TemplateCache.
func RenderTemplate(w http.ResponseWriter, tc *TemplateCache, tmpl string, data interface{}) error {
	t, ok := tc.Get(tmpl)
	if !ok {
		var err error
		t, err = template.ParseFiles(tmpl)
		if err != nil {
			return err // Return error to the caller
		}
		tc.Set(tmpl, t)
	}

	err := t.Execute(w, data)
	if err != nil {
		return err // Return error to the caller
	}
	return nil
}
