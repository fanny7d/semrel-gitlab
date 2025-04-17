# semrel-gitlab

基于语义化版本的 GitLab 发布工具

## 简介

semrel-gitlab 是一个用于 GitLab 项目的语义化版本发布工具。它可以自动分析提交记录，根据提交信息自动确定版本号，并在 GitLab 上创建相应的标签和发布说明。

## 特性

- 自动分析提交记录，确定版本号
- 支持自定义 GitLab 实例
- 支持 SSL 验证配置
- 自动生成发布说明
- 支持自定义发布说明模板
- 支持多种 shell 的自动补全
- 支持预发布版本和构建元数据
- 支持多平台构建

## 安装

### 使用 Go 安装

```bash
go install github.com/fanny7d/semrel-gitlab@latest
```

### 使用 Docker 安装

```bash
docker pull fanny7d/semrel-gitlab:latest
```

## 快速开始

1. 设置 GitLab 访问令牌：

```bash
export GITLAB_TOKEN=your-token
```

2. 在项目中运行：

```bash
semrel-gitlab release
```

## 使用示例

### 检查下一个版本号

```bash
semrel-gitlab next-version
```

### 生成变更日志

```bash
semrel-gitlab changelog
```

### 创建标签和发布

```bash
semrel-gitlab tag
```

### 提交并创建标签

```bash
semrel-gitlab commit-and-tag README.md
```

### 添加下载文件到发布

```bash
semrel-gitlab add-download --file dist/app.tar.gz
```

## 配置

### 环境变量

- `GITLAB_TOKEN`: GitLab 访问令牌（必需）
- `GITLAB_API_URL`: GitLab API URL（可选，默认为 https://gitlab.com/api/v4）
- `GITLAB_SKIP_SSL_VERIFY`: 是否跳过 SSL 验证（可选，默认为 false）

### 命令行选项

所有命令都支持以下全局选项：

- `--token, -t`: GitLab 访问令牌
- `--gl-api`: GitLab API URL
- `--skip-ssl-verify`: 不验证 GitLab API 的 CA 证书
- `--patch-commit-types`: 补丁版本更新的提交类型
- `--minor-commit-types`: 次要版本更新的提交类型
- `--initial-development`: 初始开发阶段标志
- `--tag-prefix`: 版本标签前缀
- `--pre-tmpl`: 预发布版本模板
- `--build-tmpl`: 构建元数据模板

## 提交消息格式

工具使用 Conventional Commits 规范来分析提交消息。支持的类型包括：

- `fix`: 修复 bug（补丁版本）
- `feat`: 新功能（次要版本）
- `refactor`: 重构（补丁版本）
- `perf`: 性能优化（补丁版本）
- `docs`: 文档更新（补丁版本）
- `style`: 代码格式调整（补丁版本）
- `test`: 测试相关（补丁版本）
- `chore`: 构建过程或辅助工具的变动（补丁版本）

破坏性变更（在提交消息中包含 `BREAKING CHANGE:`）会触发主要版本更新。

## 自动补全

工具支持为多种 shell 生成自动补全脚本：

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

## 文档

- [使用说明](docs/usage.md)
- [命令参数](docs/commands.md)
- [配置文件](docs/config.md)
- [常见问题](docs/faq.md)

## 贡献

欢迎提交 Pull Request 和 Issue。在提交之前，请：

1. 确保代码通过测试
2. 更新相关文档
3. 遵循项目的代码风格

## 许可证

MIT License