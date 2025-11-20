package tpl

import (
	"fmt"
	"html/template"
	"maps"
)

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
