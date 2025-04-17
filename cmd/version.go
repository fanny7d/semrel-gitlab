package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long: `显示 semrel-gitlab 的版本信息。

输出包括：
- 版本号
- Git 提交哈希
- 构建时间
- Go 版本
- 操作系统/架构`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("semrel-gitlab 版本 %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
