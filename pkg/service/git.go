package service

import (
	"fmt"
	"time"

	"github.com/fanny7d/semrel-gitlab/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// GitService 提供 Git 相关操作
type GitService struct {
	patchTypes []string
	minorTypes []string
	tagPrefix  string
}

// NewGitService 创建一个新的 Git 服务
func NewGitService(patchTypes, minorTypes []string, tagPrefix string) *GitService {
	return &GitService{
		patchTypes: patchTypes,
		minorTypes: minorTypes,
		tagPrefix:  tagPrefix,
	}
}

// AnalyzeCommits 分析提交历史并返回发布数据
func (s *GitService) AnalyzeCommits() (*domain.Release, error) {
	// 打开 Git 仓库
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, errors.Wrap(err, "打开 Git 仓库失败")
	}

	// 获取 HEAD 引用
	head, err := repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "获取 HEAD 引用失败")
	}

	// 获取提交历史
	iter, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		return nil, errors.Wrap(err, "获取提交历史失败")
	}

	// 创建版本对象
	version := domain.NewVersion(time.Now())
	release := domain.NewRelease(version, s.tagPrefix)

	// 分析每个提交
	err = iter.ForEach(func(commit *object.Commit) error {
		// 解析提交消息
		msg := commit.Message
		if msg == "" {
			return nil
		}

		// 创建提交对象
		c := domain.NewCommit(
			commit.Hash.String(),
			domain.CommitType("fix"), // 默认为修复类型
			"",                       // 暂不支持 scope
			msg,
			"",
			false,
		)

		// 确定版本升级级别
		level := c.DetermineLevel(s.patchTypes, s.minorTypes)
		version.Bump(level)

		// 添加到变更列表
		category := string(c.Type)
		if c.Breaking {
			category = "breaking"
		}
		release.AddChange(category, c)

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "分析提交历史失败")
	}

	return release, nil
}

// CreateTag 创建 Git 标签
func (s *GitService) CreateTag(tagName string) error {
	// 打开 Git 仓库
	repo, err := git.PlainOpen(".")
	if err != nil {
		return errors.Wrap(err, "打开 Git 仓库失败")
	}

	// 获取 HEAD 引用
	head, err := repo.Head()
	if err != nil {
		return errors.Wrap(err, "获取 HEAD 引用失败")
	}

	// 创建标签
	_, err = repo.CreateTag(tagName, head.Hash(), &git.CreateTagOptions{
		Message: fmt.Sprintf("Release %s", tagName),
	})
	if err != nil {
		return errors.Wrap(err, "创建标签失败")
	}

	return nil
}
