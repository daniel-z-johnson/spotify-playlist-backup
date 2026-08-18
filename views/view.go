package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
	log     *slog.Logger
}

func (t *Template) Execute(w http.ResponseWriter, data any, renderErrors []error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		t.log.Error("clone template failed", slog.String("error", err.Error()))
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(
		template.FuncMap{
			"errors": func() []error {
				return renderErrors
			},
		},
	)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		t.log.Error("executing template", slog.String("error", err.Error()))
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if _, err := io.Copy(w, &buf); err != nil {
		t.log.Error("writing rendered template", slog.String("error", err.Error()))
	}
}

func Must(tmpl *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return tmpl
}

func ParseFS(log *slog.Logger, fsys fs.FS, patterns ...string) (*Template, error) {
	if len(patterns) == 0 {
		return nil, errors.New("views: at least one template pattern is required")
	}
	if log == nil {
		return nil, errors.New("views: logger is required")
	}
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			"errors": func() []error {
				return nil
			},
		},
	)
	tpl, err := tpl.ParseFS(fsys, patterns...)
	if err != nil {
		return nil, err
	}
	return &Template{
		htmlTpl: tpl,
		log:     log,
	}, nil
}
