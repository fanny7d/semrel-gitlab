package cmd

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/fanny7d/semrel-gitlab/pkg/domain"
	"github.com/fanny7d/semrel-gitlab/pkg/service"
	"github.com/spf13/cobra"
)

var nextVersionCmd = &cobra.Command{
	Use:   "next-version",
	Short: "分析提交消息并打印下一个版本号",
	Long: `分析提交消息并打印下一个版本号。

此命令会遍历 HEAD 的父提交，收集未发布的更改，并与每个分支中遇到的第一个标签进行比较，
以确定下一个版本的基准。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取命令行参数
		bumpPatch, _ := cmd.Flags().GetBool("bump-patch")
		allowCurrent, _ := cmd.Flags().GetBool("allow-current")

		// 创建 Git 服务
		gitService := service.NewGitService(nil, nil, "")

		// 分析提交
		release, err := gitService.AnalyzeCommits()
		if err != nil {
			return err
		}

		// 检查是否有更改
		if !release.HasContent() {
			if allowCurrent {
				fmt.Println(release.Version.Current.String())
				return nil
			}
			if bumpPatch {
				release.Version.Bump(domain.BumpPatch)
				fmt.Println(release.Version.Next.String())
				return nil
			}
			return fmt.Errorf("没有检测到更改")
		}

		// 确保在初始开发期间遵循版本规则
		if initialDevelopment, _ := cmd.Flags().GetBool("initial-development"); initialDevelopment {
			if release.Version.Next.Major > 0 {
				release.Version.Next = semver.MustParse("0.1.0")
			}
		}

		fmt.Println(release.Version.Next.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nextVersionCmd)

	// 命令特定选项
	nextVersionCmd.Flags().Bool("bump-patch", false, "当没有提交会触发版本更新时强制增加补丁版本")
	nextVersionCmd.Flags().Bool("allow-current", false, "如果没有检测到更改，允许打印当前版本")
	nextVersionCmd.MarkFlagsMutuallyExclusive("bump-patch", "allow-current")
}
