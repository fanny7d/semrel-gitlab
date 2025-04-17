# semrel-gitlab

semrel-gitlab 是一个用于自动化 GitLab 发布流程的工具。它可以帮助你：

- 自动创建和管理 GitLab 标签
- 上传发布文件
- 生成发布说明
- 管理发布链接
- 自动化发布流程

## 安装

```bash
go get github.com/fanny7d/semrel-gitlab
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

## API 文档

### 核心功能

#### 标签管理

- `CreateTag`: 创建新的 GitLab 标签
- `GetTag`: 获取现有的 GitLab 标签
- `UpdateTagDescription`: 更新标签描述

#### 文件操作

- `UploadFile`: 上传文件到 GitLab
- `GenerateFileLink`: 生成文件链接
- `GenerateMarkdownLink`: 生成 Markdown 格式链接

#### 发布管理

- `CreateRelease`: 创建新的发布
- `UpdateRelease`: 更新现有发布
- `AddReleaseLink`: 添加发布链接

### 工作流

#### 基本工作流

1. 分析提交信息，确定版本号
2. 创建标签
3. 上传文件（如果需要）
4. 添加链接（如果需要）
5. 生成发布说明

#### 错误处理

- 所有操作都是幂等的
- 支持自动重试（对于 502 错误）
- 支持回滚操作

## 开发指南

### 项目结构

- `pkg/actions/`: 包含各种 GitLab 操作的具体实现
  - `actions.go`: 标签和链接相关操作
  - `commit.go`: 提交相关操作
  - `upload.go`: 文件上传相关操作
- `pkg/gitlabutil/`: GitLab 工具函数
  - `client.go`: GitLab 客户端相关
  - `update_release.go`: 发布更新相关
- `pkg/workflow/`: 工作流管理
  - `workflow.go`: 工作流核心逻辑

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

## 最佳实践

1. 版本控制
   - 使用语义化版本控制
   - 遵循 Git 提交规范
   - 使用标签管理发布

2. 错误处理
   - 使用 Go 的错误处理机制
   - 提供详细的错误信息
   - 支持自动重试

3. 代码组织
   - 使用清晰的包结构
   - 遵循 Go 的命名规范
   - 提供完整的文档

4. 测试
   - 编写单元测试
   - 使用测试覆盖率工具
   - 进行持续集成

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License