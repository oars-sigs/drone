package gitee

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	endpoint := fmt.Sprintf("api/v5/repos/%s/contents/%s?ref=%s", repo, url.QueryEscape(path), scm.TrimRef(ref))
	out := new(content)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	raw, berr := base64.StdEncoding.DecodeString(out.Content)
	if berr != nil {
		raw, berr = base64.RawStdEncoding.DecodeString(out.Content)
		if berr != nil {
			return nil, res, err
		}
	}
	return &scm.Content{
		Path: out.Path,
		Data: raw,
	}, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) Delete(ctx context.Context, repo, path, ref string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, _ scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	endpoint := fmt.Sprintf("api/v5/repos/%s/contents/%s?ref=%s", repo, url.QueryEscape(path), scm.TrimRef(ref))
	out := []*content{}
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertContentInfoList(out), res, err
}

type content struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	SHA      string `json:"sha"`
}

func convertContentInfoList(from []*content) []*scm.ContentInfo {
	to := []*scm.ContentInfo{}
	for _, v := range from {
		to = append(to, convertContentInfo(v))
	}
	return to
}

func convertContentInfo(from *content) *scm.ContentInfo {
	to := &scm.ContentInfo{Path: from.Path}
	switch from.Type {
	case "file":
		to.Kind = scm.ContentKindFile
	case "dir":
		to.Kind = scm.ContentKindDirectory
	case "symlink":
		to.Kind = scm.ContentKindSymlink
	case "submodule":
		to.Kind = scm.ContentKindGitlink
	default:
		to.Kind = scm.ContentKindUnsupported
	}
	return to
}
