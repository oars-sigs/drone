package gitee

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/branches/%s", repo, name)
	out := new(branch)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertBranch(out), res, err
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/commits/%s", repo, scm.TrimRef(ref))
	out := new(commit)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommitInfo(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/tags", repo)
	out := []*tag{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if err != nil {
		return nil, res, err
	}
	for _, t := range out {
		if t.Name == scm.TrimRef(name) {
			return convertTag(t), res, err
		}

	}
	return nil, res, scm.ErrNotFound
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/branches", repo)
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertBranchList(out), res, err
}

func (s *gitService) ListTags(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/tags", repo)
	out := []*tag{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertTagList(out), res, err
}

func (s *gitService) ListCommits(ctx context.Context, repo string, _ scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type (
	branch struct {
		Name   string `json:"name"`
		Commit commit `json:"commit"`
	}

	commit struct {
		Sha       string     `json:"sha"`
		URL       string     `json:"url"`
		Commit    commitInfo `json:"commit"`
		Author    user       `json:"author"`
		Committer user       `json:"committer"`
	}

	commitInfo struct {
		Message   string    `json:"message"`
		URL       string    `json:"url"`
		Author    signature `json:"author"`
		Committer signature `json:"committer"`
	}

	signature struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	tag struct {
		Name    string `json:"name"`
		Message string `json:"message"`
		Commit  commit `json:"commit"`
	}
)

func convertBranch(src *branch) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(src.Name),
		Path: scm.ExpandRef(src.Name, "refs/heads/"),
		Sha:  src.Commit.Sha,
	}
}

func convertCommitInfo(src *commit) *scm.Commit {
	return &scm.Commit{
		Sha:       src.Sha,
		Link:      src.URL,
		Message:   src.Commit.Message,
		Author:    convertUserSignature(src.Author),
		Committer: convertUserSignature(src.Committer),
	}
}

func convertTag(src *tag) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(src.Name),
		Path: src.Name,
		Sha:  src.Commit.Sha,
	}
}

func convertTagList(src []*tag) []*scm.Reference {
	var dst []*scm.Reference
	for _, v := range src {
		dst = append(dst, convertTag(v))
	}
	return dst
}

func convertBranchList(src []*branch) []*scm.Reference {
	dst := []*scm.Reference{}
	for _, v := range src {
		dst = append(dst, convertBranch(v))
	}
	return dst
}

func convertSignature(src signature) scm.Signature {
	return scm.Signature{
		Login: src.Name,
		Email: src.Email,
		Name:  src.Name,
	}
}

func convertUserSignature(src user) scm.Signature {
	return scm.Signature{
		Login:  src.Login,
		Email:  src.Email,
		Name:   src.Name,
		Avatar: src.Avatar,
	}
}
