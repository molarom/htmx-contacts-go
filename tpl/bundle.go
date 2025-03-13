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

// Data represents the data to be rendered for a template.
type Data map[string]any

// NewBundle creates a new bundle of templates.
//
// The baseTpl is the base template definition for which the all the templates
// inherit from.
func NewBundle(baseTpl, viewsDir string, sharedDirs ...string) *Bundle {
	b := &Bundle{
		base: baseTpl,
		tree: radix.New[byte, *template.Template](),
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}

	views, err := filepath.Glob(viewsDir)
	if err != nil {
		panic(err)
	}

	shared := make([][]string, 0, len(sharedDirs))
	for _, d := range sharedDirs {
		s, err := filepath.Glob(d)
		if err != nil {
			panic(err)
		}
		shared = append(shared, s)
	}

	for _, v := range views {
		var f []string
		for _, s := range shared {
			f = append(f, append(s, v)...)
		}
		b.tree.Insert([]byte(filepath.Base(v)),
			template.Must(
				template.New(f[0]).Funcs(funcs).ParseFiles(f...),
			),
		)
	}

	// TODO: tidy up tpl inheritance config.
	for _, s := range shared {
		for _, f := range s {
			b.tree.Insert([]byte(filepath.Base(f)),
				template.Must(
					template.New(f).Funcs(funcs).ParseFiles(f),
				),
			)
		}
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
	_, err := w.Write(buf.Bytes())
	return err
}

func (b *Bundle) Print() {
	b.tree.Print()
}
