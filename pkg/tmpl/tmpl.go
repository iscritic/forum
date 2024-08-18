package tmpl

import (
	"html/template"
	"net/http"
)

// RenderTemplate renders a template using the TemplateCache.
func RenderTemplate(w http.ResponseWriter, tc *TemplateCache, tmpl string, data interface{}) error {
	t, ok := tc.Get(tmpl)
	if !ok {
		var err error
		t, err = template.ParseFiles(tmpl)
		if err != nil {
			return err
		}
		tc.Set(tmpl, t)
	}

	err := t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
