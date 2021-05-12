package model

import (
	"context"

	"github.com/drone/drone/core"
	"github.com/drone/go-scm/scm"
)

// GitService provides access to the commit history from
// the external source code management service (e.g. Gitlab).
type GitService interface {
	//
	CreateFile(ctx context.Context, user *core.User, repo, path string, params *scm.ContentParams) error

	//
	UpdateFile(ctx context.Context, user *core.User, repo, path string, params *scm.ContentParams) error

	//
	FindFile(ctx context.Context, user *core.User, repo, path, branch string) (*scm.Content, *scm.Response, error)

	//
	FindBranches(ctx context.Context, user *core.User, repo string) ([]*scm.Reference, error)

	//
	FindTags(ctx context.Context, user *core.User, repo string) ([]*scm.Reference, error)
}
