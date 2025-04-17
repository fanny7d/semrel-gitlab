# 配置文件说明

semrel-gitlab 支持使用配置文件来设置默认选项，避免每次都需要在命令行中指定。配置文件使用 YAML 格式。

## 配置文件位置

工具会按以下顺序查找配置文件：

1. 命令行参数 `--config` 指定的文件
2. 当前目录下的 `.semrelrc.yml`
3. 用户主目录下的 `.semrelrc.yml`

## 配置项

```yaml
# GitLab 相关配置
gitlab:
  # API URL，默认为 https://gitlab.com/api/v4
  api_url: https://gitlab.example.com/api/v4
  # 是否跳过 SSL 验证
  skip_ssl_verify: false

# 版本控制配置
version:
  # 版本号前缀
  tag_prefix: v
  # 初始开发阶段标志
  initial_development: true
  # 预发布版本模板
  pre_templates:
    - alpha.{{.Num}}
    - beta.{{.Num}}
    - rc.{{.Num}}
  # 构建元数据模板
  build_templates:
    - build.{{.Timestamp}}
    - sha.{{.CommitSHA}}

# 提交分析配置
commit:
  # 补丁版本更新的提交类型
  patch_types:
    - fix
    - refactor
    - perf
    - docs
    - style
    - test
  # 次要版本更新的提交类型
  minor_types:
    - feat
  # 忽略的提交类型
  ignore_types:
    - chore
    - ci
    - build

# 发布说明配置
release:
  # 发布说明模板文件
  template: .github/release-template.md
  # 是否包含作者信息
  include_authors: true
  # 是否包含提交链接
  include_links: true
  # 变更类型分组
  groups:
    - title: "🚀 新功能"
      types: [feat]
    - title: "🐛 修复"
      types: [fix]
    - title: "♻️ 重构"
      types: [refactor]
    - title: "⚡️ 性能优化"
      types: [perf]
    - title: "📚 文档"
      types: [docs]
    - title: "💄 样式"
      types: [style]
    - title: "✅ 测试"
      types: [test]
    - title: "🔧 构建系统"
      types: [chore]

# CI/CD 配置
ci:
  # 发布分支
  release_branches:
    - main
    - master
  # 预发布分支
  prerelease_branches:
    - develop
    - staging
  # 版本更新提交消息模板
  bump_commit_template: "chore: 版本更新为 {{.Version}} [skip ci]"
```

## 环境变量

配置文件中的所有选项都可以通过环境变量覆盖。环境变量的命名规则是将配置路径转换为大写，并用下划线连接。例如：

- `GITLAB_API_URL` 对应 `gitlab.api_url`
- `VERSION_TAG_PREFIX` 对应 `version.tag_prefix`
- `COMMIT_PATCH_TYPES` 对应 `commit.patch_types`

## 模板变量

在模板中可以使用以下变量：

### 版本相关
- `{{.Version}}`: 完整的版本号
- `{{.Major}}`: 主版本号
- `{{.Minor}}`: 次版本号
- `{{.Patch}}`: 补丁版本号
- `{{.PreRelease}}`: 预发布版本标识符
- `{{.Build}}`: 构建元数据

### 提交相关
- `{{.CommitSHA}}`: 提交的完整哈希值
- `{{.CommitShortSHA}}`: 提交的短哈希值
- `{{.CommitMessage}}`: 提交消息
- `{{.CommitAuthor}}`: 提交作者
- `{{.CommitDate}}`: 提交日期

### 其他
- `{{.Timestamp}}`: 当前时间戳
- `{{.Date}}`: 当前日期
- `{{.Time}}`: 当前时间
- `{{.Num}}`: 序号（在预发布版本中使用）
- `{{.Env "VAR"}}`: 环境变量值

## 示例

### 最小配置

```yaml
gitlab:
  api_url: https://gitlab.example.com/api/v4

version:
  tag_prefix: v
  initial_development: true

commit:
  patch_types: [fix, refactor, perf, docs, style, test]
  minor_types: [feat]
```

### 完整配置

```yaml
gitlab:
  api_url: https://gitlab.example.com/api/v4
  skip_ssl_verify: false

version:
  tag_prefix: v
  initial_development: true
  pre_templates:
    - alpha.{{.Num}}
    - beta.{{.Num}}
    - rc.{{.Num}}
  build_templates:
    - build.{{.Timestamp}}
    - sha.{{.CommitSHA}}

commit:
  patch_types:
    - fix
    - refactor
    - perf
    - docs
    - style
    - test
  minor_types:
    - feat
  ignore_types:
    - chore
    - ci
    - build

release:
  template: .github/release-template.md
  include_authors: true
  include_links: true
  groups:
    - title: "🚀 新功能"
      types: [feat]
    - title: "🐛 修复"
      types: [fix]
    - title: "♻️ 重构"
      types: [refactor]
    - title: "⚡️ 性能优化"
      types: [perf]
    - title: "📚 文档"
      types: [docs]
    - title: "💄 样式"
      types: [style]
    - title: "✅ 测试"
      types: [test]
    - title: "🔧 构建系统"
      types: [chore]

ci:
  release_branches:
    - main
    - master
  prerelease_branches:
    - develop
    - staging
  bump_commit_template: "chore: 版本更新为 {{.Version}} [skip ci]"
``` 