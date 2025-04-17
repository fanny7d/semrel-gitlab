package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "生成自动补全脚本",
	Long: `为 semrel-gitlab 生成自动补全脚本。

支持以下 shell：
  - bash
  - zsh
  - fish
  - powershell

使用方法:

Bash:
  $ source <(semrel-gitlab completion bash)

  # 永久启用自动补全，需要将上述命令添加到 .bashrc 文件中：
  $ semrel-gitlab completion bash > ~/.bash_completion.d/semrel-gitlab

Zsh:
  # 如果 shell 补全尚未启用，需要执行以下命令：
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # 然后加载 semrel-gitlab 补全：
  $ source <(semrel-gitlab completion zsh)

  # 永久启用自动补全，需要将补全脚本复制到补全目录：
  $ semrel-gitlab completion zsh > "${fpath[1]}/_semrel-gitlab"

Fish:
  $ semrel-gitlab completion fish | source

  # 永久启用自动补全：
  $ semrel-gitlab completion fish > ~/.config/fish/completions/semrel-gitlab.fish

PowerShell:
  PS> semrel-gitlab completion powershell | Out-String | Invoke-Expression

  # 永久启用自动补全：
  PS> semrel-gitlab completion powershell > semrel-gitlab.ps1
  # 然后将 semrel-gitlab.ps1 添加到 PowerShell 配置文件中`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
