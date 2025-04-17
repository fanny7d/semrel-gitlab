package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fanny7d/semrel-gitlab/pkg/domain"
	"github.com/pkg/errors"
)

// RenderService 提供渲染相关操作
type RenderService struct {
	changelogFile string
}

// NewRenderService 创建一个新的渲染服务
func NewRenderService(changelogFile string) *RenderService {
	return &RenderService{
		changelogFile: changelogFile,
	}
}

// RenderReleaseNote 渲染发布说明
func (s *RenderService) RenderReleaseNote(release *domain.Release) error {
	var buf bytes.Buffer

	// 渲染标题
	buf.WriteString(fmt.Sprintf("# %s\n\n", release.TagName))

	// 渲染变更列表
	for category, changes := range release.Changes {
		if category == "other" {
			continue
		}

		buf.WriteString(fmt.Sprintf("## %s\n\n", strings.Title(category)))
		for _, change := range changes {
			if change.IsPreReleased() {
				continue
			}
			if change.Scope != "" {
				buf.WriteString(fmt.Sprintf("* **%s:** %s (%s)\n", change.Scope, change.Subject, change.Hash[:7]))
			} else {
				buf.WriteString(fmt.Sprintf("* %s (%s)\n", change.Subject, change.Hash[:7]))
			}
		}
		buf.WriteString("\n")
	}

	// 渲染下载链接
	if len(release.Links) > 0 {
		buf.WriteString("## 下载\n\n")
		for _, link := range release.Links {
			if link.Description != "" {
				buf.WriteString(fmt.Sprintf("* [%s](%s) - %s\n", link.Name, link.URL, link.Description))
			} else {
				buf.WriteString(fmt.Sprintf("* [%s](%s)\n", link.Name, link.URL))
			}
		}
		buf.WriteString("\n")
	}

	release.Message = buf.String()
	return nil
}

// UpdateChangelog 更新更新日志文件
func (s *RenderService) UpdateChangelog(release *domain.Release) error {
	entry := s.renderChangelogEntry(release)

	// 读取现有的更新日志文件
	content, err := ioutil.ReadFile(s.changelogFile)
	if err != nil {
		if !errors.Is(err, errors.New("file does not exist")) {
			return errors.Wrap(err, "读取更新日志文件失败")
		}
		// 创建新的更新日志文件
		data := strings.Join([]string{
			"# CHANGELOG",
			"<!--- next entry here -->",
			entry,
		}, "\n\n")
		return ioutil.WriteFile(s.changelogFile, []byte(data), 0644)
	}

	// 在标记处插入新条目
	parts := strings.Split(string(content), "<!--- next entry here -->")
	if len(parts) != 2 {
		return errors.New("更新日志文件格式错误")
	}

	data := strings.Join([]string{
		strings.TrimRight(parts[0], " \n\r\t"),
		"<!--- next entry here -->",
		entry,
		strings.TrimLeft(parts[1], " \n\r\t"),
	}, "\n\n")

	return ioutil.WriteFile(s.changelogFile, []byte(data), 0644)
}

// renderChangelogEntry 渲染更新日志条目
func (s *RenderService) renderChangelogEntry(release *domain.Release) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("## %s\n\n", release.TagName))

	for category, changes := range release.Changes {
		if category == "other" {
			continue
		}

		buf.WriteString(fmt.Sprintf("### %s\n\n", strings.Title(category)))
		for _, change := range changes {
			if change.IsPreReleased() {
				continue
			}
			if change.Scope != "" {
				buf.WriteString(fmt.Sprintf("* **%s:** %s\n", change.Scope, change.Subject))
			} else {
				buf.WriteString(fmt.Sprintf("* %s\n", change.Subject))
			}
		}
		buf.WriteString("\n")
	}

	return buf.String()
}
