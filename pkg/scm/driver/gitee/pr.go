package gitee

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(context.Context, string, int) (*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) FindComment(context.Context, string, int, int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) List(context.Context, string, scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListComments(context.Context, string, int, scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListChanges(context.Context, string, int, scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) Create(context.Context, string, *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) CreateComment(context.Context, string, int, *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Merge(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Close(context.Context, string, int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}
