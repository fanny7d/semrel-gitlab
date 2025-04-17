package actions

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"testing"
	"time"

	gitlab "github.com/xanzy/go-gitlab"
)

const (
	apiURL = "http://gitlab/api/v4"
	user   = "root"
	pwd    = "rootpassword"
)

var (
	projectName string
	projectURL  string
	projectPath string
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func newClient(t *testing.T) *gitlab.Client {
	t.Helper()
	client, err := gitlab.NewClient("", gitlab.WithBaseURL(apiURL))
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func setupProject(t *testing.T) (*gitlab.Project, string) {
	t.Helper()
	client := newClient(t)
	projectName = fmt.Sprintf("project_%d", rand.Int())
	projectConf, _, err := client.Projects.CreateProject(
		&gitlab.CreateProjectOptions{
			Name: &projectName,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	projectURL = projectConf.WebURL
	projectPath = projectConf.PathWithNamespace
	return projectConf, projectPath
}

func tagExits(t *testing.T, project *gitlab.Project, tag string) {
	t.Helper()
	client := newClient(t)
	tags, _, err := client.Tags.ListTags(project.ID, &gitlab.ListTagsOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, t := range tags {
		if t.Name == tag {
			return
		}
	}
	t.Fatalf("tag %s not found", tag)
}

func tagExitsWithText(t *testing.T, project *gitlab.Project, tag string, text string) {
	t.Helper()
	client := newClient(t)
	tags, _, err := client.Tags.ListTags(project.ID, &gitlab.ListTagsOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, tagObj := range tags {
		if tagObj.Name == tag {
			if tagObj.Message == text {
				return
			}
			t.Fatalf("tag %s found but message '%s' != '%s'", tag, tagObj.Message, text)
		}
	}
	t.Fatalf("tag %s not found", tag)
}

func teardownProject(t *testing.T, project *gitlab.Project) {
	t.Helper()
	client := newClient(t)
	_, err := client.Projects.DeleteProject(project.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTagAction(t *testing.T) {
	project, projectPath := setupProject(t)
	defer teardownProject(t, project)

	client := newClient(t)
	tag := "v1.0.0"
	action := NewCreateTag(client, projectPath, func() string { return "master" }, tag, "test tag", true)
	err := action.Do()
	if err != nil {
		t.Fatal(err)
	}
	tagExits(t, project, tag)
}

func TestGetTagAction(t *testing.T) {
	project, projectPath := setupProject(t)
	defer teardownProject(t, project)

	client := newClient(t)
	tag := "v1.0.0"
	action := NewCreateTag(client, projectPath, func() string { return "master" }, tag, "test tag", true)
	err := action.Do()
	if err != nil {
		t.Fatal(err)
	}
	tagExits(t, project, tag)

	getAction := NewGetTag(client, projectPath, tag)
	err = getAction.Do()
	if err != nil {
		t.Fatal(err)
	}
	if getAction.tag == nil {
		t.Fatal("tag not found")
	}
	if getAction.tag.Name != tag {
		t.Fatalf("tag name %s != %s", getAction.tag.Name, tag)
	}
}

func TestAddLink(t *testing.T) {
	project, projectPath := setupProject(t)
	defer teardownProject(t, project)

	client := newClient(t)
	tag := "v1.0.0"
	action := NewCreateTag(client, projectPath, func() string { return "master" }, tag, "test tag", true)
	err := action.Do()
	if err != nil {
		t.Fatal(err)
	}
	tagExits(t, project, tag)

	link := "http://example.com"
	linkAction := NewAddLink(&AddLinkParams{
		Client:          client,
		Project:         projectPath,
		LinkDescription: "example",
		MDLinkFunc:      func() string { return fmt.Sprintf("[example](%s)", link) },
		TagFunc:         action.TagFunc(),
	})
	err = linkAction.Do()
	if err != nil {
		t.Fatal(err)
	}
	tagExitsWithText(t, project, tag, "test tag\n\n[example](http://example.com)")
}

func TestUpload(t *testing.T) {
	project, projectPath := setupProject(t)
	defer teardownProject(t, project)

	client := newClient(t)
	file := "test.txt"
	content := "test content"
	err := ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)

	projectURL, err := url.Parse(project.WebURL)
	if err != nil {
		t.Fatal(err)
	}
	action := NewUpload(client, projectPath, projectURL, file)
	err = action.Do()
	if err != nil {
		t.Fatal(err)
	}
	if action.projectFile == nil {
		t.Fatal("file not uploaded")
	}
	if action.fullprojectFileURL == nil {
		t.Fatal("file url not set")
	}
}
