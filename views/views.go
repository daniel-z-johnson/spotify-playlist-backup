package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
)

type Template struct {
	htmlTemplate *template.Template
	logger       *slog.Logger
}

func (t *Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTemplate.Clone()
	if err != nil {
		t.logger.Error("Issue cloning template", slog.String("error", err.Error()))
		http.Error(w, "There was an error with processing the request.", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"errors": func() []string {
				return errMsgs
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		t.logger.Error("Issue executing template", slog.String("error", err.Error()))
		http.Error(w, "There was an error with processing the request.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func Must(t *Template, err error) *Template {
	if err != nil {
		// this will be called before the server starts if this fails there is no reason to start the server
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, logger *slog.Logger, fileNames ...string) (*Template, error) {
	htmlTemplate := template.New(fileNames[0])
	htmlTemplate.Funcs(
		template.FuncMap{
			"errors": func() []string {
				return nil
			},
		},
	)
	htmlTemplate, err := htmlTemplate.ParseFS(fs, fileNames...)
	if err != nil {
		logger.Error("Issue parsing template", slog.String("error", err.Error()))
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	return &Template{htmlTemplate: htmlTemplate, logger: logger}, nil
}

func errMessages(errs ...error) []string {
	errMessages := make([]string, len(errs))
	for _, err := range errs {
		errMessages = append(errMessages, err.Error())
	}
	return errMessages
}
