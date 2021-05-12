package pipelines

import (
	"context"
	"database/sql"
	"errors"

	"github.com/oars-sigs/drone/model"

	"github.com/drone/drone/core"
	"github.com/drone/drone/store/shared/db"
)

type pipelineStore struct {
	db *db.DB
}

func New(db *db.DB) model.PipelineStore {
	return &pipelineStore{db: db}
}

// Get returns a of pipeline from the datastore
func (s *pipelineStore) GetPipeline(ctx context.Context, slug, ref, configPath string) (*model.Pipeline, bool, error) {
	out := &model.Pipeline{
		Slug: slug,
		Ref:  ref,
	}
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := toParams(out)
		query, args, err := binder.BindNamed(queryBySlugRef, params)
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
	}
	return out, true, err
}

func (s *pipelineStore) FindPipelineBySlug(ctx context.Context, slug string) ([]*model.Pipeline, error) {
	var out []*model.Pipeline
	err := s.db.View(func(queryer db.Queryer, binder db.Binder) error {
		params := toParams(&model.Pipeline{Slug: slug})
		query, args, err := binder.BindNamed(queryBySlug, params)
		if err != nil {
			return err
		}
		rows, err := queryer.Query(query, args...)
		if err != nil {
			return err
		}
		out, err = scanRows(rows)
		return err
	})
	return out, err
}

func (s *pipelineStore) CreatePipeline(ctx context.Context, pipe *model.Pipeline) error {
	_, isExist, err := s.GetPipeline(ctx, pipe.Slug, pipe.Ref, "")
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("pipeline has existed")
	}
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(pipe)
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

func (s *pipelineStore) UpdatePipeline(ctx context.Context, pipe *model.Pipeline) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(pipe)
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

func (s *pipelineStore) Find(ctx context.Context, r *core.ConfigArgs) (*core.Config, error) {
	pipe, isExist, err := s.GetPipeline(ctx, r.Repo.Slug, r.Build.Ref, "")
	if err != nil {
		return nil, err
	}
	if isExist {
		return &core.Config{
			Kind: "pipeline",
			Data: pipe.Content,
		}, nil
	}
	pipe, isExist, err = s.GetPipeline(ctx, r.Repo.Slug, "default", "")
	if err != nil {
		return nil, err
	}
	if isExist {
		return &core.Config{
			Kind: "pipeline",
			Data: pipe.Content,
		}, nil
	}
	return nil, errors.New(r.Repo.Slug + r.Build.Ref + "pipeline not found")
}

const queryBySlugRef = `
SELECT * FROM tpipe_pipelines 
WHERE pipeline_slug=:pipeline_slug AND pipeline_ref=:pipeline_ref
`

const queryBySlug = `
SELECT * FROM tpipe_pipelines 
WHERE pipeline_slug=:pipeline_slug
`

const stmtInsert = `
INSERT INTO tpipe_pipelines (
 pipeline_uuid
,pipeline_name
,pipeline_repo
,pipeline_slug
,pipeline_ref
,pipeline_content
,pipeline_created
,pipeline_updated
,pipeline_sync
) VALUES (
 :pipeline_uuid
,:pipeline_name
,:pipeline_repo
,:pipeline_slug
,:pipeline_ref
,:pipeline_content
,:pipeline_created
,:pipeline_updated
,:pipeline_sync
)
`

const stmtUpdate = `
UPDATE tpipe_pipelines SET
pipeline_name=:pipeline_name
,pipeline_content=:pipeline_content
,pipeline_updated=:pipeline_updated
,pipeline_sync=:pipeline_sync
WHERE pipeline_slug=:pipeline_slug AND pipeline_ref=:pipeline_ref
`
