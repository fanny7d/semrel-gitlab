# 常见问题

## 一般问题

### Q: semrel-gitlab 是什么？
A: semrel-gitlab 是一个基于语义化版本的 GitLab 发布工具，它可以自动分析提交记录，根据提交信息自动确定版本号，并在 GitLab 上创建相应的标签和发布说明。

### Q: 为什么要使用 semrel-gitlab？
A: 使用 semrel-gitlab 可以：
- 自动化版本管理流程
- 确保版本号的一致性
- 自动生成变更日志
- 简化发布流程
- 提高开发效率

## 安装问题

### Q: 如何安装 semrel-gitlab？
A: 有两种安装方式：
1. 使用 Go 安装：`go install github.com/fanny7d/semrel-gitlab@latest`
2. 使用 Docker：`docker pull fanny7d/semrel-gitlab:latest`

### Q: 安装时出现 "package not found" 错误怎么办？
A: 请确保：
1. Go 版本 >= 1.22
2. GOPATH 已正确设置
3. 使用了正确的包名

## 配置问题

### Q: 如何配置 GitLab 访问令牌？
A: 有两种方式：
1. 环境变量：`export GITLAB_TOKEN=your-token`
2. 命令行参数：`--token your-token`

### Q: 如何跳过 SSL 验证？
A: 有三种方式：
1. 环境变量：`export GITLAB_SKIP_SSL_VERIFY=true`
2. 命令行参数：`--skip-ssl-verify`
3. 配置文件：
   ```yaml
   gitlab:
     skip_ssl_verify: true
   ```

## 使用问题

### Q: 如何确定下一个版本号？
A: 使用 `next-version` 命令：
```bash
semrel-gitlab next-version
```

### Q: 如何生成变更日志？
A: 使用 `changelog` 命令：
```bash
semrel-gitlab changelog
```

### Q: 如何创建标签和发布？
A: 使用 `tag` 命令：
```bash
semrel-gitlab tag
```

### Q: 如何在提交时自动创建标签？
A: 使用 `commit-and-tag` 命令：
```bash
semrel-gitlab commit-and-tag file1 file2
```

### Q: 如何添加下载文件到发布？
A: 使用 `add-download` 命令：
```bash
semrel-gitlab add-download --file path/to/file
```

## 版本控制问题

### Q: 什么时候会增加主版本号？
A: 当提交包含破坏性变更时（在提交消息中包含 `BREAKING CHANGE:`）。

### Q: 什么时候会增加次版本号？
A: 当提交类型为 `feat` 时。

### Q: 什么时候会增加补丁版本号？
A: 当提交类型为 `fix`、`refactor`、`perf`、`docs`、`style` 或 `test` 时。

### Q: 如何处理预发布版本？
A: 使用 `--pre-tmpl` 参数或在配置文件中设置 `version.pre_templates`：
```yaml
version:
  pre_templates:
    - alpha.{{.Num}}
    - beta.{{.Num}}
    - rc.{{.Num}}
```

## CI/CD 问题

### Q: 如何在 GitLab CI 中使用？
A: 在 `.gitlab-ci.yml` 中添加：
```yaml
release:
  stage: release
  script:
    - semrel-gitlab release
  only:
    - main
```

### Q: 如何避免发布时触发新的 CI 管道？
A: 在提交消息中添加 `[skip ci]`，可以通过配置文件设置：
```yaml
ci:
  bump_commit_template: "chore: 版本更新为 {{.Version}} [skip ci]"
```

## 错误处理

### Q: 提示 "token is required" 怎么办？
A: 需要设置 GitLab 访问令牌，可以：
1. 设置环境变量：`export GITLAB_TOKEN=your-token`
2. 使用命令行参数：`--token your-token`

### Q: 提示 "no changes found" 怎么办？
A: 这意味着没有找到会触发版本更新的提交。可以：
1. 确保提交消息遵循 Conventional Commits 规范
2. 使用 `--bump-patch` 强制增加补丁版本
3. 使用 `--allow-current` 允许使用当前版本

### Q: SSL 证书验证失败怎么办？
A: 如果使用自签名证书，可以：
1. 设置环境变量：`export GITLAB_SKIP_SSL_VERIFY=true`
2. 使用命令行参数：`--skip-ssl-verify`

## 其他问题

### Q: 如何查看工具版本？
A: 使用 `version` 命令：
```bash
semrel-gitlab version
```

### Q: 如何启用自动补全？
A: 使用 `completion` 命令生成补全脚本：
```bash
# Bash
source <(semrel-gitlab completion bash)

# Zsh
source <(semrel-gitlab completion zsh)

# Fish
semrel-gitlab completion fish | source

# PowerShell
semrel-gitlab completion powershell | Out-String | Invoke-Expression
```

### Q: 在哪里可以获取更多帮助？
A: 可以：
1. 使用 `--help` 查看命令帮助
2. 查看[文档](docs/)
3. 提交 [Issue](https://gitlab.com/fanny7d/semrel-gitlab/issues)
4. 查看[示例](examples/) 