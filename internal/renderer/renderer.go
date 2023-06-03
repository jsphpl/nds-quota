package renderer

import (
	"bytes"
	"html/template"
	"path"
)

type Params struct {
	Query    map[string]string
	Info     string
	Error    string
	Token    string
	Title    string
	ShowForm bool
}

type Renderer struct {
	template *template.Template
	title    string
}

func NewRenderer(templatePath string) (*Renderer, error) {
	var err error
	r := Renderer{
		title: "ndsquota",
	}
	r.template, err = template.ParseGlob(path.Join(templatePath, "*.html"))

	return &r, err
}

func (r *Renderer) Render(params Params, showForm bool) (string, error) {
	if params.Title == "" {
		params.Title = r.title
	}

	params.ShowForm = showForm

	w := bytes.NewBuffer(nil)
	err := r.template.ExecuteTemplate(w, "base", params)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}
