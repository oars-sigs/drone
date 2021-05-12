package model

import (
	"context"
)

type TemplateJSON struct {
	Total     int        `json:"total"`
	Templates []Template `json:"templates"`
}

type Template struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Format  string `json:"format"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Created int64  `json:"-"`
	Updated int64  `json:"-"`
}

type TemplateStore interface {

	//Get returns a list of templates from the datastore.
	GetTemplate() ([]*Template, error)

	//Create persists a new template to the datastore.
	CreateTemplate(ctx context.Context, tmpl *Template) error

	//Put a template to the datastore.
	PutTemplate(ctx context.Context, tmpl *Template) error

	//Delete a template from the datastore.
	DeleteTemplate(ctx context.Context, uuid string) error

	//Find a template from the datastore.
	FindTemplate(ctx context.Context, uuid string) (*Template, bool, error)
}
