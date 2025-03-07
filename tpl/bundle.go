package tpl

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"gitlab.com/romalor/radix"
)

// Bundle represents a collection of templates that
// extend a shared base.
type Bundle struct {
	base string
	tree *radix.Tree[byte, *template.Template]
	pool sync.Pool
}

type Data map[string]any

func NewBundle(baseTpl, layoutDir, contentDir string) *Bundle {
	b := &Bundle{
		base: baseTpl,
		tree: radix.New[byte, *template.Template](),
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}

	layouts, err := filepath.Glob(layoutDir)
	if err != nil {
		panic(err)
	}

	views, err := filepath.Glob(contentDir)
	if err != nil {
		panic(err)
	}

	for _, v := range views {
		f := append(layouts, v)
		b.tree.Insert([]byte(filepath.Base(v)),
			template.Must(
				template.New(f[0]).Funcs(funcs).ParseFiles(f...),
			),
		)
	}

	return b
}

func (b *Bundle) Render(w http.ResponseWriter, name string, data Data) error {
	t, ok := b.tree.Get([]byte(name))
	if !ok {
		return fmt.Errorf("render: %s not a valid template", name)
	}

	buf := b.pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer b.pool.Put(buf)

	f := template.FuncMap{
		"include": func(name string) (string, error) {
			v, ok := b.tree.Get([]byte(name))
			if !ok {
				return "", fmt.Errorf("no template named: %s", name)
			}

			var w *strings.Builder
			if err := v.New(name).ExecuteTemplate(w, name, data); err != nil {
				return "", err
			}
			return w.String(), nil
		},
	}

	if err := t.Funcs(f).ExecuteTemplate(buf, b.base, data); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write(buf.Bytes())
	return err
}
