package cmd

import (
	"fmt"
	"strings"

	"github.com/fanny7d/semrel-gitlab/pkg/service"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "创建标签和发布",
	Long: `分析提交信息，创建标签并发布到 GitLab。

此命令将：
1. 分析提交信息
2. 确定下一个版本号
3. 创建 Git 标签
4. 在 GitLab 上创建发布
5. 添加发布说明和下载链接`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取命令选项
		listOtherChanges, _ := cmd.Flags().GetBool("list-other-changes")
		ciCommitTag, _ := cmd.Flags().GetString("ci-commit-tag")

		// 获取全局选项
		patchTypes := strings.Split(cmd.Flag("patch-commit-types").Value.String(), ",")
		minorTypes := strings.Split(cmd.Flag("minor-commit-types").Value.String(), ",")
		tagPrefix := cmd.Flag("tag-prefix").Value.String()
		token, _ := cmd.Flags().GetString("token")
		glAPI, _ := cmd.Flags().GetString("gl-api")
		skipSSLVerify, _ := cmd.Flags().GetBool("skip-ssl-verify")
		ciProjectPath, _ := cmd.Flags().GetString("ci-project-path")

		// 创建 Git 服务
		gitService := service.NewGitService(patchTypes, minorTypes, tagPrefix)

		// 分析提交
		release, err := gitService.AnalyzeCommits()
		if err != nil {
			return err
		}

		// 检查是否有变更
		if !release.HasContent() && !listOtherChanges {
			return fmt.Errorf("提交日志中没有发现会改变版本的变更")
		}

		// 创建 GitLab 客户端
		client, err := service.NewGitLabClient(token, glAPI, skipSSLVerify)
		if err != nil {
			return fmt.Errorf("创建 GitLab 客户端失败: %v", err)
		}

		// 创建标签
		tagName := ciCommitTag
		if tagName == "" {
			tagName = release.TagName
		}

		// 创建 Git 标签
		if err := gitService.CreateTag(tagName); err != nil {
			return fmt.Errorf("创建 Git 标签失败: %v", err)
		}

		// 创建 GitLab 发布
		if err := client.CreateRelease(ciProjectPath, tagName, release); err != nil {
			return fmt.Errorf("创建 GitLab 发布失败: %v", err)
		}

		fmt.Printf("已创建标签 %s 并发布到 GitLab\n", tagName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

	// 命令特定选项
	tagCmd.Flags().Bool("list-other-changes", false, "列出不影响版本控制的更改")
}
