# semrel-gitlab

这是一个用于 GitLab 的语义化发布工具，基于 [go-semrel](https://github.com/go-semrel/semrel) 开发。

## 功能特点

- 支持语义化版本控制
- 自动生成变更日志
- 支持 GitLab 标签管理
- 支持文件上传和链接生成
- 支持提交操作管理
- 专为 GitLab CI 流水线设计
- 提供 Docker 镜像和单文件二进制，易于集成到流水线中

## 主要功能

- 根据提交信息自动确定下一个版本号
- 创建/更新变更日志
- 创建带发布说明的标签
- 为发布说明附加文件
- 提交版本更新（如 `package.json`、`pom.xml`、`CHANGELOG.md` 等）

## 安装

```bash
go get gitlab.com/fanny7d/semrel-gitlab
```

## 使用方法

### 基本配置

1. 设置 GitLab 访问令牌：
```bash
export GITLAB_TOKEN=your_token_here
```

2. 设置 GitLab API URL（可选）：
```bash
export GITLAB_API_URL=https://gitlab.example.com/api/v4
```

### 命令行参数

- `-token`: GitLab 访问令牌（也可以通过环境变量 `GITLAB_TOKEN` 设置）
- `-api-url`: GitLab API URL（也可以通过环境变量 `GITLAB_API_URL` 设置）
- `-skip-ssl-verify`: 跳过 SSL 验证（默认为 false）
- `-project`: GitLab 项目路径（例如：`group/project`）
- `-branch`: 要发布的分支（默认为 `master`）
- `-force`: 强制创建标签（默认为 false）
- `-message`: 标签消息（默认为空）
- `-files`: 要上传的文件列表（逗号分隔）
- `-link-description`: 链接描述（默认为空）

### 示例

1. 创建标签并上传文件：
```bash
semrel-gitlab -token $GITLAB_TOKEN -project group/project -branch master -force true -message "Release v1.0.0" -files file1.txt,file2.txt
```

2. 仅创建标签：
```bash
semrel-gitlab -token $GITLAB_TOKEN -project group/project -branch master -force true -message "Release v1.0.0"
```

## 开发指南

### 项目结构

- `pkg/actions/`: 包含各种 GitLab 操作的具体实现
  - `actions.go`: 标签和链接相关操作
  - `commit.go`: 提交相关操作
  - `upload.go`: 文件上传相关操作
- `pkg/gitlabutil/`: GitLab 工具函数
  - `client.go`: GitLab 客户端相关
  - `update_release.go`: 发布更新相关

### 主要功能模块

1. 标签管理
   - 创建标签
   - 获取标签
   - 更新标签描述

2. 文件操作
   - 文件上传
   - 生成文件链接
   - 生成 Markdown 格式链接

3. 提交管理
   - 创建提交
   - 文件状态映射
   - 提交操作生成

## 更多信息

- [使用指南和命令参考](https://fanny7d.gitlab.io/semrel-gitlab)

## 待办事项

- 完善文档
- 可配置的发布说明

## 许可证

MIT License
