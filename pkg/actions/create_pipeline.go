package actions

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
)

// CreatePipelineOptions GitLab 创建流水线 API 的选项
// 参考: https://docs.gitlab.com/ee/api/pipelines.html#create-a-new-pipeline
type CreatePipelineOptions struct {
	Ref       string            `json:"ref"`
	Variables map[string]string `json:"variables,omitempty"`
}

// CreatePipelineResponse GitLab 创建流水线 API 的响应
// 参考: https://docs.gitlab.com/ee/api/pipelines.html#create-a-new-pipeline
type CreatePipelineResponse struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Ref       string `json:"ref"`
	Sha       string `json:"sha"`
	WebURL    string `json:"web_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreatePipelineAPI 创建流水线
// client: GitLab 客户端
// project: 项目路径
// options: 流水线选项
func CreatePipelineAPI(client *gitlab.Client, project string, options *CreatePipelineOptions) (*CreatePipelineResponse, *gitlab.Response, error) {
	u := fmt.Sprintf("projects/%s/pipeline", url.QueryEscape(project))

	req, err := client.NewRequest("POST", u, options, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create pipeline make request")
	}

	pipeline := new(CreatePipelineResponse)
	resp, err := client.Do(req, pipeline)
	if err != nil {
		return pipeline, resp, errors.Wrap(err, "create pipeline execute request")
	}

	return pipeline, resp, nil
}
