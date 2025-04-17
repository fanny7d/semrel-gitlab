package actions

import (
	"fmt"
	"net/url"

	"github.com/fanny7d/semrel-gitlab/pkg/gitlabutil"
	"github.com/fanny7d/semrel-gitlab/pkg/workflow"
	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
)

// FuncOfString 是一个返回字符串的函数类型
type FuncOfString func() string

// NewFuncOfString 创建一个返回固定字符串的函数
func NewFuncOfString(s string) func() string {
	return func() string {
		return s
	}
}

// CreateTagOptions GitLab 创建标签 API 的选项
// 参考: https://docs.gitlab.com/ee/api/tags.html#create-a-new-tag
type CreateTagOptions struct {
	TagName     string `json:"tag_name"`
	Ref         string `json:"ref"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateTagResponse GitLab 创建标签 API 的响应
// 参考: https://docs.gitlab.com/ee/api/tags.html#create-a-new-tag
type CreateTagResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Target  string `json:"target"`
	Commit  struct {
		ID             string   `json:"id"`
		ShortID        string   `json:"short_id"`
		Title          string   `json:"title"`
		CreatedAt      string   `json:"created_at"`
		ParentIDs      []string `json:"parent_ids"`
		Message        string   `json:"message"`
		AuthorName     string   `json:"author_name"`
		AuthorEmail    string   `json:"author_email"`
		AuthoredDate   string   `json:"authored_date"`
		CommitterName  string   `json:"committer_name"`
		CommitterEmail string   `json:"committer_email"`
		CommittedDate  string   `json:"committed_date"`
	} `json:"commit"`
	Release struct {
		TagName     string `json:"tag_name"`
		Description string `json:"description"`
	} `json:"release"`
	Protected bool `json:"protected"`
}

// CreateTag 表示创建 GitLab 标签的操作
type CreateTag struct {
	client     *gitlab.Client
	project    string
	branch     func() string
	tag        string
	message    string
	force      bool
	createdTag *gitlab.Tag
}

// Do 实现 Action 接口，执行创建标签的操作
func (action *CreateTag) Do() *workflow.ActionError {
	if action.createdTag != nil {
		return nil
	}
	options := &gitlab.CreateTagOptions{
		TagName: &action.tag,
		Ref:     gitlab.String(action.branch()),
		Message: &action.message,
	}
	tag, _, err := action.client.Tags.CreateTag(action.project, options)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "create tag"), true)
	}
	action.createdTag = tag
	return nil
}

// Undo 实现 Action 接口，撤销创建标签的操作
func (action *CreateTag) Undo() error {
	if action.createdTag == nil {
		return nil
	}
	_, err := action.client.Tags.DeleteTag(action.project, action.tag)
	return err
}

// TagFunc 返回一个函数，用于获取创建的标签名称
func (action *CreateTag) TagFunc() func() string {
	return func() string {
		if action.createdTag == nil {
			return ""
		}
		return action.createdTag.Name
	}
}

// NewCreateTag 创建一个新的创建标签操作
func NewCreateTag(client *gitlab.Client, project string, branch func() string, tag string, message string, force bool) *CreateTag {
	return &CreateTag{
		client:  client,
		project: project,
		branch:  branch,
		tag:     tag,
		message: message,
		force:   force,
	}
}

// GetTagOptions GitLab 获取标签 API 的选项
// 参考: https://docs.gitlab.com/ee/api/tags.html#get-a-single-tag
type GetTagOptions struct {
	TagName string `json:"tag_name"`
}

// GetTagResponse GitLab 获取标签 API 的响应
// 参考: https://docs.gitlab.com/ee/api/tags.html#get-a-single-tag
type GetTagResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Target  string `json:"target"`
	Commit  struct {
		ID             string   `json:"id"`
		ShortID        string   `json:"short_id"`
		Title          string   `json:"title"`
		CreatedAt      string   `json:"created_at"`
		ParentIDs      []string `json:"parent_ids"`
		Message        string   `json:"message"`
		AuthorName     string   `json:"author_name"`
		AuthorEmail    string   `json:"author_email"`
		AuthoredDate   string   `json:"authored_date"`
		CommitterName  string   `json:"committer_name"`
		CommitterEmail string   `json:"committer_email"`
		CommittedDate  string   `json:"committed_date"`
	} `json:"commit"`
	Release struct {
		TagName     string `json:"tag_name"`
		Description string `json:"description"`
	} `json:"release"`
	Protected bool `json:"protected"`
}

// GetTag 表示获取 GitLab 标签的操作
type GetTag struct {
	client  *gitlab.Client
	project string
	tag     string
	tagObj  *gitlab.Tag
}

// Do 实现 Action 接口，执行获取标签的操作
func (action *GetTag) Do() *workflow.ActionError {
	if action.tagObj != nil {
		return nil
	}
	tag, _, err := action.client.Tags.GetTag(action.project, action.tag)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "get tag"), true)
	}
	action.tagObj = tag
	return nil
}

// Undo 实现 Action 接口，撤销获取标签的操作（无需实现）
func (action *GetTag) Undo() error {
	return nil
}

// TagFunc 返回一个函数，用于获取标签对象
func (action *GetTag) TagFunc() func() string {
	return func() string {
		if action.tagObj == nil {
			return ""
		}
		return action.tagObj.Name
	}
}

// NewGetTag 创建一个新的获取标签操作
func NewGetTag(client *gitlab.Client, project string, tag string) *GetTag {
	return &GetTag{
		client:  client,
		project: project,
		tag:     tag,
	}
}

// AddLinkParams 表示添加链接操作的参数
type AddLinkParams struct {
	Client               *gitlab.Client
	Project              string
	LinkDescription      string
	MDLinkFunc           func() string
	TagFunc              func() string
	LinkURLFunc          func() string
	ReleasesAPIAvailable bool
}

// AddLink 表示添加 GitLab 标签链接的操作
type AddLink struct {
	client               *gitlab.Client
	project              string
	linkDescription      string
	mdLinkFunc           func() string
	tagFunc              func() string
	linkURLFunc          func() string
	releasesAPIAvailable bool
}

// Do 实现 Action 接口，执行添加链接的操作
func (action *AddLink) Do() *workflow.ActionError {
	tag := action.tagFunc()
	if tag == "" {
		return workflow.NewActionError(errors.New("tag not set"), false)
	}
	link := action.mdLinkFunc()
	if link == "" {
		return workflow.NewActionError(errors.New("link not set"), false)
	}
	_, _, err := gitlabutil.UpdateTagDescription(action.client, action.project, tag, action.linkDescription+"\n\n"+link)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "add link"), true)
	}
	return nil
}

// Undo 实现 Action 接口，撤销添加链接的操作
func (action *AddLink) Undo() error {
	tag := action.tagFunc()
	if tag == "" {
		return nil
	}
	_, _, err := gitlabutil.UpdateTagDescription(action.client, action.project, tag, action.linkDescription)
	return err
}

// NewAddLink 创建一个新的添加链接操作
func NewAddLink(params *AddLinkParams) *AddLink {
	return &AddLink{
		client:               params.Client,
		project:              params.Project,
		linkDescription:      params.LinkDescription,
		mdLinkFunc:           params.MDLinkFunc,
		tagFunc:              params.TagFunc,
		linkURLFunc:          params.LinkURLFunc,
		releasesAPIAvailable: params.ReleasesAPIAvailable,
	}
}

// AddReleaseLinkOptions GitLab 添加发布链接 API 的选项
// 参考: https://docs.gitlab.com/ee/api/releases/links.html#create-a-link
type AddReleaseLinkOptions struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	FilePath string `json:"filepath,omitempty"`
	LinkType string `json:"link_type,omitempty"`
}

// AddReleaseLinkResponse GitLab 添加发布链接 API 的响应
// 参考: https://docs.gitlab.com/ee/api/releases/links.html#create-a-link
type AddReleaseLinkResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	DirectURL string `json:"direct_asset_url"`
	External  bool   `json:"external"`
	LinkType  string `json:"link_type"`
}

// AddReleaseLink 添加发布链接
// client: GitLab 客户端
// project: 项目路径
// tagName: 标签名称
// options: 链接选项
func AddReleaseLink(client *gitlab.Client, project string, tagName string, options *AddReleaseLinkOptions) (*AddReleaseLinkResponse, *gitlab.Response, error) {
	u := fmt.Sprintf("projects/%s/releases/%s/assets/links", url.QueryEscape(project), tagName)

	req, err := client.NewRequest("POST", u, options, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "add release link make request")
	}

	link := new(AddReleaseLinkResponse)
	resp, err := client.Do(req, link)
	if err != nil {
		return link, resp, errors.Wrap(err, "add release link execute request")
	}

	return link, resp, nil
}
