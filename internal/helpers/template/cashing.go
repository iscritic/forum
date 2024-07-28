package template

import (
	"path/filepath"
	"sync"
	"text/template"
)

// кэш мэш
type TemplateCache struct {
	cache map[string]*template.Template
	mu    sync.RWMutex
}

// инит кэша темплейта
func NewTemplateCache() *TemplateCache {
	return &TemplateCache{
		cache: make(map[string]*template.Template),
	}
}

// проверка на наличие
func (tc *TemplateCache) Get(tmpl string) (*template.Template, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	t, ok := tc.cache[tmpl]
	return t, ok
}

// кешируем новый темплейт
func (tc *TemplateCache) Set(tmpl string, t *template.Template) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.cache[tmpl] = t
}

func LoadTemplates(tc *TemplateCache, dir string) error {
	pattern := filepath.Join(dir, "*.html")
	templates, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, tmpl := range templates {
		t, err := template.ParseFiles(tmpl)
		if err != nil {
			return err
		}
		tc.Set(tmpl, t)
	}

	return nil
}
