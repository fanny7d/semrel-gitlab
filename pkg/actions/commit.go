package actions

import (
	"io/ioutil"

	"github.com/fanny7d/semrel-gitlab/pkg/workflow"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	git "gopkg.in/src-d/go-git.v4"
)

var (
	// fileStatusMap 将 git 状态码映射到 GitLab 提交操作
	// 在这个上下文中只支持创建和更新操作
	fileStatusMap = map[git.StatusCode]gitlab.FileActionValue{
		git.Untracked: gitlab.FileCreate,
		git.Modified:  gitlab.FileUpdate,
	}
)

// Commit 表示 GitLab 提交操作
type Commit struct {
	client  *gitlab.Client
	project string
	files   []string
	message string
	branch  string
	commit  *gitlab.Commit
}

// Do 实现 Action 接口，执行提交操作
func (action *Commit) Do() *workflow.ActionError {
	if action.commit != nil {
		return nil
	}
	commitActions, err := getCommitActionsForFiles(action.files)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "commit"), false)
	}
	if len(commitActions) == 0 {
		return workflow.NewActionError(errors.New("No changes to commit"), false)
	}
	options := &gitlab.CreateCommitOptions{
		Branch:        &action.branch,
		CommitMessage: &action.message,
		Actions:       commitActions,
	}
	commit, _, err := action.client.Commits.CreateCommit(action.project, options)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "commit"), true)
	}
	action.commit = commit
	return nil
}

// Undo 实现 Action 接口，撤销提交操作
func (action *Commit) Undo() error {
	if action.commit == nil {
		return nil
	}
	// TODO: 实现撤销操作
	return nil
}

// RefFunc 返回一个函数，用于获取提交的引用
func (action *Commit) RefFunc() func() string {
	return func() string {
		if action.commit == nil {
			return ""
		}
		return action.commit.ID
	}
}

// NewCommit 创建一个新的提交操作
func NewCommit(client *gitlab.Client, project string, files []string, message string, branch string) *Commit {
	return &Commit{
		client:  client,
		project: project,
		files:   files,
		message: message,
		branch:  branch,
	}
}

// getCommitActionsForFiles 为文件列表生成提交操作
func getCommitActionsForFiles(files []string) ([]*gitlab.CommitActionOptions, error) {
	actions := []*gitlab.CommitActionOptions{}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Wrapf(err, "read file %s", file)
		}
		action := &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(gitlab.FileUpdate),
			FilePath: gitlab.String(file),
			Content:  gitlab.String(string(content)),
		}
		actions = append(actions, action)
	}
	return actions, nil
}
