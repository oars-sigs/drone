package gitee

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"
)

type repositoryService struct {
	client *wrapper
}

func (s *repositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("/api/v5/repos/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out), res, err
}

func (s *repositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/hooks/%s", repo, id)
	out := new(hook)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertHook(out), res, err
}

func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	path := fmt.Sprintf("/api/v5/repos/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out).Perm, res, err
}

func (s *repositoryService) List(ctx context.Context, _ scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	path := "api/v5/user/repos?visibility=all&affiliation=owner%2C%20collaborator%2C%20organization_member&sort=full_name&direction=asc&page=1&per_page=100"
	out := []*repository{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertRepositoryList(out), res, err
}

func (s *repositoryService) ListHooks(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/hooks", repo)
	out := []*hook{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertHookList(out), res, err
}

func (s *repositoryService) ListStatus(context.Context, string, string, scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/hooks", repo)
	in := &hook{
		URL:                 input.Target,
		EncryptionType:      0,
		Password:            input.Secret,
		PushEvents:          true,
		TagPushEvents:       true,
		MergeRequestsEvents: true,
	}
	out := new(hook)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertHook(out), res, err
}

func (s *repositoryService) CreateStatus(context.Context, string, string, *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) UpdateHook(ctx context.Context, repo, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/hooks/%s", repo, id)
	in := &hook{
		URL:                 input.Target,
		EncryptionType:      0,
		Password:            input.Secret,
		PushEvents:          true,
		TagPushEvents:       true,
		MergeRequestsEvents: true,
	}
	out := new(hook)
	res, err := s.client.do(ctx, "PATCH", path, in, out)
	return convertHook(out), res, err
}

func (s *repositoryService) DeleteHook(ctx context.Context, repo string, id string) (*scm.Response, error) {
	path := fmt.Sprintf("api/v1/repos/%s/hooks/%s", repo, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

type (
	repository struct {
		ID            int       `json:"id"`
		Owner         user      `json:"owner"`
		Name          string    `json:"name"`
		FullName      string    `json:"full_name"`
		Private       bool      `json:"private"`
		Fork          bool      `json:"fork"`
		HTMLURL       string    `json:"html_url"`
		SSHURL        string    `json:"ssh_url"`
		DefaultBranch string    `json:"default_branch"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Permissions   perm      `json:"permissions"`
		Namespace     namespace `json:"namespace"`
	}

	perm struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	}

	namespace struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}

	hook struct {
		ID                  int    `json:"id"`
		EncryptionType      int    `json:"encryption_type"`
		Password            string `json:"password"`
		PushEvents          bool   `json:"push_events"`
		TagPushEvents       bool   `json:"tag_push_events"`
		IssuesEvents        bool   `json:"issues_events"`
		NoteEvents          bool   `json:"note_events"`
		MergeRequestsEvents bool   `json:"merge_requests_events"`
		URL                 string `json:"url"`
	}
)

func convertRepositoryList(src []*repository) []*scm.Repository {
	var dst []*scm.Repository
	for _, v := range src {
		dst = append(dst, convertRepository(v))
	}
	return dst
}

func convertRepository(src *repository) *scm.Repository {
	return &scm.Repository{
		ID:        strconv.Itoa(src.ID),
		Namespace: src.Namespace.Path,
		Name:      src.Name,
		Perm:      convertPerm(src.Permissions),
		Branch:    src.DefaultBranch,
		Private:   src.Private,
		Clone:     src.HTMLURL,
		CloneSSH:  src.SSHURL,
		Link:      strings.TrimSuffix(src.HTMLURL, ".git"),
	}
}

func convertPerm(src perm) *scm.Perm {
	return &scm.Perm{
		Push:  src.Push,
		Pull:  src.Pull,
		Admin: src.Admin,
	}
}

func convertHookList(src []*hook) []*scm.Hook {
	var dst []*scm.Hook
	for _, v := range src {
		dst = append(dst, convertHook(v))
	}
	return dst
}

func convertHook(from *hook) *scm.Hook {

	return &scm.Hook{
		ID:     strconv.Itoa(from.ID),
		Active: true,
		Target: from.URL,
		Events: convertHookEvent(from),
	}
}

func convertHookEvent(from *hook) []string {
	var events []string
	if from.MergeRequestsEvents {
		events = append(events, "pull_request")
	}
	if from.IssuesEvents {
		events = append(events, "issues")
	}
	if from.IssuesEvents {
		events = append(events, "issue_comment")
	}
	if from.PushEvents {
		events = append(events, "push")
	}
	return events
}
