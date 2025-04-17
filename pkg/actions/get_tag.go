package actions

import (
	gitlab "github.com/xanzy/go-gitlab"
	"gitlab.com/juhani/go-semrel-gitlab/pkg/workflow"
)

// GetTag ..
type GetTag struct {
	client  *gitlab.Client
	project string
	tagID   string
	tag     *gitlab.Tag
}

// Do implements Action for GetTag
func (action *GetTag) Do() *workflow.ActionError {
	tag, resp, err := action.client.Tags.GetTag(action.project, action.tagID)
	if err != nil {
		retry := false
		if resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(err, retry)
	}
	action.tag = tag
	return nil
}

// Undo implements Action for GetTag
func (action *GetTag) Undo() error {
	return nil
}

// TagFunc returns accessor func for the tag
func (action *GetTag) TagFunc() func() *gitlab.Tag {
	return func() *gitlab.Tag {
		return action.tag
	}
}

// NewGetTag creates a get tag action
func NewGetTag(client *gitlab.Client, project string, tagID string) *GetTag {
	return &GetTag{
		client:  client,
		project: project,
		tagID:   tagID,
	}
}
