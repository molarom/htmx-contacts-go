package tpl

import (
	"html/template"
)

var funcs = template.FuncMap{
	// arithmetic
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
	"mul": func(a, b float64) int { return int(a * b) },
}
