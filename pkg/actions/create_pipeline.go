package actions

import (
	"fmt"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
	"gitlab.com/juhani/go-semrel-gitlab/pkg/workflow"
)

// CreatePipeline ..
type CreatePipeline struct {
	client   *gitlab.Client
	project  string
	refFunc  func() string
	pipeline *gitlab.Pipeline
}

// Do implements Action for CreatePipeline
func (action *CreatePipeline) Do() *workflow.ActionError {
	if action.pipeline != nil {
		return nil
	}

	ref := action.refFunc()
	pipeline, _, err := action.client.Pipelines.CreatePipeline(action.project, &gitlab.CreatePipelineOptions{
		Ref: &ref,
	})
	if err != nil {
		// retry if reference was not found
		if strings.Contains(err.Error(), "Reference not found") {
			return workflow.NewActionError(err, true)
		}
		// continue normally if pipeline creation fails because there's nothing to do
		if !strings.Contains(err.Error(), "No stages / jobs for this pipeline.") {
			return workflow.NewActionError(err, false)
		}
	}
	action.pipeline = pipeline

	return nil
}

// Undo implements Action for CreatePipeline
func (action *CreatePipeline) Undo() error {
	if action.pipeline == nil {
		return nil
	}
	fmt.Printf(`
MANUAL ACTION REQUIRED!
Attempting to cancel pipeline %d.
Jobs may have been executed, already.`, action.pipeline.ID)
	_, _, err := action.client.Pipelines.CancelPipelineBuild(action.project, action.pipeline.ID)
	if err != nil {
		return err
	}
	action.pipeline = nil
	return nil
}

// NewCreatePipeline ..
func NewCreatePipeline(client *gitlab.Client, project string, refFunc func() string) *CreatePipeline {
	return &CreatePipeline{
		client:  client,
		project: project,
		refFunc: refFunc,
	}
}
