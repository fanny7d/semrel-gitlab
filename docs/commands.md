# 命令参数说明

## 全局选项

这些选项适用于所有命令：

| 选项 | 环境变量 | 说明 | 默认值 |
|------|----------|------|--------|
| `--token` | `GITLAB_TOKEN` | GitLab 访问令牌 | - |
| `--api-url` | `GITLAB_API_URL` | GitLab API URL | https://gitlab.com |
| `--skip-ssl-verify` | `GITLAB_SKIP_SSL_VERIFY` | 跳过 SSL 验证 | false |
| `--debug` | `GITLAB_DEBUG` | 启用调试输出 | false |

## release 命令

用于创建新的发布。

### 基本用法

```bash
semrel-gitlab release [选项]
```

### 选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `--version` | 指定版本号（不指定则自动确定） | - |
| `--branch` | 目标分支 | master |
| `--template` | 发布说明模板文件路径 | - |
| `--files` | 要上传的文件（支持 glob 模式） | - |
| `--dry-run` | 预览模式，不实际创建发布 | false |
| `--force` | 强制创建发布，即使版本已存在 | false |
| `--no-merge-commits` | 忽略合并提交 | false |
| `--tag-prefix` | 标签前缀 | v |
| `--prerelease` | 预发布标识符 | - |
| `--build` | 构建元数据 | - |

### 示例

1. 基本用法：
```bash
semrel-gitlab release
```

2. 指定版本号：
```bash
semrel-gitlab release --version 1.2.3
```

3. 使用自定义模板：
```bash
semrel-gitlab release --template release-notes.md
```

4. 上传文件：
```bash
semrel-gitlab release --files "dist/*.tar.gz" --files "dist/*.zip"
```

5. 预览模式：
```bash
semrel-gitlab release --dry-run
```

6. 预发布版本：
```bash
semrel-gitlab release --prerelease beta
```

7. 完整示例：
```bash
semrel-gitlab release \
  --token your-token \
  --api-url https://gitlab.example.com \
  --branch develop \
  --template custom-template.md \
  --files "dist/*" \
  --tag-prefix release- \
  --prerelease beta.1 \
  --build "20240318.1"
```

## version 命令

显示工具版本信息。

### 用法

```bash
semrel-gitlab version
```

## init 命令

初始化项目配置。

### 用法

```bash
semrel-gitlab init [选项]
```

### 选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `--template` | 发布说明模板 | default |
| `--commit-format` | 提交消息格式 | conventional |

### 示例

```bash
semrel-gitlab init --template custom
```

## 配置文件

工具支持使用配置文件（`.semrelrc.yml`）来设置默认选项：

```yaml
gitlab:
  api_url: https://gitlab.example.com
  skip_ssl_verify: false

release:
  branch: master
  tag_prefix: v
  template: .github/release-template.md
  files:
    - dist/*.tar.gz
    - dist/*.zip
  
commit:
  format: conventional
  ignore_merge_commits: true

ci:
  dry_run: false
  debug: false
``` 