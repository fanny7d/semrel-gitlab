package actions

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
	"gitlab.com/juhani/go-semrel-gitlab/pkg/gitlabutil"
	"gitlab.com/juhani/go-semrel-gitlab/pkg/workflow"
)

// AddLinkParams ..
type AddLinkParams struct {
	Client          *gitlab.Client
	Project         string
	LinkDescription string
	LinkURLFunc     func() string
	// MDLinkFunc returns markdown link of form '[link text](https://foo.bar....)'
	MDLinkFunc           func() string
	TagFunc              func() *gitlab.Tag
	ReleasesAPIAvailable bool
}

// AddLink action
type AddLink struct {
	*AddLinkParams
	origRelNote string
	linkID      *int
}

// Do implements Action for AddLink
func (action *AddLink) Do() *workflow.ActionError {
	if action.ReleasesAPIAvailable {
		return action.doWithReleases()
	}
	return action.doWithoutReleases()
}

func (action *AddLink) doWithReleases() *workflow.ActionError {
	if action.linkID != nil {
		return nil
	}

	tag := action.TagFunc()

	options := gitlab.CreateReleaseLinkOptions{
		Name: gitlab.String(action.LinkDescription),
		URL:  gitlab.String(action.LinkURLFunc()),
	}
	link, resp, err := action.Client.ReleaseLinks.CreateReleaseLink(action.Project, tag.Name, &options)
	if err != nil {
		retry := false
		if resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(errors.Wrap(err, "create link"), retry)
	}
	action.linkID = &link.ID
	return nil
}

func (action *AddLink) doWithoutReleases() *workflow.ActionError {
	if len(action.origRelNote) > 0 {
		return nil
	}
	markdownLink := action.MDLinkFunc()
	tag := action.TagFunc()

	newRelNote := addDownloadLinkToDescription(tag.Release.Description, markdownLink, action.LinkDescription)

	_, resp, err := gitlabutil.UpdateTagDescription(action.Client, action.Project, tag.Name, newRelNote)
	if err != nil {
		retry := false
		if resp.StatusCode == 502 {
			retry = true
		}
		return workflow.NewActionError(errors.Wrap(err, "update description"), retry)
	}
	action.origRelNote = tag.Release.Description
	return nil
}

// Undo implements Action for AddLink
func (action *AddLink) Undo() error {
	if action.ReleasesAPIAvailable {
		return action.undoWithReleases()
	}
	return action.undoWithoutReleases()
}

func (action *AddLink) undoWithReleases() error {
	if action.linkID == nil {
		return nil
	}
	tag := action.TagFunc()

	_, _, err := action.Client.ReleaseLinks.DeleteReleaseLink(action.Project, tag.Name, *action.linkID)
	if err != nil {
		fmt.Printf(`
MANUAL ACTION REQUIRED!!
Removing link '%s' from release %s failed.
`, action.LinkDescription, tag.Name)
		return errors.Wrap(err, "remove link")
	}
	action.linkID = nil
	return nil
}

func (action *AddLink) undoWithoutReleases() error {
	if len(action.origRelNote) == 0 {
		return nil
	}
	tag := action.TagFunc()

	_, _, err := gitlabutil.UpdateTagDescription(action.Client, action.Project, tag.Name, action.origRelNote)
	if err != nil {
		fmt.Printf(`
MANUAL ACTION REQUIRED!!
Restoring the description for tag %s failed.
The text of the original description is:
-----
%s
-----\n`, tag.Name, action.origRelNote)
		return errors.Wrap(err, "update description")
	}
	action.origRelNote = ""
	return nil
}

func addDownloadLinkToDescription(origReleaseText string, markdownLink string, description string) string {
	parts := []string{}
	marker := "\n<!--- download here -->"

	if strings.Index(origReleaseText, "<!--- downloads here -->") > -1 {
		parts = append(parts, "### Downloads", "")
		marker = "<!--- downloads here -->"
	} else if strings.Index(origReleaseText, "<!--- download here -->") == -1 {
		// marker was not found, add it to the end of release note
		origReleaseText += marker
	}

	parts = append(parts, fmt.Sprintf("- **%s:** %s", markdownLink, description), "", "<!--- download here -->")
	return strings.Replace(origReleaseText, marker, strings.Join(parts, "\n"), 1)
}

// NewAddLink ..
func NewAddLink(params *AddLinkParams) *AddLink {
	return &AddLink{
		AddLinkParams: params,
	}
}
