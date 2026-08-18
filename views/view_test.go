package views

import (
	"errors"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func parseTemplate(t *testing.T, source string) *Template {
	t.Helper()
	return Must(ParseFS(testLogger(), fstest.MapFS{
		"test.gohtml": {Data: []byte(source)},
	}, "test.gohtml"))
}

func TestExecuteEscapesHTML(t *testing.T) {
	tpl := parseTemplate(t, `{{.}}`)
	recorder := httptest.NewRecorder()

	tpl.Execute(recorder, `<script>alert("xss")</script>`, nil)

	if recorder.Code != 200 {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if got, want := recorder.Body.String(), `&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;`; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q", got)
	}
	if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("X-Content-Type-Options = %q", got)
	}
}

type failingData struct{}

func (failingData) Value() (string, error) {
	return "", errors.New("render failed")
}

func TestExecuteDoesNotWritePartialTemplateOnError(t *testing.T) {
	tpl := parseTemplate(t, `partial output{{.Value}}`)
	recorder := httptest.NewRecorder()

	tpl.Execute(recorder, failingData{}, nil)

	if recorder.Code != 500 {
		t.Fatalf("status = %d, want 500", recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), "partial output") {
		t.Fatalf("body contains partial template output: %q", recorder.Body.String())
	}
}

func TestExecuteUsesErrorsForCurrentRender(t *testing.T) {
	tpl := parseTemplate(t, `{{range errors}}{{.}};{{end}}`)
	recorder := httptest.NewRecorder()

	tpl.Execute(recorder, nil, []error{errors.New("first"), errors.New("second")})

	if got, want := recorder.Body.String(), "first;second;"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestParseFSRequiresPattern(t *testing.T) {
	_, err := ParseFS(testLogger(), fstest.MapFS{})
	if err == nil || !strings.Contains(err.Error(), "at least one template pattern") {
		t.Fatalf("error = %v, want missing-pattern error", err)
	}
}

func TestMustPanicsWithParseError(t *testing.T) {
	want := errors.New("parse failed")
	defer func() {
		if got := recover(); got != want {
			t.Fatalf("panic = %v, want %v", got, want)
		}
	}()

	Must(nil, want)
}
