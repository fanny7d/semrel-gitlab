package actions

import (
	"fmt"

	"github.com/fanny7d/semrel-gitlab/pkg/workflow"
	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
)

// CreatePipeline 表示创建 GitLab 流水线的操作
type CreatePipeline struct {
	client     *gitlab.Client
	project    string
	branch     func() string
	pipelineID int
}

// Do 实现 Action 接口，执行创建流水线的操作
func (action *CreatePipeline) Do() *workflow.ActionError {
	if action.pipelineID != 0 {
		return nil
	}
	options := &gitlab.CreatePipelineOptions{
		Ref: gitlab.String(action.branch()),
	}
	pipeline, _, err := action.client.Pipelines.CreatePipeline(action.project, options)
	if err != nil {
		return workflow.NewActionError(errors.Wrap(err, "create pipeline"), true)
	}
	action.pipelineID = pipeline.ID
	return nil
}

// Undo 实现 Action 接口，撤销创建流水线的操作
func (action *CreatePipeline) Undo() error {
	if action.pipelineID == 0 {
		return nil
	}
	fmt.Printf("Pipeline cannot be rolled back: %d.\n", action.pipelineID)
	return nil
}

// NewCreatePipeline 创建一个新的创建流水线操作
func NewCreatePipeline(client *gitlab.Client, project string, branch func() string) *CreatePipeline {
	return &CreatePipeline{
		client:  client,
		project: project,
		branch:  branch,
	}
}
