package web

import (
	"errors"
	"html/template"
	"io"
	"sync/atomic"
)

type Renderer struct {
	tmpl atomic.Value
}

var (
	ErrNilRenderer error = errors.New("renderer is nil")
	ErrNoTemplates error = errors.New("renderer doesn't have any templates")

	renderer Renderer
)

func HasTemplate(name string) bool {
	return renderer.HasTemplate(name)
}

func HasAnyTemplates() bool {
	return renderer.HasAnyTemplates()
}

func LoadTemplates() error {
	return renderer.LoadTemplates(GlobalSettings.TemplatesPattern)
}

func CleanUpTemplates() {
	renderer.CleanUpTemplates()
}

func Render(w io.Writer, name string, data interface{}) error {
	if !HasAnyTemplates() {
		err := LoadTemplates()
		if err != nil {
			return err
		}
	}
	return renderer.Render(w, name, data)
}

func (r *Renderer) HasTemplate(name string) bool {
	if r == nil {
		return false
	}
	x := r.tmpl.Load()
	if x == nil {
		return false
	}
	t := x.(*template.Template)
	return t != nil && t.Lookup(name) != nil
}

func (r *Renderer) HasAnyTemplates() bool {
	if r == nil {
		return false
	}
	return r.tmpl.Load() != nil
}

func (r *Renderer) LoadTemplates(pattern string) error {
	if r == nil {
		return ErrNilRenderer
	}
	t, err := template.ParseGlob(pattern)
	if err != nil {
		return err
	}
	r.tmpl.Store(t)
	return nil
}

func (r *Renderer) CleanUpTemplates() {
	if r == nil {
		return
	}
	var t *template.Template = nil
	r.tmpl.Store(t)
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}) error {
	if r == nil {
		return ErrNilRenderer
	}
	x := r.tmpl.Load()
	if x == nil {
		return ErrNoTemplates
	}
	t := x.(*template.Template)
	return t.ExecuteTemplate(w, name, data)
}
