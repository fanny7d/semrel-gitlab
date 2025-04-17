package actions

import (
	"fmt"

	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
	"gitlab.com/juhani/go-semrel-gitlab/pkg/workflow"
)

// CreateTag ..
type CreateTag struct {
	client               *gitlab.Client
	project              string
	refFunc              func() string
	tagID                string
	note                 string
	tag                  *gitlab.Tag
	release              *gitlab.Release
	releasesAPIAvailable bool
}

// Do implements Action for CreateTag
func (action *CreateTag) Do() *workflow.ActionError {
	if action.releasesAPIAvailable {
		return action.doWithReleases()
	}
	return action.doWithoutReleases()
}

func (action *CreateTag) doWithReleases() *workflow.ActionError {
	if action.release != nil {
		return nil
	}
	ref := action.refFunc()
	message := fmt.Sprintf("Release %s", action.tagID)
	options := &gitlab.CreateReleaseOptions{
		Description: gitlab.String(action.note),
		Name:        gitlab.String(message),
		TagName:     gitlab.String(action.tagID),
		Ref:         gitlab.String(ref),
	}
	release, resp, err := action.client.Releases.CreateRelease(action.project, options)
	if err != nil {
		retry := false
		if resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(errors.Wrap(err, "release create"), retry)
	}
	action.release = release
	return nil
}

func (action *CreateTag) doWithoutReleases() *workflow.ActionError {
	if action.tag != nil {
		return nil
	}
	ref := action.refFunc()
	message := fmt.Sprintf("Release %s", action.tagID)
	options := &gitlab.CreateTagOptions{
		TagName: &action.tagID,
		Ref:     &ref,
		Message: &message,
	}
	if len(action.note) > 0 {
		options.ReleaseDescription = &action.note
	}
	// tag might already exist if 502 occurred on earlier attempt
	tag, resp, err := action.client.Tags.GetTag(action.project, action.tagID)
	// tag already exists
	if err == nil {
		action.tag = tag
		return nil
	}
	// something else than NOT FOUND was returned
	if resp.StatusCode != 404 {
		retry := false
		if resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(err, retry)
	}
	tag, resp, err = action.client.Tags.CreateTag(action.project, options)
	if err != nil {
		retry := false
		if resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(errors.Wrap(err, "release create tag"), retry)
	}
	action.tag = tag
	return nil
}

// Undo implements Action for CreateTag
func (action *CreateTag) Undo() error {
	if action.releasesAPIAvailable {
		return action.undoWithReleases()
	}
	return action.undoWithoutReleases()
}

func (action *CreateTag) undoWithReleases() error {
	if action.release == nil {
		return nil
	}
	_, _, err := action.client.Releases.DeleteRelease(action.project, action.tagID)
	if err != nil {
		fmt.Printf(`
MANUAL ACTION REQUIRED!!
Removing release '%s' failed.\n`, action.release.Name)
		return err
	}
	return nil
}

func (action *CreateTag) undoWithoutReleases() error {
	if action.tag == nil {
		return nil
	}
	_, err := action.client.Tags.DeleteTag(action.project, action.tagID)
	if err != nil {
		fmt.Printf(`
MANUAL ACTION REQUIRED!!
Removing tag '%s' failed.\n`, action.tag.Name)
		return err
	}
	return nil
}

// TagFunc returns accessor func for the tag
func (action *CreateTag) TagFunc() func() *gitlab.Tag {
	return func() *gitlab.Tag {
		return action.tag
	}
}

// NewCreateTag ..
func NewCreateTag(client *gitlab.Client, project string, refFunc func() string, tagID string, note string, releasesAPIAvailable bool) *CreateTag {
	return &CreateTag{
		client:               client,
		project:              project,
		refFunc:              refFunc,
		tagID:                tagID,
		note:                 note,
		releasesAPIAvailable: releasesAPIAvailable,
	}
}

// FuncOfString wraps given string in a function
func FuncOfString(ref string) func() string {
	return func() string {
		return ref
	}
}
