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

## 安装

```bash
go install github.com/fanny7d/semrel-gitlab@latest
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

## 文档

- [使用说明](docs/usage.md)
- [命令参数](docs/commands.md)

## 许可证

MIT License