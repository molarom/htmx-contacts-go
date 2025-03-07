package tpl

import (
	"fmt"
	"html/template"
)

var funcs = template.FuncMap{
	// arithmetic
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },

	// Template nesting
	"include": func(name string) (string, error) {
		return "", fmt.Errorf("empty include called")
	},
}
