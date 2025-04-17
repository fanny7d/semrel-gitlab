package gitlabutil

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

// UpdateReleaseOptions GitLab 发布更新 API 的选项
// 参考: https://docs.gitlab.com/ee/api/tags.html#update-a-release
type UpdateReleaseOptions struct {
	ID          string `json:"id"`
	TagName     string `json:"tag_name"`
	Description string `json:"description"`
}

// UpdateReleaseResponse GitLab 发布更新 API 的响应
// 参考: https://docs.gitlab.com/ee/api/tags.html#update-a-release
type UpdateReleaseResponse struct {
	TagName     string `json:"tag_name"`
	Description string `json:"description"`
}

// UpdateTagDescription 更新标签的描述信息
// client: GitLab 客户端
// project: 项目路径
// tagID: 标签 ID
// description: 新的描述信息
func UpdateTagDescription(client *gitlab.Client, project string, tagID string, description string) (*UpdateReleaseResponse, *gitlab.Response, error) {
	updateOptions := &UpdateReleaseOptions{
		ID:          project,
		TagName:     tagID,
		Description: description,
	}
	updateResp := UpdateReleaseResponse{}
	u := fmt.Sprintf("projects/%s/repository/tags/%s/release", url.QueryEscape(project), tagID)
	req, err := client.NewRequest("PUT", u, updateOptions, []gitlab.RequestOptionFunc{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "add-download make request")
	}

	resp, err := client.Do(req, &updateResp)
	if err != nil {
		return &updateResp, resp, errors.Wrap(err, "add-download execute request")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &updateResp, resp, errors.Errorf("unexpected status for release update %d", resp.StatusCode)
	}

	return &updateResp, resp, nil
}
