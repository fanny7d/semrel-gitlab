package service

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fanny7d/semrel-gitlab/pkg/domain"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

// httpClientWithTimeout 创建一个带有超时设置的 HTTP 客户端
func httpClientWithTimeout(skipSSLVerify bool) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSLVerify,
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// GitLabService 提供 GitLab 相关操作
type GitLabService struct {
	client     *gitlab.Client
	project    string
	projectURL *url.URL
}

// GitLabClient 提供 GitLab API 操作
type GitLabClient struct {
	client *gitlab.Client
}

// NewGitLabService 创建一个新的 GitLab 服务
func NewGitLabService(token, apiURL, project string, projectURL *url.URL, skipSSLVerify bool) (*GitLabService, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL))
	if err != nil {
		return nil, errors.Wrap(err, "创建 GitLab 客户端失败")
	}

	return &GitLabService{
		client:     client,
		project:    project,
		projectURL: projectURL,
	}, nil
}

// NewGitLabClient 创建一个新的 GitLab 客户端
func NewGitLabClient(token, apiURL string, skipSSLVerify bool) (*GitLabClient, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL), gitlab.WithHTTPClient(httpClientWithTimeout(skipSSLVerify)))
	if err != nil {
		return nil, fmt.Errorf("创建 GitLab 客户端失败: %v", err)
	}

	return &GitLabClient{client: client}, nil
}

// GetTag 获取标签信息
func (s *GitLabService) GetTag(tagName string) (*domain.Release, error) {
	tag, _, err := s.client.Tags.GetTag(s.project, tagName)
	if err != nil {
		return nil, errors.Wrap(err, "获取标签失败")
	}

	// 创建发布对象
	version := domain.NewVersion(*tag.Commit.CreatedAt)
	release := domain.NewRelease(version, "")
	release.TagName = tagName
	release.Message = tag.Message

	return release, nil
}

// CreateTag 创建标签和发布说明
func (s *GitLabService) CreateTag(release *domain.Release, ref string) error {
	// 创建标签
	_, _, err := s.client.Tags.CreateTag(s.project, &gitlab.CreateTagOptions{
		TagName: gitlab.String(release.TagName),
		Ref:     gitlab.String(ref),
		Message: gitlab.String(release.Message),
	})
	if err != nil {
		return errors.Wrap(err, "创建标签失败")
	}

	return nil
}

// CreateCommit 创建提交
func (s *GitLabService) CreateCommit(branch, message string, files []string) error {
	// 创建提交
	actions := make([]*gitlab.CommitActionOptions, 0, len(files))
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return errors.Wrap(err, "读取文件失败")
		}

		action := &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(gitlab.FileUpdate),
			FilePath: gitlab.String(file),
			Content:  gitlab.String(string(content)),
		}
		actions = append(actions, action)
	}

	_, _, err := s.client.Commits.CreateCommit(s.project, &gitlab.CreateCommitOptions{
		Branch:        gitlab.String(branch),
		CommitMessage: gitlab.String(message),
		Actions:       actions,
	})
	if err != nil {
		return errors.Wrap(err, "创建提交失败")
	}

	return nil
}

// UploadFile 上传文件并添加到发布说明
func (s *GitLabService) UploadFile(release *domain.Release, filePath string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "打开文件失败")
	}
	defer file.Close()

	// 上传文件
	upload, _, err := s.client.Projects.UploadFile(s.project, file, filePath)
	if err != nil {
		return errors.Wrap(err, "上传文件失败")
	}

	// 添加下载链接
	release.AddLink(
		upload.Alt,
		s.projectURL.String()+upload.URL,
		"",
	)

	return nil
}

// CreatePipeline 创建管道
func (s *GitLabService) CreatePipeline(ref string) error {
	_, _, err := s.client.Pipelines.CreatePipeline(s.project, &gitlab.CreatePipelineOptions{
		Ref: gitlab.String(ref),
	})
	if err != nil {
		return errors.Wrap(err, "创建管道失败")
	}

	return nil
}

// CreateRelease 在 GitLab 上创建发布
func (c *GitLabClient) CreateRelease(projectPath, tagName string, release *domain.Release) error {
	// 创建发布说明
	description := generateReleaseDescription(release)

	// 创建发布
	_, _, err := c.client.Releases.CreateRelease(projectPath, &gitlab.CreateReleaseOptions{
		Name:        gitlab.String(release.Version.Next.String()),
		TagName:     gitlab.String(tagName),
		Description: gitlab.String(description),
	})
	if err != nil {
		return fmt.Errorf("创建发布失败: %v", err)
	}

	// 添加下载链接
	for _, link := range release.Links {
		_, _, err := c.client.ReleaseLinks.CreateReleaseLink(projectPath, tagName, &gitlab.CreateReleaseLinkOptions{
			Name: gitlab.String(link.Name),
			URL:  gitlab.String(link.URL),
		})
		if err != nil {
			return fmt.Errorf("添加下载链接失败: %v", err)
		}
	}

	return nil
}

func generateReleaseDescription(release *domain.Release) string {
	var description strings.Builder

	// 添加版本信息
	description.WriteString(fmt.Sprintf("# %s\n\n", release.Version.Next.String()))

	// 添加变更类型
	for category, changes := range release.Changes {
		if len(changes) == 0 {
			continue
		}

		// 转换类别名称
		categoryName := getCategoryName(category)

		// 添加类别标题
		description.WriteString(fmt.Sprintf("## %s\n\n", categoryName))

		// 添加变更列表
		for _, change := range changes {
			description.WriteString(fmt.Sprintf("- %s\n", change.Subject))
		}

		description.WriteString("\n")
	}

	return description.String()
}

func getCategoryName(category string) string {
	switch category {
	case "feat":
		return "新功能"
	case "fix":
		return "修复"
	case "refactor":
		return "重构"
	case "perf":
		return "性能优化"
	case "docs":
		return "文档更新"
	case "style":
		return "代码格式调整"
	case "test":
		return "测试相关"
	case "chore":
		return "构建过程或辅助工具的变动"
	case "breaking":
		return "破坏性变更"
	default:
		return category
	}
}
