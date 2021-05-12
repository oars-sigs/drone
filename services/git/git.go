package git

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/oars-sigs/drone/model"

	"github.com/drone/drone/core"
	"github.com/drone/drone/handler/api/request"
	"github.com/drone/go-scm/scm"
)

// New returns a new CommitServiceFactory.
func New(client *scm.Client, renew core.Renewer) model.GitService {
	return &service{
		client: client,
		renew:  renew,
	}
}

type service struct {
	renew  core.Renewer
	client *scm.Client
}

func (s *service) userContext(ctx context.Context) (context.Context, error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.New("user info null")
	}
	err := s.renew.Renew(ctx, user, false)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, scm.TokenKey{}, &scm.Token{
		Token:   user.Token,
		Refresh: user.Refresh,
	}), nil

}

func (s *service) FindBranches(ctx context.Context, user *core.User, repo string) ([]*scm.Reference, error) {
	ctx, err := s.userContext(ctx)
	if err != nil {
		return nil, err
	}
	references, _, err := s.client.Git.ListBranches(ctx, repo, scm.ListOptions{})
	return references, err
}
func (s *service) FindTags(ctx context.Context, user *core.User, repo string) ([]*scm.Reference, error) {
	err := s.renew.Renew(ctx, user, false)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, scm.TokenKey{}, &scm.Token{
		Token:   user.Token,
		Refresh: user.Refresh,
	})
	references, _, err := s.client.Git.ListTags(ctx, repo, scm.ListOptions{})
	return references, nil
}

func (s *service) FindFile(ctx context.Context, user *core.User, repo, path, branch string) (*scm.Content, *scm.Response, error) {
	err := s.renew.Renew(ctx, user, false)
	if err != nil {
		return nil, nil, err
	}
	ctx = context.WithValue(ctx, scm.TokenKey{}, &scm.Token{
		Token:   user.Token,
		Refresh: user.Refresh,
	})
	return s.client.Contents.Find(ctx, repo, path, branch)
}

func (s *service) CreateFile(ctx context.Context, user *core.User, repo, path string, params *scm.ContentParams) error {
	err := s.renew.Renew(ctx, user, false)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, scm.TokenKey{}, &scm.Token{
		Token:   user.Token,
		Refresh: user.Refresh,
	})
	resp, err := s.client.Contents.Create(ctx, repo, path, params)
	if err != nil {
		return err
	}
	if resp.Status > 300 {
		data, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(data))
	}
	return err
}

func (s *service) UpdateFile(ctx context.Context, user *core.User, repo, path string, params *scm.ContentParams) error {
	err := s.renew.Renew(ctx, user, false)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, scm.TokenKey{}, &scm.Token{
		Token:   user.Token,
		Refresh: user.Refresh,
	})
	_, err = s.client.Contents.Update(ctx, repo, path, params)
	return err
}
