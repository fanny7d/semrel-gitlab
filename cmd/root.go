package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "DEV"
)

var rootCmd = &cobra.Command{
	Use:   "semrel-gitlab",
	Short: "GitLab 语义化发布工具",
	Long: `一组用于帮助自动化发布的工具

推荐的使用方式是在 Gitlab CI 变量中为相关环境变量赋值。
请注意 ci-* 选项是由 Gitlab CI 自动填充的。

支持的功能：
- 自动版本号管理
- 变更日志生成
- GitLab 发布管理
- 多平台构建支持
- 预发布版本支持`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 验证必要的环境变量
		if token, _ := cmd.Flags().GetString("token"); token == "" {
			return fmt.Errorf("必须提供 GitLab 访问令牌")
		}

		// 验证 CI 环境变量
		if os.Getenv("CI_PROJECT_PATH") == "" {
			return fmt.Errorf("必须提供 CI_PROJECT_PATH 环境变量")
		}

		if os.Getenv("CI_COMMIT_SHA") == "" {
			return fmt.Errorf("必须提供 CI_COMMIT_SHA 环境变量")
		}

		return nil
	},
}

func Execute() {
	cobra.AddTemplateFunc("translate", func(s string) string {
		translations := map[string]string{
			"completion":                "生成指定 shell 的自动补全脚本",
			"help":                      "获取任意命令的帮助信息",
			"version for semrel-gitlab": "semrel-gitlab 的版本信息",
			"Generate the autocompletion script for the specified shell": "生成指定 shell 的自动补全脚本",
			"Help about any command":                                     "获取任意命令的帮助信息",
		}
		if t, ok := translations[s]; ok {
			return t
		}
		return s
	})

	rootCmd.SetUsageTemplate(`用法:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [命令]{{end}}{{if gt (len .Aliases) 0}}

别名:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

示例:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

可用命令:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{translate .Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

标志:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

全局标志:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

其他帮助主题:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

使用 "{{.CommandPath}} [命令] --help" 获取更多关于命令的信息。{{end}}
`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 全局选项
	rootCmd.PersistentFlags().StringP("token", "t", "", "GitLab 私有令牌 (必需)")
	rootCmd.PersistentFlags().String("gl-api", os.Getenv("CI_API_V4_URL"), "GitLab API URL。如果未定义，则使用 CI_API_V4_URL 环境变量")
	rootCmd.PersistentFlags().Bool("skip-ssl-verify", false, "不验证 GitLab API 的 CA 证书")
	rootCmd.PersistentFlags().String("patch-commit-types", "fix,refactor,perf,docs,style,test", "逗号分隔的提交消息类型列表，表示补丁版本更新")
	rootCmd.PersistentFlags().String("minor-commit-types", "feat", "逗号分隔的提交消息类型列表，表示次要版本更新")
	rootCmd.PersistentFlags().Bool("initial-development", true, "当你准备发布 1.0.0 时设置为 false，如果版本已经 >= 1.0.0 则忽略")
	rootCmd.PersistentFlags().Bool("bump-patch", false, "当没有提交会触发版本更新时强制增加补丁版本")
	rootCmd.PersistentFlags().String("release-branches", "main,master", "逗号分隔的分支名称列表")
	rootCmd.PersistentFlags().String("tag-prefix", "v", "版本标签使用的前缀")
	rootCmd.PersistentFlags().String("bump-commit-tmpl", "chore: 版本更新为 {{tag}} [skip ci]", "版本更新提交消息的模板")
	rootCmd.PersistentFlags().String("pre-tmpl", "", "预发布版本模板。逗号分隔的 ID 模板列表")
	rootCmd.PersistentFlags().String("build-tmpl", "", "构建元数据模板。逗号分隔的 ID 模板列表")

	// 从环境变量中读取默认值
	rootCmd.PersistentFlags().SetAnnotation("token", cobra.BashCompOneRequiredFlag, []string{"true"})
}
