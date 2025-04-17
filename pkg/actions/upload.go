package actions

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/fanny7d/semrel-gitlab/pkg/workflow"
	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
)

// UploadFileOptions GitLab 上传文件 API 的选项
// 参考: https://docs.gitlab.com/ee/api/projects.html#upload-a-file
type UploadFileOptions struct {
	File     io.Reader `json:"file"`
	FileName string    `json:"filename"`
}

// UploadFileResponse GitLab 上传文件 API 的响应
// 参考: https://docs.gitlab.com/ee/api/projects.html#upload-a-file
type UploadFileResponse struct {
	Alt      string `json:"alt"`
	URL      string `json:"url"`
	Markdown string `json:"markdown"`
}

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

// UploadFile 上传文件
// client: GitLab 客户端
// project: 项目路径
// options: 上传选项
func UploadFile(client *gitlab.Client, project string, options *UploadFileOptions) (*UploadFileResponse, *gitlab.Response, error) {
	u := fmt.Sprintf("projects/%s/uploads", url.QueryEscape(project))

	req, err := client.NewRequest("POST", u, options, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "upload file make request")
	}

	upload := new(UploadFileResponse)
	resp, err := client.Do(req, upload)
	if err != nil {
		return upload, resp, errors.Wrap(err, "upload file execute request")
	}

	return upload, resp, nil
}

// GenerateFileLink 生成文件链接
// project: 项目路径
// filePath: 文件路径
func GenerateFileLink(project string, filePath string) string {
	return fmt.Sprintf("/%s/-/blob/master/%s", project, filePath)
}

// GenerateMarkdownLink 生成 Markdown 格式链接
// text: 链接文本
// url: 链接地址
func GenerateMarkdownLink(text string, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}
