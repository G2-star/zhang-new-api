# 📦 Docker Hub 发布完整指南

## 🎯 目标

将你的 New API 对话记录功能版本打包并上传到 Docker Hub，让其他人可以使用。

---

## 📋 前置准备

### 1. 注册 Docker Hub 账号

访问 https://hub.docker.com/ 注册账号（如果还没有）

- 记住你的用户名（例如：`zhang`）
- 记住密码

### 2. 确认项目状态

```bash
# 进入项目目录
cd D:\Users\Zhang\Desktop\new-api

# 确认 Dockerfile 存在
ls -lh Dockerfile

# 确认代码已提交（可选，但推荐）
git status
```

---

## 🚀 完整上传流程

### 步骤 1: 登录 Docker Hub

#### Windows PowerShell

```powershell
# 登录 Docker Hub
docker login

# 输入你的 Docker Hub 用户名
# 输入你的 Docker Hub 密码
```

#### Windows CMD

```cmd
docker login
```

**注意事项：**
- 用户名和密码输入时不会显示（正常现象）
- 看到 `Login Succeeded` 表示登录成功

---

### 步骤 2: 构建 Docker 镜像

#### 方案 A: 单标签构建（快速）

```bash
# 替换 zhang 为你的 Docker Hub 用户名
docker build -t zhang/new-api:v1.0.0-conversation .
```

#### 方案 B: 多标签构建（推荐）⭐

```bash
# Windows PowerShell（推荐）
docker build `
  -t zhang/new-api:latest `
  -t zhang/new-api:v1.0.0-conversation `
  -t zhang/new-api:conversation `
  -t zhang/new-api:v1.0.0 `
  .

# Windows CMD
docker build ^
  -t zhang/new-api:latest ^
  -t zhang/new-api:v1.0.0-conversation ^
  -t zhang/new-api:conversation ^
  -t zhang/new-api:v1.0.0 ^
  .

# Linux/Mac
docker build \
  -t zhang/new-api:latest \
  -t zhang/new-api:v1.0.0-conversation \
  -t zhang/new-api:conversation \
  -t zhang/new-api:v1.0.0 \
  .
```

**标签说明：**
- `latest` - 最新版本（默认）
- `v1.0.0-conversation` - 完整版本标识（推荐生产环境使用）
- `conversation` - 功能标识
- `v1.0.0` - 版本号

**构建时间：**
- 首次构建：10-20 分钟（取决于网络速度）
- 后续构建：5-10 分钟（会使用缓存）

**可能遇到的问题：**

1. **网络超时**
   ```bash
   # 使用国内镜像加速
   # 编辑 Docker Desktop 设置 → Docker Engine
   # 添加以下内容：
   {
     "registry-mirrors": [
       "https://docker.mirrors.ustc.edu.cn",
       "https://registry.docker-cn.com"
     ]
   }
   ```

2. **磁盘空间不足**
   ```bash
   # 清理未使用的镜像
   docker system prune -a
   ```

---

### 步骤 3: 验证镜像构建成功

```bash
# 查看本地镜像
docker images zhang/new-api

# 应该看到类似输出：
# REPOSITORY       TAG                      IMAGE ID       CREATED         SIZE
# zhang/new-api    latest                   abc123def456   2 minutes ago   50MB
# zhang/new-api    v1.0.0-conversation     abc123def456   2 minutes ago   50MB
# zhang/new-api    conversation             abc123def456   2 minutes ago   50MB
# zhang/new-api    v1.0.0                   abc123def456   2 minutes ago   50MB
```

---

### 步骤 4: 本地测试镜像（重要！）

```bash
# 使用构建的镜像启动容器
docker run -d \
  --name new-api-test \
  -p 3001:3000 \
  -e CONVERSATION_LOG_ENABLED=true \
  zhang/new-api:v1.0.0-conversation

# 等待 10 秒让服务启动
# 访问 http://localhost:3001 测试

# 测试成功后停止并删除容器
docker stop new-api-test
docker rm new-api-test
```

**Windows PowerShell 版本：**
```powershell
docker run -d `
  --name new-api-test `
  -p 3001:3000 `
  -e CONVERSATION_LOG_ENABLED=true `
  zhang/new-api:v1.0.0-conversation
```

---

### 步骤 5: 推送镜像到 Docker Hub

#### 方案 A: 推送所有标签（推荐）

```bash
# 推送所有标签
docker push zhang/new-api --all-tags
```

**推送时间：**
- 首次推送：5-15 分钟（取决于镜像大小和网络速度）
- 后续推送：1-5 分钟（只推送变更的层）

#### 方案 B: 推送单个标签

```bash
# 只推送特定标签
docker push zhang/new-api:v1.0.0-conversation

# 或推送 latest
docker push zhang/new-api:latest
```

**推送进度显示：**
```
The push refers to repository [docker.io/zhang/new-api]
abc123: Pushed
def456: Pushed
ghi789: Pushed
v1.0.0-conversation: digest: sha256:xxx size: 1234
```

看到 `digest: sha256:xxx` 表示推送成功！

---

### 步骤 6: 验证上传成功

#### 方法 1: 浏览器访问

访问：`https://hub.docker.com/r/zhang/new-api`

（替换 `zhang` 为你的用户名）

#### 方法 2: 命令行验证

```bash
# 删除本地镜像
docker rmi zhang/new-api:v1.0.0-conversation

# 从 Docker Hub 拉取
docker pull zhang/new-api:v1.0.0-conversation

# 如果能成功拉取，说明上传成功
```

---

## 📝 添加镜像描述（推荐）

### 在 Docker Hub 网站上添加说明

1. 访问 https://hub.docker.com/r/zhang/new-api
2. 点击 "Edit Repository"
3. 在 "Description" 中添加简短描述
4. 在 "Full Description" 中添加详细说明（支持 Markdown）

### 推荐的描述内容

**简短描述（Short Description）：**
```
New API with conversation logging - 基于 new-api 的二开版本，增加完整对话记录功能
```

**详细描述（Full Description）：**
```markdown
# New API - Conversation Logging Edition

基于 [QuantumNous/new-api](https://github.com/QuantumNous/new-api) 的二开版本。

## ✨ 新增功能

- ✅ 完整的对话内容记录（请求 + 响应）
- ✅ 后台管理界面（按用户、模型、时间筛选）
- ✅ 批量删除和按条件删除
- ✅ 自动归档和数据压缩
- ✅ 性能优化（查询速度提升 5-10 倍，存储节省 70-90%）
- ✅ 自动维护任务（每日归档、每周清理、每月优化）

## 🚀 快速开始

### 使用 Docker Run

\`\`\`bash
docker run -d \
  --name new-api \
  -p 3000:3000 \
  -e CONVERSATION_LOG_ENABLED=true \
  -e CONVERSATION_ARCHIVE_DAYS=30 \
  -v ./data:/data \
  zhang/new-api:v1.0.0-conversation
\`\`\`

### 使用 Docker Compose

\`\`\`yaml
version: '3.8'

services:
  new-api:
    image: zhang/new-api:v1.0.0-conversation
    container_name: new-api
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - CONVERSATION_LOG_ENABLED=true
      - CONVERSATION_ARCHIVE_DAYS=30
      - CONVERSATION_CLEANUP_DAYS=365
      - SQL_DSN=postgresql://user:pass@postgres:5432/new-api
    volumes:
      - ./data:/data
      - ./logs:/app/logs
\`\`\`

## 📖 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `CONVERSATION_LOG_ENABLED` | `false` | 启用对话记录功能 |
| `CONVERSATION_ARCHIVE_DAYS` | `30` | 归档天数 |
| `CONVERSATION_CLEANUP_DAYS` | `365` | 清理归档天数 |
| `SQL_DSN` | - | 数据库连接字符串 |

## 📊 版本标签

- `latest` - 最新版本
- `v1.0.0-conversation` - 稳定版本（推荐生产环境）
- `conversation` - 功能标识版本
- `v1.0.0` - 版本号

## 📚 文档

详细文档请访问：https://github.com/zhang/new-api

## 🔧 API 接口

新增 12 个对话管理 API 接口：

- `GET /api/conversation/` - 获取对话列表
- `GET /api/conversation/:id` - 获取对话详情
- `DELETE /api/conversation/` - 批量删除
- `POST /api/conversation/delete_by_condition` - 按条件删除
- `GET /api/conversation/stats` - 获取统计信息
- `POST /api/conversation/archive` - 手动归档
- 更多...

## ⚙️ 性能优化

- 查询速度提升：5-10 倍
- 存储空间节省：70-90%
- 支持百万级对话记录
- 自动归档和压缩

## 📄 许可证

与原项目相同

## 🙏 致谢

基于 [QuantumNous/new-api](https://github.com/QuantumNous/new-api) 项目
\`\`\`

---

## 🔄 更新镜像版本

### 发布新版本（v1.1.0）

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

# 3. 在 Docker Hub 上添加 Release Notes（可选）
```

---

## 🎨 自动化构建（可选，高级）

### 使用 GitHub Actions 自动构建

创建 `.github/workflows/docker-build.yml`：

```yaml
name: Build and Push Docker Image

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Log in to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Extract version
        id: version
        run: echo "VERSION=${GITHUB_REF#refs/tags/}" >> $GITHUB_OUTPUT

      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: |
            zhang/new-api:latest
            zhang/new-api:${{ steps.version.outputs.VERSION }}
            zhang/new-api:${{ steps.version.outputs.VERSION }}-conversation
            zhang/new-api:conversation
          cache-from: type=registry,ref=zhang/new-api:latest
          cache-to: type=inline
```

**使用步骤：**

1. 在 GitHub 仓库的 Settings → Secrets 添加：
   - `DOCKER_USERNAME`: 你的 Docker Hub 用户名
   - `DOCKER_PASSWORD`: 你的 Docker Hub 密码或 Access Token

2. 推送 tag 时自动构建：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

---

## 🛡️ 安全建议

### 使用 Access Token 代替密码

1. 访问 Docker Hub → Account Settings → Security
2. 点击 "New Access Token"
3. 创建 token（权限选择 Read & Write）
4. 使用 token 登录：
   ```bash
   docker login -u zhang -p <your-token>
   ```

### 扫描镜像漏洞

```bash
# 使用 Docker Scout 扫描
docker scout cves zhang/new-api:v1.0.0-conversation

# 或使用 Trivy
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy image zhang/new-api:v1.0.0-conversation
```

---

## 📊 镜像大小优化（可选）

当前 Dockerfile 已经优化良好，但如果需要进一步优化：

### 查看镜像大小

```bash
docker images zhang/new-api
```

### 查看镜像层

```bash
docker history zhang/new-api:v1.0.0-conversation
```

### 优化建议

1. **多阶段构建** - ✅ 已使用
2. **Alpine 基础镜像** - ✅ 已使用
3. **清理缓存** - 可以添加

---

## 🆘 常见问题

### Q1: 构建失败 - "network timeout"

**解决方案：**
```bash
# 使用镜像加速
# 在 Docker Desktop → Settings → Docker Engine 添加：
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn"
  ]
}
```

### Q2: 推送失败 - "unauthorized"

**解决方案：**
```bash
# 重新登录
docker logout
docker login
```

### Q3: 推送失败 - "denied: requested access to the resource is denied"

**原因：** 镜像名称格式错误或没有权限

**解决方案：**
```bash
# 确保格式正确：<用户名>/<仓库名>:<标签>
docker tag zhang/new-api:v1.0.0-conversation zhang/new-api:latest
docker push zhang/new-api:latest
```

### Q4: 如何删除 Docker Hub 上的镜像？

1. 访问 https://hub.docker.com/r/zhang/new-api/tags
2. 点击要删除的标签后面的删除按钮
3. 确认删除

### Q5: 如何设置镜像为私有？

1. 访问 https://hub.docker.com/r/zhang/new-api/settings
2. 将 "Visibility" 设置为 "Private"
3. 保存

**注意：** 免费账户只能有 1 个私有仓库

---

## ✅ 完整检查清单

发布前检查：

- [ ] Dockerfile 存在且正确
- [ ] 代码已提交到 Git（推荐）
- [ ] Docker Desktop 正在运行
- [ ] 已登录 Docker Hub
- [ ] 镜像构建成功
- [ ] 本地测试通过
- [ ] 镜像标签正确

发布后检查：

- [ ] 推送成功
- [ ] Docker Hub 上可以看到镜像
- [ ] 能从 Docker Hub 拉取镜像
- [ ] 添加了镜像描述
- [ ] 添加了使用文档

---

## 🎯 快速命令总结

```bash
# 1. 登录
docker login

# 2. 构建（选择一个）
# Windows PowerShell
docker build -t zhang/new-api:v1.0.0-conversation .

# 或多标签
docker build `
  -t zhang/new-api:latest `
  -t zhang/new-api:v1.0.0-conversation `
  .

# 3. 测试
docker run -d --name test -p 3001:3000 zhang/new-api:v1.0.0-conversation
# 访问 http://localhost:3001
docker stop test && docker rm test

# 4. 推送
docker push zhang/new-api --all-tags

# 5. 验证
# 访问 https://hub.docker.com/r/zhang/new-api
```

---

## 🎉 完成！

你的镜像现在已经在 Docker Hub 上了！

**分享给其他人：**
```bash
docker pull zhang/new-api:v1.0.0-conversation
```

**Docker Hub 链接：**
https://hub.docker.com/r/zhang/new-api

**祝发布顺利！** 🚀

---

*提示：记得将本文档中所有的 `zhang` 替换为你的实际 Docker Hub 用户名*
