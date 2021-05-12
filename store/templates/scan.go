package templates

import (
	"database/sql"
	"time"

	"github.com/oars-sigs/drone/model"

	"github.com/drone/drone/store/shared/db"
	"github.com/google/uuid"
)

// helper function converts the Plugin structure to a set
// of named query parameters.
func toParams(t *model.Template) map[string]interface{} {
	uuid := uuid.New().String()
	created_time := time.Now().Unix()
	updated_time := time.Now().Unix()
	return map[string]interface{}{
		"template_uuid":    uuid,
		"template_name":    t.Name,
		"template_format":  t.Format,
		"template_type":    t.Type,
		"template_content": t.Content,
		"template_updated": updated_time,
		"template_created": created_time,
	}
}

func toParam(t *model.Template) map[string]interface{} {
	updated_time := time.Now().Unix()
	return map[string]interface{}{
		"template_uuid":    t.UUID,
		"template_name":    t.Name,
		"template_format":  t.Format,
		"template_type":    t.Type,
		"template_content": t.Content,
		"template_updated": updated_time,
	}
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *model.Template) error {
	err := scanner.Scan(
		&dest.UUID,
		&dest.Name,
		&dest.Format,
		&dest.Type,
		&dest.Content,
		&dest.Updated,
		&dest.Created,
	)
	return err
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*model.Template, error) {
	defer rows.Close()

	templates := []*model.Template{}
	for rows.Next() {
		template := new(model.Template)
		err := scanRow(rows, template)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}
