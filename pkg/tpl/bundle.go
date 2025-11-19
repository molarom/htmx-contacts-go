package tpl

import (
	"bytes"
	"fmt"
	"html/template"
	"maps"
	"net/http"
	"path/filepath"
	"sync"
)

// Bundle represents a collection of templates that extend a shared base.
type Bundle struct {
	base      string
	templates map[string]*template.Template
	pool      sync.Pool
	mu        sync.RWMutex
	funcs     template.FuncMap
}

// Data represents the data to be rendered for a template.
type Data map[string]any

// Option is a functional option for configuring a Bundle.
type Option func(*Bundle) error

// WithBaseTpl sets the base template name that other templates extend.
func WithBaseTpl(name string) Option {
	return func(b *Bundle) error {
		if name == "" {
			return fmt.Errorf("base template name cannot be empty")
		}
		b.base = name
		return nil
	}
}

// WithFuncs adds custom template functions.
func WithFuncs(funcMap template.FuncMap) Option {
	return func(b *Bundle) error {
		if b.funcs == nil {
			b.funcs = make(template.FuncMap)
		}

		maps.Copy(b.funcs, funcMap)
		return nil
	}
}

// NewBundle creates a new bundle of templates with the given options.
//
// The viewsDir is a glob pattern for view templates (e.g., "views/*.html").
// The sharedDirs are glob patterns for shared/partial templates.
func NewBundle(viewsDir string, sharedDirs []string, opts ...Option) (*Bundle, error) {
	b := &Bundle{
		base:      "base", // default base template name
		templates: make(map[string]*template.Template),
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
		funcs: funcs,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(b); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	// Load templates
	if err := b.load(viewsDir, sharedDirs); err != nil {
		return nil, fmt.Errorf("loading templates: %w", err)
	}

	return b, nil
}

// load handles the template loading logic.
func (b *Bundle) load(viewsDir string, sharedDirs []string) error {
	// Load views
	views, err := filepath.Glob(viewsDir)
	if err != nil {
		return fmt.Errorf("globbing views directory %q: %w", viewsDir, err)
	}

	if len(views) == 0 {
		return fmt.Errorf("no templates found matching pattern %q", viewsDir)
	}

	// Load shared templates
	var allShared []string
	for _, pattern := range sharedDirs {
		shared, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("globbing shared directory %q: %w", pattern, err)
		}
		allShared = append(allShared, shared...)
	}

	// Parse view templates (with shared templates as dependencies)
	for _, viewPath := range views {
		name := filepath.Base(viewPath)
		files := append([]string{viewPath}, allShared...)

		tpl, err := b.parseTemplate(files)
		if err != nil {
			return fmt.Errorf("parsing template %q: %w", name, err)
		}

		b.templates[name] = tpl
	}

	// Also register shared templates individually
	for _, sharedPath := range allShared {
		name := filepath.Base(sharedPath)

		tpl, err := b.parseTemplate([]string{sharedPath})
		if err != nil {
			return fmt.Errorf("parsing shared template %q: %w", name, err)
		}

		b.templates[name] = tpl
	}

	return nil
}

// parseTemplate creates a new template from the given files.
func (b *Bundle) parseTemplate(files []string) (*template.Template, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	tpl := template.New(filepath.Base(files[0]))

	if len(b.funcs) > 0 {
		tpl = tpl.Funcs(b.funcs)
	}

	tpl, err := tpl.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

// Render renders the named template with the given data to the http.ResponseWriter.
func (b *Bundle) Render(w http.ResponseWriter, name string, data Data) error {
	b.mu.RLock()
	t, ok := b.templates[name]
	b.mu.RUnlock()

	if !ok {
		return fmt.Errorf("template %q not found", name)
	}

	buf := b.pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer b.pool.Put(buf)

	// Determine which template to execute
	tplName := name
	if t.Lookup(b.base) != nil {
		tplName = b.base
	}

	if err := t.ExecuteTemplate(buf, tplName, data); err != nil {
		return fmt.Errorf("executing template %q: %w", name, err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}

// HasTemplate checks if a template with the given name exists.
func (b *Bundle) HasTemplate(name string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.templates[name]
	return ok
}

// TemplateNames returns a list of all registered template names.
func (b *Bundle) TemplateNames() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	names := make([]string, 0, len(b.templates))
	for name := range b.templates {
		names = append(names, name)
	}
	return names
}

// Validate checks that all templates can be executed without errors.
func (b *Bundle) Validate() error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var buf bytes.Buffer
	testData := Data{}

	for name, t := range b.templates {
		tplName := name
		if t.Lookup(b.base) != nil {
			tplName = b.base
		}

		buf.Reset()
		if err := t.ExecuteTemplate(&buf, tplName, testData); err != nil {
			return fmt.Errorf("validating template %q: %w", name, err)
		}
	}

	return nil
}
