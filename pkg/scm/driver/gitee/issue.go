package gitee

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type issueService struct {
	client *wrapper
}

func (s *issueService) Find(ctx context.Context, repo string, number int) (*scm.Issue, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) FindComment(ctx context.Context, repo string, index, id int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) List(ctx context.Context, repo string, _ scm.IssueListOptions) ([]*scm.Issue, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) ListComments(ctx context.Context, repo string, index int, _ scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) Create(ctx context.Context, repo string, input *scm.IssueInput) (*scm.Issue, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) CreateComment(ctx context.Context, repo string, index int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) DeleteComment(ctx context.Context, repo string, index, id int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *issueService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *issueService) Lock(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *issueService) Unlock(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}
