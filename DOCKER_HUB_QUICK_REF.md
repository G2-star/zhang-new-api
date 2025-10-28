# 🚀 Docker Hub 发布快速参考卡

## ⚡ 5 分钟快速发布

### 步骤 1: 登录 Docker Hub
```bash
docker login
# 输入用户名和密码
```

### 步骤 2: 构建镜像

**Windows PowerShell:**
```powershell
cd D:\Users\Zhang\Desktop\new-api

docker build `
  -t zhang/new-api:latest `
  -t zhang/new-api:v1.0.0-conversation `
  .
```

**Windows CMD:**
```cmd
cd D:\Users\Zhang\Desktop\new-api

docker build ^
  -t zhang/new-api:latest ^
  -t zhang/new-api:v1.0.0-conversation ^
  .
```

### 步骤 3: 推送到 Docker Hub
```bash
docker push zhang/new-api --all-tags
```

### 步骤 4: 验证
访问: https://hub.docker.com/r/zhang/new-api

---

## 📋 完整命令（复制粘贴即可）

**Windows PowerShell 完整版:**
```powershell
# 1. 登录
docker login

# 2. 进入项目目录
cd D:\Users\Zhang\Desktop\new-api

# 3. 构建镜像
docker build `
  -t zhang/new-api:latest `
  -t zhang/new-api:v1.0.0-conversation `
  -t zhang/new-api:conversation `
  -t zhang/new-api:v1.0.0 `
  .

# 4. 本地测试
docker run -d --name test -p 3001:3000 `
  -e CONVERSATION_LOG_ENABLED=true `
  zhang/new-api:v1.0.0-conversation

# 等待 10 秒，然后访问 http://localhost:3001

# 5. 停止测试容器
docker stop test
docker rm test

# 6. 推送到 Docker Hub
docker push zhang/new-api --all-tags

# 完成！访问 https://hub.docker.com/r/zhang/new-api
```

---

## 🔧 常用命令

### 查看本地镜像
```bash
docker images zhang/new-api
```

### 测试镜像
```bash
docker run -d --name test -p 3001:3000 zhang/new-api:v1.0.0-conversation
docker logs -f test
docker stop test && docker rm test
```

### 删除本地镜像
```bash
docker rmi zhang/new-api:v1.0.0-conversation
```

### 从 Docker Hub 拉取
```bash
docker pull zhang/new-api:v1.0.0-conversation
```

### 查看镜像详情
```bash
docker inspect zhang/new-api:v1.0.0-conversation
```

### 查看镜像大小
```bash
docker images zhang/new-api --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
```

---

## ⏱️ 预计时间

| 步骤 | 时间 |
|------|------|
| 登录 Docker Hub | 30 秒 |
| 构建镜像（首次） | 10-20 分钟 |
| 构建镜像（后续） | 5-10 分钟 |
| 本地测试 | 2 分钟 |
| 推送到 Docker Hub（首次） | 5-15 分钟 |
| 推送到 Docker Hub（后续） | 1-5 分钟 |
| **总计（首次）** | **20-40 分钟** |
| **总计（后续）** | **10-20 分钟** |

---

## 🆘 快速故障排查

### 问题 1: "Cannot connect to Docker daemon"
```bash
# 确保 Docker Desktop 正在运行
# 重启 Docker Desktop
```

### 问题 2: "unauthorized: authentication required"
```bash
# 重新登录
docker logout
docker login
```

### 问题 3: "network timeout during build"
```bash
# 配置 Docker 镜像加速
# Docker Desktop → Settings → Docker Engine
# 添加国内镜像源
```

### 问题 4: "denied: requested access to the resource is denied"
```bash
# 检查镜像名称格式
# 必须是: <你的用户名>/<仓库名>:<标签>
# 例如: zhang/new-api:latest
```

### 问题 5: 推送很慢
```bash
# 检查网络连接
# 考虑使用 GitHub Container Registry (ghcr.io) 作为替代
```

---

## 📝 推荐的镜像标签策略

| 标签 | 用途 | 示例 |
|------|------|------|
| `latest` | 最新版本 | `zhang/new-api:latest` |
| `v1.0.0-conversation` | 完整版本标识 ⭐ | `zhang/new-api:v1.0.0-conversation` |
| `conversation` | 功能标识 | `zhang/new-api:conversation` |
| `v1.0.0` | 版本号 | `zhang/new-api:v1.0.0` |

**生产环境推荐：** `v1.0.0-conversation`

---

## 🌟 发布后的下一步

### 1. 添加镜像描述
访问: https://hub.docker.com/r/zhang/new-api/settings

### 2. 创建 README
在 Docker Hub 上添加使用说明

### 3. 设置自动构建（可选）
使用 GitHub Actions 自动构建和推送

### 4. 分享镜像
```bash
# 告诉其他人如何使用
docker pull zhang/new-api:v1.0.0-conversation
docker run -d -p 3000:3000 zhang/new-api:v1.0.0-conversation
```

---

## 📞 获取帮助

- **详细指南**: 查看 `DOCKER_HUB_PUBLISH_GUIDE.md`
- **Docker 文档**: https://docs.docker.com/
- **Docker Hub**: https://hub.docker.com/

---

**记得将所有 `zhang` 替换为你的实际 Docker Hub 用户名！** 🎯
