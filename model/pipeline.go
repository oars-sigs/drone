package model

import (
	"context"

	"github.com/drone/drone/core"
)

type Pipeline struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Repo       string `json:"repo"`
	Slug       string `json:"slug"`
	Ref        string `json:"ref"`
	ConfigPath string `json:"config_path"`
	Content    string `json:"content"`
	Sync       int    `json:"sync"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
}

type PipelineStore interface {
	Find(ctx context.Context, r *core.ConfigArgs) (*core.Config, error)
	GetPipeline(ctx context.Context, slug, ref, configPath string) (*Pipeline, bool, error)
	UpdatePipeline(ctx context.Context, pipe *Pipeline) error
	CreatePipeline(ctx context.Context, pipe *Pipeline) error
}

type PipelineService interface {
	Find(ctx context.Context, r *core.ConfigArgs) (*core.Config, error)
}
