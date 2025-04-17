package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/fanny7d/semrel-gitlab/pkg/service"
	"github.com/spf13/cobra"
)

var commitAndTagCmd = &cobra.Command{
	Use:   "commit-and-tag",
	Short: "提交文件并标记新提交",
	Long: `提交并推送列出的文件。
如果文件不包含任何更改，命令将失败。

默认的提交消息模板包含 [skip ci]，
以防止提交管道运行。你可以使用全局选项
--bump-commit-tmpl 或环境变量 GSG_BUMP_COMMIT_TMPL
覆盖默认模板。

为新提交创建标签和发布说明
(查看 'release help tag' 获取更多详细信息)。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取全局选项
		token := cmd.Flag("token").Value.String()
		if token == "" {
			return fmt.Errorf("token 是必需的")
		}

		apiURL := cmd.Flag("gl-api").Value.String()
		if apiURL == "" {
			return fmt.Errorf("gl-api 是必需的")
		}

		project := cmd.Flag("ci-project-path").Value.String()
		if project == "" {
			return fmt.Errorf("ci-project-path 是必需的")
		}

		projectURLStr := cmd.Flag("ci-project-url").Value.String()
		if projectURLStr == "" {
			return fmt.Errorf("ci-project-url 是必需的")
		}

		projectURL, err := url.Parse(projectURLStr)
		if err != nil {
			return fmt.Errorf("解析 project-url 失败: %v", err)
		}

		skipSSLVerify, _ := cmd.Flags().GetBool("skip-ssl-verify")
		patchTypes := strings.Split(cmd.Flag("patch-commit-types").Value.String(), ",")
		minorTypes := strings.Split(cmd.Flag("minor-commit-types").Value.String(), ",")
		tagPrefix := cmd.Flag("tag-prefix").Value.String()
		branch := cmd.Flag("ci-commit-ref-name").Value.String()
		commitTmpl := cmd.Flag("bump-commit-tmpl").Value.String()

		// 创建服务
		gitService := service.NewGitService(patchTypes, minorTypes, tagPrefix)
		gitlabService, err := service.NewGitLabService(token, apiURL, project, projectURL, skipSSLVerify)
		if err != nil {
			return err
		}

		// 分析提交
		release, err := gitService.AnalyzeCommits()
		if err != nil {
			return err
		}

		// 检查是否有变更
		if !release.HasContent() {
			return fmt.Errorf("提交日志中没有发现会改变版本的变更")
		}

		// 渲染提交消息
		tmpl, err := template.New("commit").Parse(commitTmpl)
		if err != nil {
			return fmt.Errorf("解析提交消息模板失败: %v", err)
		}

		var message strings.Builder
		err = tmpl.Execute(&message, map[string]string{
			"tag": release.TagName,
		})
		if err != nil {
			return fmt.Errorf("渲染提交消息失败: %v", err)
		}

		// 创建提交
		if err := gitlabService.CreateCommit(branch, message.String(), args); err != nil {
			return err
		}

		// 渲染发布说明
		renderService := service.NewRenderService("")
		if err := renderService.RenderReleaseNote(release); err != nil {
			return err
		}

		// 创建标签
		if err := gitlabService.CreateTag(release, branch); err != nil {
			return err
		}

		// 如果需要，创建管道
		createTagPipeline, _ := cmd.Flags().GetBool("create-tag-pipeline")
		if createTagPipeline && strings.Contains(message.String(), "[skip ci]") {
			if err := gitlabService.CreatePipeline(branch); err != nil {
				return err
			}
		}

		fmt.Printf("已创建标签 %s\n", release.TagName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitAndTagCmd)

	// 命令特定选项
	commitAndTagCmd.Flags().Bool("create-tag-pipeline", false, "当标记的提交消息包含 [skip ci] 并且你想要执行标签管道时需要")
	commitAndTagCmd.Flags().Bool("list-other-changes", false, "列出不影响版本控制的更改")
}
