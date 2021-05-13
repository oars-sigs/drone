package gitee

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	out := new(user)
	res, err := s.client.do(ctx, "GET", "api/v5/user", nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/users/%s", login)
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	user, res, err := s.Find(ctx)
	return user.Email, res, err
}

type user struct {
	ID     int    `json:"id"`
	Login  string `json:"login"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar_url"`
}

func convertUser(src *user) *scm.User {
	return &scm.User{
		Login:  src.Login,
		Avatar: src.Avatar,
		Email:  src.Email,
		Name:   src.Name,
	}
}
