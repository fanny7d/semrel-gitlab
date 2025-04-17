package actions

import (
	"fmt"
	"net/url"

	"github.com/fanny7d/semrel-gitlab/pkg/workflow"
	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
)

// CommitOptions GitLab 提交 API 的选项
// 参考: https://docs.gitlab.com/ee/api/commits.html#create-a-commit-with-multiple-files-and-actions
type CommitOptions struct {
	Branch        string         `json:"branch"`
	CommitMessage string         `json:"commit_message"`
	Actions       []CommitAction `json:"actions"`
	AuthorEmail   string         `json:"author_email"`
	AuthorName    string         `json:"author_name"`
	StartBranch   string         `json:"start_branch,omitempty"`
	StartSHA      string         `json:"start_sha,omitempty"`
	StartProject  string         `json:"start_project,omitempty"`
	Stats         bool           `json:"stats,omitempty"`
	Force         bool           `json:"force,omitempty"`
}

// CommitAction GitLab 提交操作
// 参考: https://docs.gitlab.com/ee/api/commits.html#create-a-commit-with-multiple-files-and-actions
type CommitAction struct {
	Action   string `json:"action"`
	FilePath string `json:"file_path"`
	Content  string `json:"content,omitempty"`
	Encoding string `json:"encoding,omitempty"`
}

// CommitResponse GitLab 提交 API 的响应
// 参考: https://docs.gitlab.com/ee/api/commits.html#create-a-commit-with-multiple-files-and-actions
type CommitResponse struct {
	ID             string   `json:"id"`
	ShortID        string   `json:"short_id"`
	Title          string   `json:"title"`
	AuthorName     string   `json:"author_name"`
	AuthorEmail    string   `json:"author_email"`
	AuthoredDate   string   `json:"authored_date"`
	CommitterName  string   `json:"committer_name"`
	CommitterEmail string   `json:"committer_email"`
	CommittedDate  string   `json:"committed_date"`
	CreatedAt      string   `json:"created_at"`
	Message        string   `json:"message"`
	ParentIDs      []string `json:"parent_ids"`
	Stats          struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
		Total     int `json:"total"`
	} `json:"stats"`
	Status string `json:"status"`
}

// CreateCommit 创建提交
// client: GitLab 客户端
// project: 项目路径
// options: 提交选项
func CreateCommit(client *gitlab.Client, project string, options *CommitOptions) (*CommitResponse, *gitlab.Response, error) {
	u := fmt.Sprintf("projects/%s/repository/commits", url.QueryEscape(project))

	req, err := client.NewRequest("POST", u, options, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create commit make request")
	}

	commit := new(CommitResponse)
	resp, err := client.Do(req, commit)
	if err != nil {
		return commit, resp, errors.Wrap(err, "create commit execute request")
	}

	return commit, resp, nil
}

// FileStatus 文件状态
type FileStatus string

const (
	// FileStatusCreated 文件已创建
	FileStatusCreated FileStatus = "created"
	// FileStatusUpdated 文件已更新
	FileStatusUpdated FileStatus = "updated"
	// FileStatusDeleted 文件已删除
	FileStatusDeleted FileStatus = "deleted"
	// FileStatusRenamed 文件已重命名
	FileStatusRenamed FileStatus = "renamed"
)

// FileStatusMap 文件状态映射
var FileStatusMap = map[FileStatus]string{
	FileStatusCreated: "create",
	FileStatusUpdated: "update",
	FileStatusDeleted: "delete",
	FileStatusRenamed: "move",
}

// GenerateCommitAction 生成提交操作
// status: 文件状态
// filePath: 文件路径
// content: 文件内容
func GenerateCommitAction(status FileStatus, filePath string, content string) CommitAction {
	return CommitAction{
		Action:   FileStatusMap[status],
		FilePath: filePath,
		Content:  content,
	}
}

// Commit 表示创建 GitLab 提交的操作
type Commit struct {
	client   *gitlab.Client
	project  string
	branch   string
	message  string
	commitID string
}

// Do 实现 Action 接口，执行创建提交的操作
func (action *Commit) Do() *workflow.ActionError {
	if action.commitID != "" {
		return nil
	}
	options := &gitlab.CreateCommitOptions{
		Branch:        &action.branch,
		CommitMessage: &action.message,
	}
	commit, _, err := action.client.Commits.CreateCommit(action.project, options)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "create commit"), true)
	}
	action.commitID = commit.ID
	return nil
}

// Undo 实现 Action 接口，撤销创建提交的操作
func (action *Commit) Undo() error {
	if action.commitID == "" {
		return nil
	}
	fmt.Printf("Commit cannot be rolled back: %s.\n", action.commitID)
	return nil
}

// BranchFunc 返回一个函数，用于获取分支名称
func (action *Commit) BranchFunc() func() string {
	return func() string {
		return action.branch
	}
}

// NewCommit 创建一个新的创建提交操作
func NewCommit(client *gitlab.Client, project string, branch string, message string) *Commit {
	return &Commit{
		client:  client,
		project: project,
		branch:  branch,
		message: message,
	}
}
