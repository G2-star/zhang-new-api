# 🚀 New API 对话记录功能 - 构建和部署指南

## ✅ 所有修复已完成

以下问题已在代码中修复：
1. ✅ OptimizeConversationTable 函数调用错误
2. ✅ Semi UI message 导入错误
3. ✅ renderTimestamp 导入错误
4. ✅ ConversationArchive 表迁移
5. ✅ CONVERSATION_LOG_ENABLED 环境变量支持
6. ✅ **关键修复**: migrateLOGDB 在 LOG_SQL_DSN 为空时也会执行

---

## 📦 构建最终版本镜像

### Windows PowerShell
```powershell
cd D:\Users\Zhang\Desktop\new-api

# 构建镜像（包含所有修复）
docker build `
  -t agaid1mnjh45/new-api:latest `
  -t agaid1mnjh45/new-api:v1.0.0-conversation `
  .
```

### Linux/Mac/Git Bash
```bash
cd D:\Users\Zhang\Desktop\new-api

# 构建镜像（包含所有修复）
docker build \
  -t agaid1mnjh45/new-api:latest \
  -t agaid1mnjh45/new-api:v1.0.0-conversation \
  .
```

**预计时间**: 10-20 分钟（首次构建）

---

## 🔍 验证镜像

```bash
# 查看镜像
docker images agaid1mnjh45/new-api

# 应该看到两个标签
agaid1mnjh45/new-api   latest                   <image-id>   Just now   xxx MB
agaid1mnjh45/new-api   v1.0.0-conversation      <image-id>   Just now   xxx MB
```

---

## 🧪 测试新镜像

```bash
# 停止旧容器
docker stop <old-container>

# 启动新镜像测试
docker run -d --name new-api-test \
  -p 3001:3000 \
  -e CONVERSATION_LOG_ENABLED=true \
  -e SQL_DSN=postgresql://root:123456@your-postgres-host:5432/new-api \
  agaid1mnjh45/new-api:v1.0.0-conversation

# 查看启动日志（重点检查迁移日志）
docker logs -f new-api-test

# 应该看到：
# - "log database migration started (using main database)"
# - "database migration started"
# - "database migrated"

# 测试对话记录 API
curl http://localhost:3001/api/conversation/setting

# 清理测试容器
docker stop new-api-test && docker rm new-api-test
```

---

## 📤 推送到 Docker Hub

```bash
# 登录 Docker Hub
docker login
# 输入用户名: agaid1mnjh45
# 输入密码: (您的 Docker Hub 密码)

# 推送所有标签
docker push agaid1mnjh45/new-api:latest
docker push agaid1mnjh45/new-api:v1.0.0-conversation
```

**验证推送成功**:
访问 https://hub.docker.com/r/agaid1mnjh45/new-api

---

## 🎯 客户部署指南

### 方式 1: Docker Run
```bash
docker run -d --name new-api \
  --restart always \
  -p 3000:3000 \
  -e CONVERSATION_LOG_ENABLED=true \
  -e SQL_DSN=postgresql://user:password@host:5432/dbname \
  agaid1mnjh45/new-api:v1.0.0-conversation
```

### 方式 2: Docker Compose
```yaml
version: '3'
services:
  new-api:
    image: agaid1mnjh45/new-api:v1.0.0-conversation
    container_name: new-api
    restart: always
    ports:
      - "3000:3000"
    environment:
      # 数据库配置
      - SQL_DSN=postgresql://user:password@postgres:5432/newapi

      # 对话记录功能（可选，默认关闭）
      - CONVERSATION_LOG_ENABLED=true
      - CONVERSATION_ARCHIVE_DAYS=30
      - CONVERSATION_CLEANUP_DAYS=365

      # 会话密钥（必须修改）
      - SESSION_SECRET=your-random-string-here
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    container_name: postgres
    restart: always
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=newapi
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## ✅ 自动创建表的保证

使用最新镜像部署时，**表会自动创建**，因为：

1. ✅ `migrateLOGDB()` 现在无论 `LOG_SQL_DSN` 是否设置都会执行
2. ✅ `Conversation` 和 `ConversationArchive` 表已注册到迁移
3. ✅ 启动时会自动运行迁移，创建所有必需的表

**客户只需要**：
- 拉取镜像
- 设置环境变量 `CONVERSATION_LOG_ENABLED=true`
- 启动容器

**不需要**：
- ❌ 手动连接数据库
- ❌ 手动执行 SQL
- ❌ 任何额外的安装步骤

---

## 🔍 故障排查

如果客户报告表未创建，检查：

1. **启动日志**
```bash
docker logs <container-name> | grep -i "migration\|database"
```

应该看到：
```
log database migration started (using main database)
database migration started
database migrated
```

2. **环境变量**
```bash
docker inspect <container-name> | grep CONVERSATION_LOG_ENABLED
```

3. **节点类型**
确保不是 slave 节点：
```bash
docker inspect <container-name> | grep NODE_TYPE
```

如果是 `NODE_TYPE=slave`，迁移会被跳过。

---

## 📝 发布说明

向客户说明：

**v1.0.0-conversation 版本特性**
- ✅ 完整的对话记录功能
- ✅ 管理后台界面（管理 → 对话记录）
- ✅ 自动数据归档和清理
- ✅ 一键启用/禁用
- ✅ 数据库表自动创建（无需手动操作）

**启用方法**
只需在 docker-compose.yml 或启动命令中添加：
```yaml
environment:
  - CONVERSATION_LOG_ENABLED=true
```

**访问地址**
```
http://your-domain:3000/console/conversation
（需要管理员权限）
```

---

## 🎉 总结

1. **立即构建最新镜像** - 包含所有修复
2. **推送到 Docker Hub** - 客户可直接拉取
3. **提供部署文档** - docker-compose 示例
4. **强调自动化** - 无需手动 SQL 操作

这样客户部署时就能**开箱即用**，表会自动创建！
