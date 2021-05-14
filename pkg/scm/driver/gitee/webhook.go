package gitee

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
)

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (scm.Webhook, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}
	var hook scm.Webhook
	switch req.Header.Get("X-Gitee-Event") {
	case "Push Hook", "Tag Push Hook":
		hook, err = s.parsePushHook(data)
	case "Merge Request Hook":
		hook, err = s.parsePullRequestHook(data)
	default:
		return nil, scm.ErrUnknownEvent
	}
	if err != nil {
		return nil, err
	}

	token, err := fn(hook)
	if err != nil {
		return hook, err
	} else if token == "" {
		return hook, nil
	}

	if token != req.Header.Get("X-Gitee-Token") {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePushHook(data []byte) (scm.Webhook, error) {
	dst := new(PushEvent)
	err := json.Unmarshal(data, dst)
	var commits []scm.Commit
	for _, c := range dst.Commits {
		commits = append(commits,
			scm.Commit{
				Sha:     c.Id,
				Message: c.Message,
				Link:    c.Url,
				Author: scm.Signature{
					Login: c.Author.Login,
					Email: c.Author.Email,
					Name:  c.Author.Name,
					Date:  c.Timestamp,
				},
				Committer: scm.Signature{
					Login: c.Committer.Login,
					Email: c.Committer.Email,
					Name:  c.Committer.Name,
					Date:  c.Timestamp,
				},
			})
	}
	return &scm.PushHook{
		Ref: *dst.Ref,
		Commit: scm.Commit{
			Sha:     *dst.After,
			Message: dst.Commits[0].Message,
			Link:    *dst.Compare,
			Author: scm.Signature{
				Login: dst.Commits[0].Author.Login,
				Email: dst.Commits[0].Author.Email,
				Name:  dst.Commits[0].Author.Name,
				Date:  dst.Commits[0].Timestamp,
			},
			Committer: scm.Signature{
				Login: dst.Commits[0].Committer.Login,
				Email: dst.Commits[0].Committer.Email,
				Name:  dst.Commits[0].Committer.Name,
				Date:  dst.Commits[0].Timestamp,
			},
		},
		Repo:    *convertHookRepository(dst.Repository),
		Sender:  *convertUser(dst.Sender),
		Commits: commits,
	}, err
}

func (s *webhookService) parsePullRequestHook(data []byte) (scm.Webhook, error) {
	dst := new(PullRequestEvent)
	err := json.Unmarshal(data, dst)
	return &scm.PullRequestHook{
		Action: convertAction(*dst.Action),
		PullRequest: scm.PullRequest{
			Number: dst.PullRequest.Number,
			Title:  dst.PullRequest.Title,
			Body:   dst.PullRequest.Body,
			Closed: dst.PullRequest.State == "closed",
			Author: scm.User{
				Login:  dst.PullRequest.User.Login,
				Email:  dst.PullRequest.User.Email,
				Avatar: dst.PullRequest.User.Avatar,
			},
			Merged: dst.PullRequest.Merged,
			// Created: nil,
			// Updated: nil,
			Source: dst.PullRequest.Head.Ref,
			Target: dst.PullRequest.Base.Ref,
			Link:   dst.PullRequest.HtmlUrl,
			Fork:   dst.PullRequest.Head.Repo.FullName,
			Ref:    fmt.Sprintf("refs/pull/%d/head", dst.PullRequest.Number),
			Sha:    dst.PullRequest.Head.Sha,
		},
		Repo:   *convertHookRepository(dst.Repository),
		Sender: *convertUser(dst.Sender),
	}, err
}

func convertHookRepository(src *ProjectHook) *scm.Repository {
	return &scm.Repository{
		ID:        strconv.Itoa(src.Id),
		Namespace: src.Owner.Login,
		Name:      src.Name,
		Perm: &scm.Perm{
			Pull: true,
			Push: true,
		},
		Branch:   src.DefaultBranch,
		Private:  src.Private,
		Clone:    src.GitHttpUrl,
		CloneSSH: src.GitSshUrl,
		Link:     src.HtmlUrl,
	}
}

type PushEvent struct {
	Ref                *string         `json:"ref,omitempty"`
	Before             *string         `json:"before,omitempty"`
	After              *string         `json:"after,omitempty"`
	TotalCommitsCount  int64           `json:"total_commits_count,omitempty"`
	CommitsMoreThanTen *bool           `json:"commits_more_than_ten,omitempty"`
	Created            *bool           `json:"created,omitempty"`
	Deleted            *bool           `json:"deleted,omitempty"`
	Compare            *string         `json:"compare,omitempty"`
	Commits            []CommitHook    `json:"commits,omitempty"`
	HeadCommit         *CommitHook     `json:"head_commit,omitempty"`
	Repository         *ProjectHook    `json:"repository,omitempty"`
	Project            *ProjectHook    `json:"project,omitempty"`
	UserID             int64           `json:"user_id,omitempty"`
	UserName           *string         `json:"user_name,omitempty"`
	User               *user           `json:"user,omitempty"`
	Pusher             *user           `json:"pusher,omitempty"`
	Sender             *user           `json:"sender,omitempty"`
	Enterprise         *EnterpriseHook `json:"enterprise,omitempty"`
	HookName           *string         `json:"hook_name,omitempty"`
	Password           *string         `json:"password,omitempty"`
}

type RepoInfo struct {
	Project    *ProjectHook `json:"project,omitempty"`
	Repository *ProjectHook `json:"repository,omitempty"`
}

type PullRequestEvent struct {
	Action         *string          `json:"action,omitempty"`
	ActionDesc     *string          `json:"action_desc,omitempty"`
	PullRequest    *PullRequestHook `json:"pull_request,omitempty"`
	Number         int64            `json:"number,omitempty"`
	IID            int64            `json:"iid,omitempty"`
	Title          *string          `json:"title,omitempty"`
	Body           *string          `json:"body,omitempty"`
	State          *string          `json:"state,omitempty"`
	MergeStatus    *string          `json:"merge_status,omitempty"`
	MergeCommitSha *string          `json:"merge_commit_sha,omitempty"`
	URL            *string          `json:"url,omitempty"`
	SourceBranch   *string          `json:"source_branch,omitempty"`
	SourceRepo     *RepoInfo        `json:"source_repo,omitempty"`
	TargetBranch   *string          `json:"target_branch,omitempty"`
	TargetRepo     *RepoInfo        `json:"target_repo,omitempty"`
	Project        *ProjectHook     `json:"project,omitempty"`
	Repository     *ProjectHook     `json:"repository,omitempty"`
	Author         *user            `json:"author,omitempty"`
	UpdatedBy      *user            `json:"updated_by,omitempty"`
	Sender         *user            `json:"sender,omitempty"`
	TargetUser     *user            `json:"target_user,omitempty"`
	Enterprise     *EnterpriseHook  `json:"enterprise,omitempty"`
	HookName       *string          `json:"hook_name,omitempty"`
	Password       *string          `json:"password,omitempty"`
}

type TagPushEvent struct {
	Action *string `json:"action,omitempty"`
}

type LabelHook struct {
	Id    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

// EnterpriseHook : 企业信息
type EnterpriseHook struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

// CommitHook : git commit 中的信息
type CommitHook struct {
	Id        string    `json:"id,omitempty"`
	TreeId    string    `json:"tree_id,omitempty"`
	ParentIds []string  `json:"parent_ids,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Url       string    `json:"url,omitempty"`
	Author    *user     `json:"author,omitempty"`
	Committer *user     `json:"committer,omitempty"`
	Distinct  bool      `json:"distinct,omitempty"`
	Added     []string  `json:"added,omitempty"`
	Removed   []string  `json:"removed,omitempty"`
	Modified  []string  `json:"modified,omitempty"`
}

// MilestoneHook : 里程碑信息
type MilestoneHook struct {
	Id           int       `json:"id,omitempty"`
	HtmlUrl      string    `json:"html_url,omitempty"`
	Number       int       `json:"number,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	OpenIssues   int       `json:"open_issues,omitempty"`
	ClosedIssues int       `json:"closed_issues,omitempty"`
	State        string    `json:"state,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	DueOn        string    `json:"due_on,omitempty"`
}

// IssueHook : issue 信息
type IssueHook struct {
	Id            int            `json:"id,omitempty"`
	HtmlUrl       string         `json:"html_url,omitempty"`
	Number        string         `json:"number,omitempty"`
	Title         string         `json:"title,omitempty"`
	User          *user          `json:"user,omitempty"`
	Labels        []LabelHook    `json:"labels,omitempty"`
	State         string         `json:"state,omitempty"`
	StateName     string         `json:"state_name,omitempty"`
	TypeName      string         `json:"type_name,omitempty"`
	Assignee      *user          `json:"assignee,omitempty"`
	Collaborators []user         `json:"collaborators,omitempty"`
	Milestone     *MilestoneHook `json:"milestone,omitempty"`
	Comments      int            `json:"comments,omitempty"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
	UpdatedAt     time.Time      `json:"updated_at,omitempty"`
	Body          string         `json:"body,omitempty"`
}

// ProjectHook : project 信息
type ProjectHook struct {
	Id              int    `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Path            string `json:"path,omitempty"`
	FullName        string `json:"full_name,omitempty"`
	Owner           *user  `json:"owner,omitempty"`
	Private         bool   `json:"private,omitempty"`
	HtmlUrl         string `json:"html_url,omitempty"`
	Url             string `json:"url,omitempty"`
	Description     string `json:"description,omitempty"`
	Fork            bool   `json:"fork,omitempty"`
	PushedAt        string `json:"pushed_at,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	SshUrl          string `json:"ssh_url,omitempty"`
	GitUrl          string `json:"git_url,omitempty"`
	CloneUrl        string `json:"clone_url,omitempty"`
	SvnUrl          string `json:"svn_url,omitempty"`
	GitHttpUrl      string `json:"git_http_url,omitempty"`
	GitSshUrl       string `json:"git_ssh_url,omitempty"`
	GitSvnUrl       string `json:"git_svn_url,omitempty"`
	Homepage        string `json:"homepage,omitempty"`
	StargazersCount int    `json:"stargazers_count,omitempty"`
	WatchersCount   int    `json:"watchers_count,omitempty"`
	ForksCount      int    `json:"forks_count,omitempty"`
	Language        string `json:"language,omitempty"`

	HasIssues bool   `json:"has_issues,omitempty"`
	HasWiki   bool   `json:"has_wiki,omitempty"`
	HasPage   bool   `json:"has_pages,omitempty"`
	License   string `json:"license,omitempty"`

	OpenIssuesCount int    `json:"open_issues_count,omitempty"`
	DefaultBranch   string `json:"default_branch,omitempty"`
	Namespace       string `json:"namespace,omitempty"`

	NameWithNamespace string `json:"name_with_namespace,omitempty"`
	PathWithNamespace string `json:"path_with_namespace,omitempty"`
}

// BranchHook : 分支信息
type BranchHook struct {
	Label string       `json:"label,omitempty"`
	Ref   string       `json:"ref,omitempty"`
	Sha   string       `json:"sha,omitempty"`
	User  *user        `json:"user,omitempty"`
	Repo  *ProjectHook `json:"repo,omitempty"`
}

// PullRequestHook : PR 信息
type PullRequestHook struct {
	Id                 int            `json:"id,omitempty"`
	Number             int            `json:"number,omitempty"`
	State              string         `json:"state,omitempty"`
	HtmlUrl            string         `json:"html_url,omitempty"`
	DiffUrl            string         `json:"diff_url,omitempty"`
	PatchUrl           string         `json:"patch_url,omitempty"`
	Title              string         `json:"title,omitempty"`
	Body               string         `json:"body,omitempty"`
	Labels             []LabelHook    `json:"labels,omitempty"`
	CreatedAt          string         `json:"created_at,omitempty"`
	UpdatedAt          string         `json:"updated_at,omitempty"`
	ClosedAt           string         `json:"closed_at,omitempty"`
	MergedAt           string         `json:"merged_at,omitempty"`
	MergeCommitSha     string         `json:"merge_commit_sha,omitempty"`
	MergeReferenceName string         `json:"merge_reference_name,omitempty"`
	User               *user          `json:"user,omitempty"`
	Assignee           *user          `json:"assignee,omitempty"`
	Assignees          []user         `json:"assignees,omitempty"`
	Tester             []user         `json:"tester,omitempty"`
	Testers            []user         `json:"testers,omitempty"`
	NeedTest           bool           `json:"need_test,omitempty"`
	NeedReview         bool           `json:"need_review,omitempty"`
	Milestone          *MilestoneHook `json:"milestone,omitempty"`
	Head               *BranchHook    `json:"head,omitempty"`
	Base               *BranchHook    `json:"base,omitempty"`
	Merged             bool           `json:"merged,omitempty"`
	Mergeable          bool           `json:"mergeable,omitempty"`
	MergeStatus        string         `json:"merge_status,omitempty"`
	UpdatedBy          *user          `json:"updated_by,omitempty"`
	Comments           int            `json:"comments,omitempty"`
	Commits            int            `json:"commits,omitempty"`
	Additions          int            `json:"additions,omitempty"`
	Deletions          int            `json:"deletions,omitempty"`
	ChangedFiles       int            `json:"changed_files,omitempty"`
}

func convertAction(src string) (action scm.Action) {
	switch src {
	case "create", "created":
		return scm.ActionCreate
	case "delete", "deleted":
		return scm.ActionDelete
	case "update", "updated", "edit", "edited":
		return scm.ActionUpdate
	case "open", "opened":
		return scm.ActionOpen
	case "reopen", "reopened":
		return scm.ActionReopen
	case "close", "closed":
		return scm.ActionClose
	case "label", "labeled":
		return scm.ActionLabel
	case "unlabel", "unlabeled":
		return scm.ActionUnlabel
	case "merge", "merged":
		return scm.ActionMerge
	case "synchronize", "synchronized":
		return scm.ActionSync
	default:
		return
	}
}
