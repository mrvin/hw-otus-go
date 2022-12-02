package httpserver

import (
	"embed"
	"html/template"
	"path/filepath"
	"time"
)

//go:embed templates/*
var files embed.FS

type templateLoader struct {
	templates map[string]*template.Template
}

func newTemplateLoader() *templateLoader {
	return &templateLoader{make(map[string]*template.Template)}
}

func dateFormat(layout string, d time.Time) string {
	return d.Format(layout)
}

func (t *templateLoader) LoadTemplates(templatesDir string) {
	t.templates["form-user.html"] = template.Must(template.ParseFS(files,
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, "form-user.html")))
	t.templates["form-event.html"] = template.Must(template.ParseFS(files,
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, "form-event.html")))
	t.templates["text.html"] = template.Must(template.ParseFS(files,
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, "text.html")))
	t.templates["list-users.html"] = template.Must(template.ParseFS(files,
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, "list-users.html")))
	t.templates["list-events.html"] = template.Must(template.ParseFS(files,
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, "list-events.html")))
	t.templates["user.html"] = template.Must(template.ParseFS(files,
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, "user.html")))

	name := "base.html"
	t.templates["event.html"] = template.Must(template.New(name).
		Funcs(template.FuncMap{"dateFormat": dateFormat}).
		ParseFS(files, filepath.Join(templatesDir, name),
			filepath.Join(templatesDir, "event.html")))
}
