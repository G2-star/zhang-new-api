# 对话记录功能 - 完整部署指南

## 📖 功能说明

本二开功能为 New API 项目添加了完整的 AI 对话内容记录和管理功能，包括：

- ✅ 记录所有模型的所有用户的完整对话（请求+响应）
- ✅ 后台管理界面支持按用户、模型、时间筛选
- ✅ 批量多选删除和按条件批量删除
- ✅ 功能开关，可随时启用/禁用
- ✅ 兼容原项目，可与上游版本合并

---

## 🏗️ 架构设计

### 1. 数据库设计

**新增表：`conversations`**

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | int | 主键 | ✓ |
| user_id | int | 用户ID | ✓ (复合) |
| username | string | 用户名 | ✓ |
| model_name | string | 模型名称 | ✓ (复合) |
| token_id | int | Token ID | ✓ |
| token_name | string | Token名称 | - |
| channel_id | int | 渠道ID | ✓ |
| request_messages | text | 请求消息(JSON) | - |
| response_content | text | 响应内容 | - |
| prompt_tokens | int | 输入token数 | - |
| completion_tokens | int | 输出token数 | - |
| total_tokens | int | 总token数 | - |
| is_stream | bool | 是否流式 | - |
| created_at | int64 | 创建时间戳 | ✓ (复合) |
| use_time | int | 响应时间(ms) | - |
| ip | string | 客户端IP | ✓ |
| group | string | 用户组 | ✓ |

**索引设计：**
- 复合索引：`(user_id, model_name, created_at)` - 优化按用户和模型筛选
- 单字段索引：`user_id`, `model_name`, `created_at`, `channel_id`, `ip`, `group`

### 2. 代码结构

```
new-api/
├── model/
│   └── conversation.go          # 数据库模型和操作函数
├── controller/
│   └── conversation.go          # HTTP API 控制器
├── router/
│   └── api-router.go            # 路由配置（已修改）
├── relay/
│   ├── conversation_helper.go   # 对话记录辅助函数
│   └── conversation_middleware.go # 对话记录中间件
├── common/
│   └── constants.go             # 全局配置（已修改）
└── web/src/pages/
    └── Conversation/
        └── index.jsx            # 前端管理页面
```

---

## 🚀 部署步骤

### 方法一：直接使用已修改的代码

1. **复制文件到项目**
   ```bash
   # 已创建的新文件（直接可用）：
   # - model/conversation.go
   # - controller/conversation.go
   # - relay/conversation_helper.go
   # - relay/conversation_middleware.go
   # - web/src/pages/Conversation/index.jsx

   # 已修改的文件：
   # - common/constants.go (添加了 ConversationLogEnabled 变量)
   # - model/main.go (添加了 Conversation 表迁移)
   # - router/api-router.go (添加了对话记录路由)
   ```

2. **数据库迁移**
   ```bash
   # 启动应用时会自动创建 conversations 表
   # 无需手动执行迁移脚本
   ```

3. **重新编译**
   ```bash
   cd new-api
   go mod tidy
   go build -o new-api
   ```

4. **前端编译**
   ```bash
   cd web
   npm install
   npm run build
   cd ..
   ```

5. **重启服务**
   ```bash
   # 方式1：直接运行
   ./new-api

   # 方式2：Docker
   docker-compose down
   docker-compose build
   docker-compose up -d
   ```

### 方法二：手动集成（如果需要自定义）

参考 `CONVERSATION_INTEGRATION_GUIDE.md` 文件中的详细说明。

---

## ⚙️ 配置说明

### 环境变量（可选）

在 `.env` 文件中添加：

```bash
# 对话记录功能开关（默认关闭，建议通过后台界面控制）
CONVERSATION_LOG_ENABLED=false

# 对话内容最大长度（字符数，超过会被截断，默认 50000）
CONVERSATION_MAX_LENGTH=50000
```

### 功能开关

**方式1：通过后台界面**
1. 登录管理员账号
2. 进入"对话记录管理"页面
3. 点击"修改设置"按钮
4. 切换启用/禁用状态

**方式2：通过 API**
```bash
# 查询状态
curl -X GET http://localhost:3000/api/conversation/setting \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# 启用功能
curl -X PUT http://localhost:3000/api/conversation/setting \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'

# 禁用功能
curl -X PUT http://localhost:3000/api/conversation/setting \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}'
```

---

## 🎯 功能使用

### 1. 后台管理界面

访问路径：`http://your-domain:3000/conversation`（需要管理员权限）

**功能列表：**
- 查看所有对话记录
- 按用户名筛选
- 按模型名称筛选
- 按时间范围筛选
- 查看对话详情（完整的请求和响应内容）
- 批量选择删除
- 按条件批量删除

### 2. API 接口

#### 获取对话记录列表
```http
GET /api/conversation/
Query参数:
  - page: 页码（默认1）
  - page_size: 每页数量（默认10，最大100）
  - user_id: 用户ID（可选）
  - username: 用户名（可选）
  - model_name: 模型名称（可选，支持模糊匹配）
  - start_time: 开始时间戳（可选）
  - end_time: 结束时间戳（可选）
```

#### 获取对话详情
```http
GET /api/conversation/:id
```

#### 批量删除对话
```http
DELETE /api/conversation/
Body: {
  "ids": [1, 2, 3, ...]
}
```

#### 按条件批量删除
```http
POST /api/conversation/delete_by_condition
Body: {
  "user_id": 123,          // 可选
  "username": "user1",     // 可选
  "model_name": "gpt-4",   // 可选
  "start_time": 1234567890, // 可选
  "end_time": 1234567999    // 可选
}
注意：至少需要一个筛选条件
```

#### 获取统计信息
```http
GET /api/conversation/stats
Query参数同上
```

---

## 🔐 安全和隐私

### 1. 权限控制

- 所有对话记录 API 均需要**管理员权限**
- 普通用户无法查看对话记录
- 遵循用户设置中的 IP 记录选项

### 2. 数据保护

- 对话内容存储在独立的 `conversations` 表中
- 支持单独的日志数据库配置（`LOG_SQL_DSN`）
- 建议定期清理历史数据

### 3. 合规建议

⚠️ **重要提示：**

1. **法律合规**：请确保启用对话记录功能符合您所在地区的法律法规（如 GDPR、CCPA 等）
2. **用户告知**：建议在用户协议和隐私政策中明确说明会记录对话内容
3. **数据安全**：建议对数据库进行加密存储
4. **定期清理**：建议定期删除旧的对话记录

---

## 🧪 测试

### 1. 功能测试

```bash
# 测试1：启用功能
# 1. 通过后台界面启用对话记录功能
# 2. 发送一个 AI 请求
# 3. 检查 conversations 表中是否有新记录

# 测试2：查询功能
# 1. 访问对话记录管理页面
# 2. 测试各种筛选条件
# 3. 查看对话详情

# 测试3：删除功能
# 1. 选择一条记录删除
# 2. 测试批量删除
# 3. 测试按条件删除
```

### 2. 性能测试

```bash
# 测试并发请求对性能的影响
ab -n 1000 -c 10 http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -p request.json

# 监控数据库大小
# conversations 表大小会随着记录增加而增长
# 建议定期清理超过 30 天的记录
```

---

## 🔧 维护

### 1. 定期清理旧数据

创建定时任务清理旧记录：

```go
// 在 main.go 中添加定时清理任务
go func() {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        // 删除 30 天前的记录
        targetTime := time.Now().AddDate(0, 0, -30).Unix()
        deleted, err := model.DeleteOldConversations(context.Background(), targetTime, 1000)
        if err != nil {
            common.SysLog("清理旧对话记录失败: " + err.Error())
        } else {
            common.SysLog(fmt.Sprintf("成功清理 %d 条旧对话记录", deleted))
        }
    }
}()
```

### 2. 数据库优化

```sql
-- MySQL 优化
OPTIMIZE TABLE conversations;

-- 检查索引使用情况
EXPLAIN SELECT * FROM conversations WHERE user_id = 1 AND created_at > 1234567890;

-- PostgreSQL 优化
VACUUM ANALYZE conversations;
```

### 3. 监控指标

关注以下指标：
- conversations 表大小
- 查询响应时间
- 磁盘使用率
- 每日新增记录数

---

## 🤝 与上游版本合并

本二开功能的设计确保了与原项目的兼容性：

### 1. 独立性

- **新增文件**：不修改现有文件的核心逻辑
- **独立表**：使用单独的 `conversations` 表
- **可选功能**：默认关闭，不影响现有功能

### 2. 合并策略

当需要合并上游更新时：

```bash
# 1. 备份修改的文件
git stash

# 2. 拉取上游更新
git pull upstream main

# 3. 恢复修改
git stash pop

# 4. 解决冲突（如果有）
# 主要可能冲突的文件：
# - common/constants.go (ConversationLogEnabled 变量)
# - model/main.go (Conversation 表迁移)
# - router/api-router.go (conversation 路由)
```

### 3. 冲突处理

如果上游版本修改了以下文件，需要手动合并：

- `common/constants.go`: 保留 `ConversationLogEnabled` 变量
- `model/main.go`: 保留 `Conversation` 表的 AutoMigrate 调用
- `router/api-router.go`: 保留 conversation 路由组

---

## 📝 更新日志

### v1.0.0 (2025-01-XX)

初始版本，包含以下功能：
- ✅ 完整的对话记录数据库模型
- ✅ 后台管理 API
- ✅ 前端管理界面
- ✅ 功能开关和隐私保护
- ✅ 批量删除和条件筛选
- ✅ 异步记录，不影响主流程性能

---

## 🆘 常见问题

### Q1: 对话记录功能会影响性能吗？

A: 不会。对话记录采用异步方式（goroutine），不阻塞主请求流程。即使记录失败也不会影响 API 响应。

### Q2: 数据库会很快变大吗？

A: 是的。对话内容通常较长，建议：
- 定期清理旧数据（如保留 30 天）
- 使用单独的日志数据库
- 考虑使用 PostgreSQL 或 MySQL（比 SQLite 性能更好）

### Q3: 可以只记录特定用户或模型的对话吗？

A: 当前版本记录所有对话。如需定制化，可修改 `relay/conversation_helper.go` 中的 `RecordConversationHelper` 函数添加过滤条件。

### Q4: 如何备份对话数据？

A:
```bash
# MySQL
mysqldump -u root -p database_name conversations > conversations_backup.sql

# PostgreSQL
pg_dump -U postgres -t conversations database_name > conversations_backup.sql

# SQLite
sqlite3 new_api.db ".dump conversations" > conversations_backup.sql
```

### Q5: 前端页面在哪里添加到菜单？

A: 需要在 `web/src/App.jsx` 或相应的路由配置文件中添加菜单项和路由。

---

## 📞 支持

如有问题，请：
1. 检查日志文件中的错误信息
2. 确认数据库连接正常
3. 验证管理员权限
4. 查看浏览器控制台错误

---

## 📄 许可证

本二开功能遵循 New API 项目的原始许可证。
