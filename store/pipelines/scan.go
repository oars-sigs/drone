package pipelines

import (
	"database/sql"
	"encoding/json"

	"github.com/oars-sigs/drone/model"

	"github.com/drone/drone/store/shared/db"
	"github.com/jmoiron/sqlx/types"
)

// helper function converts the Plugin structure to a set
// of named query parameters.
func toParams(p *model.Pipeline) map[string]interface{} {

	return map[string]interface{}{
		"pipeline_uuid":    p.UUID,
		"pipeline_name":    p.Name,
		"pipeline_repo":    p.Repo,
		"pipeline_slug":    p.Slug,
		"pipeline_ref":     p.Ref,
		"pipeline_sync":    p.Sync,
		"pipeline_content": p.Content,
		"pipeline_created": p.Created,
		"pipeline_updated": p.Updated,
	}
}

func encode(v interface{}) types.JSONText {
	raw, _ := json.Marshal(v)
	return types.JSONText(raw)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *model.Pipeline) error {
	err := scanner.Scan(
		&dest.UUID,
		&dest.Name,
		&dest.Repo,
		&dest.Slug,
		&dest.Ref,
		&dest.Content,
		&dest.Updated,
		&dest.Created,
		&dest.Sync,
	)
	return err
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*model.Pipeline, error) {
	defer rows.Close()

	pipelines := []*model.Pipeline{}
	for rows.Next() {
		pipeline := new(model.Pipeline)
		err := scanRow(rows, pipeline)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, pipeline)
	}
	return pipelines, nil
}
