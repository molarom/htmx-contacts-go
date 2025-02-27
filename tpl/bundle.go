package tpl

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
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

// TODO: improve error handling, over just panicing
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

	if err := t.ExecuteTemplate(buf, b.base, data); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
	return nil
}
