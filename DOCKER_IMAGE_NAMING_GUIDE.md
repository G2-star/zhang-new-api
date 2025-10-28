# Docker 镜像构建和发布指南

## 📦 镜像命名方案

### 推荐命名格式

```
<组织名/用户名>/<项目名>:<版本标签>
```

### 具体示例

```bash
# 方案1：简洁版
your-name/new-api:latest

# 方案2：带功能标识（推荐）
your-name/new-api:conversation-v1.0.0

# 方案3：完整描述
your-name/new-api-conversation:v1.0.0
```

---

## 🏗️ 构建 Docker 镜像

### 方法 1：使用现有 Dockerfile

```bash
# 进入项目目录
cd D:\Users\Zhang\Desktop\new-api

# 构建镜像（多标签）
docker build \
  -t your-name/new-api:latest \
  -t your-name/new-api:v1.0.0 \
  -t your-name/new-api:conversation \
  -t your-name/new-api:v1.0.0-conversation \
  .

# 示例：如果你的名字是 zhang
docker build \
  -t zhang/new-api:latest \
  -t zhang/new-api:v1.0.0-conversation \
  .
```

### 方法 2：修改 Dockerfile（可选）

项目应该已有 Dockerfile，如需自定义，可以添加标签：

```dockerfile
# 在 Dockerfile 开头添加元数据
LABEL maintainer="your-email@example.com"
LABEL version="1.0.0"
LABEL description="New API with conversation logging feature"
LABEL features="conversation-log,auto-archive,compression"

# 其余内容保持不变
```

---

## 🚀 发布到 Docker Hub

### 步骤 1：登录 Docker Hub

```bash
docker login

# 输入你的 Docker Hub 用户名和密码
```

### 步骤 2：推送镜像

```bash
# 推送单个标签
docker push your-name/new-api:latest

# 推送所有标签
docker push your-name/new-api --all-tags
```

### 步骤 3：验证

访问：`https://hub.docker.com/r/your-name/new-api`

---

## 📝 镜像命名最佳实践

### 1. 语义化版本号（推荐）

```bash
# 主版本.次版本.修订号
v1.0.0    # 初始版本
v1.0.1    # Bug 修复
v1.1.0    # 新增功能
v2.0.0    # 重大更新
```

### 2. 功能标识

```bash
# 格式：功能名-版本号
conversation-v1.0.0          # 对话记录功能 v1.0.0
conversation-archive-v1.1.0  # 增加归档功能 v1.1.0
```

### 3. 基于原版本

```bash
# 格式：原版本-自定义版本
upstream-v0.6.5-custom-v1.0.0  # 基于原项目 v0.6.5，自定义 v1.0.0
```

### 4. 环境标识

```bash
v1.0.0-prod    # 生产环境
v1.0.0-dev     # 开发环境
v1.0.0-test    # 测试环境
```

---

## 🎯 完整的镜像标签策略

### 推荐的标签组合

```bash
# 1. 构建镜像（打多个标签）
docker build \
  -t zhang/new-api:latest \                    # 最新版
  -t zhang/new-api:v1 \                        # 主版本
  -t zhang/new-api:v1.0 \                      # 次版本
  -t zhang/new-api:v1.0.0 \                    # 完整版本
  -t zhang/new-api:conversation \              # 功能标识
  -t zhang/new-api:v1.0.0-conversation \       # 版本+功能
  .

# 2. 推送到 Docker Hub
docker push zhang/new-api --all-tags
```

### 使用场景

```yaml
# 开发环境：使用 latest
image: zhang/new-api:latest

# 测试环境：使用功能标签
image: zhang/new-api:conversation

# 生产环境：使用固定版本号（最推荐）
image: zhang/new-api:v1.0.0-conversation
```

---

## 🔄 更新镜像版本

### 发布新版本流程

```bash
# 1. 修改代码后，构建新版本
docker build \
  -t zhang/new-api:latest \
  -t zhang/new-api:v1.1.0 \
  -t zhang/new-api:v1.1.0-conversation \
  .

# 2. 推送新版本
docker push zhang/new-api:v1.1.0
docker push zhang/new-api:v1.1.0-conversation
docker push zhang/new-api:latest  # 更新 latest 标签

# 3. 更新 docker-compose.yml
# 修改版本号为 v1.1.0

# 4. 重启服务
docker-compose pull
docker-compose up -d
```

---

## 📋 推荐的命名示例

### 场景 1：个人项目

```bash
# 如果你的 GitHub 用户名是 zhangsan
image: zhangsan/new-api:v1.0.0-conversation

# 或者使用真实姓名
image: zhang-wei/new-api:v1.0.0-conversation
```

### 场景 2：公司项目

```bash
# 如果公司叫 acme
image: acme/new-api:v1.0.0-conversation

# 或者更具体的项目名
image: acme/ai-gateway:v1.0.0
```

### 场景 3：团队项目

```bash
# 如果团队叫 dev-team
image: dev-team/new-api:v1.0.0-conversation
```

### 场景 4：私有仓库

```bash
# 使用私有 Docker Registry
image: registry.your-company.com/new-api:v1.0.0-conversation

# 或者使用其他云服务
image: your-account.azurecr.io/new-api:v1.0.0-conversation  # Azure
image: gcr.io/your-project/new-api:v1.0.0-conversation      # Google
```

---

## 🎨 镜像描述和文档

### 在 Docker Hub 上添加描述

创建 `README.docker.md` 文件：

```markdown
# New API - Conversation Logging Edition

基于 [calciumion/new-api](https://github.com/QuantumNous/new-api) 的二开版本。

## 新增功能

- ✅ 完整的对话内容记录
- ✅ 后台管理界面（按用户、模型、时间筛选）
- ✅ 批量删除和按条件删除
- ✅ 自动归档和数据压缩
- ✅ 性能优化（查询速度提升 5-10 倍）

## 快速开始

\`\`\`bash
docker run -d \\
  --name new-api \\
  -p 3000:3000 \\
  -e CONVERSATION_LOG_ENABLED=true \\
  -e CONVERSATION_ARCHIVE_DAYS=30 \\
  zhang/new-api:v1.0.0-conversation
\`\`\`

## 版本说明

- `latest` - 最新版本
- `v1.0.0` - 稳定版本
- `conversation` - 功能标识

## 文档

详细文档请访问：https://github.com/your-name/new-api

## 许可证

与原项目相同
```

---

## 🔍 镜像仓库对比

### Docker Hub（推荐，免费）

```bash
# 命名格式
your-username/new-api:v1.0.0

# 优点
- 免费（公开仓库无限制）
- 全球 CDN 加速
- 易于使用

# 缺点
- 私有仓库有数量限制（免费版只有1个）
```

### GitHub Container Registry（推荐，免费）

```bash
# 命名格式
ghcr.io/your-username/new-api:v1.0.0

# 优点
- 完全免费（公开和私有都无限制）
- 与 GitHub 集成
- 支持 GitHub Actions 自动构建

# 缺点
- 相对较新，生态不如 Docker Hub
```

### 阿里云容器镜像服务（国内推荐）

```bash
# 命名格式
registry.cn-hangzhou.aliyuncs.com/your-namespace/new-api:v1.0.0

# 优点
- 国内访问速度快
- 免费（个人版）
- 稳定性好

# 缺点
- 需要注册阿里云账号
```

---

## 🎉 最终推荐

### 如果你的名字/组织是 "Zhang"

```yaml
# docker-compose.yml

# 推荐方案 1：简洁明了
image: zhang/new-api:v1.0.0

# 推荐方案 2：功能标识（最推荐）✅
image: zhang/new-api:v1.0.0-conversation

# 推荐方案 3：完整描述
image: zhang/new-api-conversation:v1.0.0
```

### 构建和发布命令

```bash
# 1. 构建镜像
docker build \
  -t zhang/new-api:latest \
  -t zhang/new-api:v1.0.0-conversation \
  .

# 2. 测试镜像
docker run -d -p 3000:3000 zhang/new-api:v1.0.0-conversation

# 3. 推送到 Docker Hub
docker login
docker push zhang/new-api:v1.0.0-conversation
docker push zhang/new-api:latest

# 4. 更新 docker-compose.yml
# 修改 image: zhang/new-api:v1.0.0-conversation

# 5. 重新部署
docker-compose pull
docker-compose up -d
```

---

## ✅ 总结

**我的建议**：

```yaml
# 格式：<你的名字>/new-api:v<版本号>-conversation
image: zhang/new-api:v1.0.0-conversation
```

**原因**：
- ✅ 清晰表明这是你的定制版本
- ✅ 版本号便于管理和回滚
- ✅ `conversation` 标识表明包含对话记录功能
- ✅ 符合 Docker 命名最佳实践
- ✅ 方便与原版本区分

**替代选择**（根据你的实际情况）：
- 个人项目：`your-name/new-api:v1.0.0-conversation`
- 公司项目：`company-name/new-api:v1.0.0-conversation`
- 简化版本：`your-name/new-api:conversation`（不推荐，缺少版本号）
