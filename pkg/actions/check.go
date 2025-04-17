package actions

import (
	"github.com/fanny7d/semrel-gitlab/pkg/workflow"
	gitlab "github.com/xanzy/go-gitlab"
)

// Check action tests api connection
type Check struct {
	client   *gitlab.Client
	Version  string
	Revision string
}

// Do implements Action for Check
func (action *Check) Do() *workflow.ActionError {

	_, resp, err := action.client.Users.CurrentUser()
	if err != nil {
		retry := false
		if resp != nil && resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(err, retry)
	}
	v, resp, err := action.client.Version.GetVersion()
	if err != nil {
		retry := false
		if resp != nil && resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(err, retry)
	}
	action.Version = v.Version
	action.Revision = v.Revision
	return nil
}

// Undo implements Action for Check
func (action *Check) Undo() error {
	return nil
}

// NewCheck creates a Check action
func NewCheck(client *gitlab.Client) *Check {
	return &Check{
		client: client,
	}
}
