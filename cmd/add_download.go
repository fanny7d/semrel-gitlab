package cmd

import (
	"fmt"
	"net/url"

	"github.com/fanny7d/semrel-gitlab/pkg/service"
	"github.com/spf13/cobra"
)

var addDownloadCmd = &cobra.Command{
	Use:   "add-download",
	Short: "添加下载到发布说明",
	Long: `上传文件到项目上传并添加链接到发布说明。
            
需要 CI_COMMIT_TAG 环境变量或 --ci-commit-tag 标志。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取命令选项
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			return fmt.Errorf("file 是必需的")
		}

		tag, _ := cmd.Flags().GetString("ci-commit-tag")
		if tag == "" {
			return fmt.Errorf("ci-commit-tag 是必需的")
		}

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

		// 创建服务
		gitlabService, err := service.NewGitLabService(token, apiURL, project, projectURL, skipSSLVerify)
		if err != nil {
			return err
		}

		// 获取标签
		release, err := gitlabService.GetTag(tag)
		if err != nil {
			return err
		}

		// 上传文件
		if err := gitlabService.UploadFile(release, file); err != nil {
			return err
		}

		fmt.Printf("已添加下载链接到标签 %s\n", tag)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addDownloadCmd)

	// 命令特定选项
	addDownloadCmd.Flags().StringP("file", "f", "", "要上传的文件")
	addDownloadCmd.Flags().String("ci-commit-tag", "", "要添加下载的标签")
}
