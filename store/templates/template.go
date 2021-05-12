package templates

import (
	"context"
	"database/sql"
	"errors"

	"github.com/drone/drone/store/shared/db"
	"github.com/oars-sigs/drone/model"
	"github.com/sirupsen/logrus"
)

func New(db *db.DB) model.TemplateStore {
	return &tpmlStore{
		db: db,
	}
}

type tpmlStore struct {
	db *db.DB
}

//Get returns a list of templates from the datastore.
func (s *tpmlStore) GetTemplate() ([]*model.Template, error) {
	var out []*model.Template
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		rows, err := queryer.Query(queryAll)
		if err != nil {
			return err
		}
		out, err = scanRows(rows)
		return err
	})
	return out, err
}

//Create persists a new template to the datastore.
func (s *tpmlStore) CreateTemplate(ctx context.Context, tmpl *model.Template) error {
	_, isExist, err := s.FindTemplateName(ctx, tmpl.Name)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("template name is exist")
	}
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(tmpl)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		return err
	})
}

//Put a template to the datastore.
func (s *tpmlStore) PutTemplate(ctx context.Context, tmpl *model.Template) error {
	t, isExist, err := s.FindTemplateName(ctx, tmpl.Name)
	if err != nil {
		return err
	}
	if isExist && t.UUID != tmpl.UUID {
		return errors.New("template name is exist")
	}
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParam(tmpl)
		stmt, args, err := binder.BindNamed(stmtUpdate, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		return err
	})
}

//Delete a template from the datastore.
func (s *tpmlStore) DeleteTemplate(ctx context.Context, uuid string) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		tmp := &model.Template{
			UUID: uuid,
		}
		params := toParam(tmp)
		stmt, args, err := binder.BindNamed(stmtDelete, params)
		if err != nil {
			return err
		}
		_, err = execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		return err
	})
}

//Find a template from the datastore.
func (s *tpmlStore) FindTemplate(ctx context.Context, uuid string) (*model.Template, bool, error) {
	out := &model.Template{
		UUID: uuid,
	}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := toParam(out)
		query, args, err := binder.BindNamed(queryByUuid, params)
		if err != nil {
			return err
		}
		row := queryer.QueryRow(query, args...)
		if err != nil {
			return err
		}

		err = scanRow(row, out)
		return err
	})
	logrus.Debug(out)
	return out, true, err
}

//Find a template from the datastore by name.
func (s *tpmlStore) FindTemplateName(ctx context.Context, name string) (*model.Template, bool, error) {
	out := &model.Template{
		Name: name,
	}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := toParam(out)
		query, args, err := binder.BindNamed(queryByName, params)
		if err != nil {
			return err
		}
		row := queryer.QueryRow(query, args...)
		if err != nil {
			return err
		}

		err = scanRow(row, out)
		return err
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	return out, true, err
}

const queryAll = `
SELECT * FROM tpipe_templates;

`

const stmtInsert = `
INSERT INTO tpipe_templates (
 template_uuid
,template_name
,template_format
,template_type
,template_content
,template_created
,template_updated
) VALUES (
 :template_uuid
,:template_name
,:template_format
,:template_type
,:template_content
,:template_created
,:template_updated
)
`

const stmtUpdate = `
UPDATE tpipe_templates
SET
template_uuid           = :template_uuid
,template_name          = :template_name
,template_format        = :template_format
,template_type          = :template_type
,template_content       = :template_content
,template_updated       = :template_updated
WHERE template_uuid     = :template_uuid
`

const stmtDelete = `
DELETE FROM tpipe_templates WHERE template_uuid = :template_uuid
`

const queryByUuid = `
SELECT * FROM tpipe_templates WHERE template_uuid = :template_uuid
`
const queryByName = `
SELECT * FROM tpipe_templates WHERE template_name = :template_name
`
