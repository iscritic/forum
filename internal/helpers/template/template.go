package template

import (
	"net/http"
	"text/template"
)

func RenderTemplate(w http.ResponseWriter, tc *TemplateCache, tmpl string, data interface{}) {
	t, ok := tc.Get(tmpl)
	if !ok {
		var err error
		t, err = template.ParseFiles(tmpl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tc.Set(tmpl, t)
	}
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
