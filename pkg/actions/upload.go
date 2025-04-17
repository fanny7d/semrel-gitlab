package actions

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/fanny7d/semrel-gitlab/pkg/workflow"
	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
)

// Upload 表示 GitLab 文件上传操作
type Upload struct {
	client             *gitlab.Client
	project            string
	projectURL         *url.URL
	file               string
	projectFile        *gitlab.ProjectFile
	fullprojectFileURL *url.URL
}

// Do 实现 Action 接口，执行文件上传操作
func (action *Upload) Do() *workflow.ActionError {
	if action.projectFile != nil {
		return nil
	}

	file, err := os.Open(action.file)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "open file"), false)
	}
	defer file.Close()

	projectFile, resp, err := action.client.Projects.UploadFile(
		action.project,
		file,
		filepath.Base(action.file),
	)
	if err != nil {
		retry := false
		if resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(errors.Wrap(err, "upload file"), retry)
	}
	action.projectFile = projectFile

	action.fullprojectFileURL, err = url.Parse(action.projectURL.String() + action.projectFile.URL)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "parse url"), false)
	}
	return nil
}

// Undo 实现 Action 接口，撤销文件上传操作
func (action *Upload) Undo() error {
	if action.projectFile == nil {
		return nil
	}
	fmt.Printf("File upload cannot be rolled back: %s.\n", action.file)
	return nil
}

// MDLinkFunc 返回一个函数，用于获取 Markdown 格式的文件链接
func (action *Upload) MDLinkFunc() func() string {
	return func() string {
		if action.projectFile == nil {
			return ""
		}
		return fmt.Sprintf("[%s](%s)", action.file, action.fullprojectFileURL.String())
	}
}

// LinkURLFunc 返回一个函数，用于获取文件的完整 URL
func (action *Upload) LinkURLFunc() func() string {
	return func() string {
		if action.projectFile == nil {
			return ""
		}
		return action.fullprojectFileURL.String()
	}
}

// NewUpload 创建一个新的文件上传操作
func NewUpload(client *gitlab.Client, project string, projectURL *url.URL, file string) *Upload {
	return &Upload{
		client:     client,
		project:    project,
		projectURL: projectURL,
		file:       file,
	}
}
