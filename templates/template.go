package templates

import (
	"html/template"

	"github.com/pkg/errors"
)

// Templates holds parsed a templates data
type Templates struct {
	Tpls *template.Template
}

// NewTemplates returns an instance of Templates
func NewTemplates() (*Templates, error) {
	path := "./templates/index.tmpl"
	t, err := template.New("index.tmpl").ParseFiles(path)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &Templates{
		Tpls: t,
	}, nil
}
