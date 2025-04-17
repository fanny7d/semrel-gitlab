package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fanny7d/semrel-gitlab/pkg/domain"
	"github.com/fanny7d/semrel-gitlab/pkg/service"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "生成变更日志",
	Long: `分析提交信息并生成变更日志。

变更日志将包含以下内容：
- 版本号
- 发布日期
- 变更类型（新功能、修复、重构等）
- 提交信息
- 提交者信息

变更日志将写入到 CHANGELOG.md 文件中。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取全局选项
		patchTypes := strings.Split(cmd.Flag("patch-commit-types").Value.String(), ",")
		minorTypes := strings.Split(cmd.Flag("minor-commit-types").Value.String(), ",")
		tagPrefix := cmd.Flag("tag-prefix").Value.String()

		// 创建 Git 服务
		gitService := service.NewGitService(patchTypes, minorTypes, tagPrefix)

		// 分析提交
		release, err := gitService.AnalyzeCommits()
		if err != nil {
			return err
		}

		// 生成变更日志
		changelog := generateChangelog(release)

		// 写入文件
		if err := os.WriteFile("CHANGELOG.md", []byte(changelog), 0644); err != nil {
			return fmt.Errorf("写入变更日志失败: %v", err)
		}

		return nil
	},
}

func generateChangelog(release *domain.Release) string {
	var changelog strings.Builder

	// 添加标题
	changelog.WriteString("# 变更日志\n\n")

	// 添加当前版本
	changelog.WriteString(fmt.Sprintf("## [%s] - %s\n\n", release.Version.Next.String(), time.Now().Format("2006-01-02")))

	// 添加变更类型
	for category, changes := range release.Changes {
		if len(changes) == 0 {
			continue
		}

		// 转换类别名称
		categoryName := getCategoryName(category)

		// 添加类别标题
		changelog.WriteString(fmt.Sprintf("### %s\n\n", categoryName))

		// 添加变更列表
		for _, change := range changes {
			changelog.WriteString(fmt.Sprintf("- %s\n", change.Subject))
		}

		changelog.WriteString("\n")
	}

	return changelog.String()
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

func init() {
	rootCmd.AddCommand(changelogCmd)
}
