package tmpl

import (
	"html/template"
	"path/filepath"
	"sync"
)

// TemplateCache provides a thread-safe cache for storing parsed templates.
type TemplateCache struct {
	cache map[string]*template.Template
	mu    sync.RWMutex
}

// NewTemplateCache creates and initializes a new TemplateCache.
func NewTemplateCache() *TemplateCache {
	return &TemplateCache{
		cache: make(map[string]*template.Template),
	}
}

// Get retrieves a template from the cache by its file path.
// It returns the template and a boolean indicating whether it was found.
func (tc *TemplateCache) Get(tmpl string) (*template.Template, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	t, ok := tc.cache[tmpl]
	return t, ok
}

// Set adds a template to the cache, associating it with its file path.
func (tc *TemplateCache) Set(tmpl string, t *template.Template) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.cache[tmpl] = t
}

// LoadTemplates parses HTML templates from a directory and adds them to the cache.
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
