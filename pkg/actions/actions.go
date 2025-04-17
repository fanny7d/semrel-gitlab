package actions

import (
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
