# 使用说明

## 环境要求

- Go 1.22 或更高版本
- GitLab 实例（支持自托管或 gitlab.com）
- GitLab 访问令牌（需要 API 权限）

## 配置

### GitLab 访问令牌

你需要一个具有 API 权限的 GitLab 访问令牌。可以通过以下方式配置：

1. 环境变量（推荐）：
```bash
export GITLAB_TOKEN=your-token
```

2. 命令行参数：
```bash
semrel-gitlab release --token your-token
```

### GitLab API URL

如果你使用自托管的 GitLab 实例，需要配置 API URL：

1. 环境变量：
```bash
export GITLAB_API_URL=https://gitlab.example.com
```

2. 命令行参数：
```bash
semrel-gitlab release --api-url https://gitlab.example.com
```

### SSL 验证

对于自签名证书的 GitLab 实例，可以选择跳过 SSL 验证：

```bash
semrel-gitlab release --skip-ssl-verify
```

## 基本用法

### 创建发布

在项目目录下运行：

```bash
semrel-gitlab release
```

这将：
1. 分析提交历史
2. 确定下一个版本号
3. 创建标签
4. 生成发布说明
5. 创建 GitLab 发布

### 指定版本号

手动指定版本号：

```bash
semrel-gitlab release --version 1.2.3
```

### 自定义发布说明

使用自定义模板：

```bash
semrel-gitlab release --template path/to/template.md
```

模板示例：
```markdown
# 版本 {{.Version}}

## 新特性
{{range .Features}}
- {{.}}
{{end}}

## 修复
{{range .Fixes}}
- {{.}}
{{end}}

## 其他更改
{{range .Others}}
- {{.}}
{{end}}
```

### 发布文件

上传文件到发布：

```bash
semrel-gitlab release --files dist/*
```

### 预览模式

在不实际创建发布的情况下预览结果：

```bash
semrel-gitlab release --dry-run
```

## 提交消息格式

工具使用提交消息来确定版本号和生成发布说明。建议使用以下格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型（type）可以是：
- feat: 新特性（触发次版本更新）
- fix: 修复（触发修订版本更新）
- BREAKING CHANGE: 重大更改（触发主版本更新）
- chore: 维护性更改
- docs: 文档更改
- style: 代码格式更改
- refactor: 重构
- test: 测试相关
- ci: CI 相关

示例：
```
feat(api): 添加用户认证接口

- 实现 JWT 认证
- 添加用户登录接口
- 添加用户注册接口

BREAKING CHANGE: 认证方式从 session 改为 JWT
```

## 持续集成

### GitLab CI 示例

```yaml
release:
  stage: release
  script:
    - export GITLAB_TOKEN=${CI_JOB_TOKEN}
    - semrel-gitlab release
  only:
    - master
``` 