# 对话记录功能 - 部署检查清单

## ✅ 文件清单

### 已创建的新文件（可直接使用）

- [x] `model/conversation.go` - 对话记录数据库模型
- [x] `controller/conversation.go` - 对话记录API控制器
- [x] `relay/conversation_helper.go` - 对话记录辅助函数
- [x] `relay/conversation_middleware.go` - 对话记录中间件
- [x] `web/src/pages/Conversation/index.jsx` - 前端管理页面
- [x] `CONVERSATION_FEATURE_README.md` - 完整部署文档
- [x] `CONVERSATION_INTEGRATION_GUIDE.md` - 集成指南

### 已修改的文件

- [x] `common/constants.go` - 添加了 `ConversationLogEnabled` 变量
- [x] `model/main.go` - 添加了 `Conversation` 表迁移
- [x] `router/api-router.go` - 添加了对话记录路由

---

## 📋 部署前检查

### 1. 代码完整性

- [ ] 确认所有新文件已复制到项目目录
- [ ] 确认所有修改已应用到对应文件
- [ ] 运行 `go mod tidy` 检查依赖

### 2. 编译测试

```bash
# 后端编译
cd /path/to/new-api
go build -o new-api

# 前端编译
cd web
npm install
npm run build
```

- [ ] 后端编译成功，无错误
- [ ] 前端编译成功，无错误

### 3. 数据库准备

- [ ] 确认数据库连接正常（MySQL/PostgreSQL/SQLite）
- [ ] 确认有足够的存储空间
- [ ] （可选）配置单独的日志数据库 `LOG_SQL_DSN`

---

## 🚀 部署步骤检查

### 步骤 1: 停止服务

```bash
# 方式1：直接运行的服务
pkill new-api

# 方式2：Docker
docker-compose down

# 方式3：Systemd
systemctl stop new-api
```

- [ ] 服务已停止

### 步骤 2: 备份数据

```bash
# 备份数据库
mysqldump -u root -p database_name > backup_$(date +%Y%m%d).sql

# 备份配置文件
cp .env .env.backup
```

- [ ] 数据库已备份
- [ ] 配置文件已备份

### 步骤 3: 更新代码

```bash
# 复制新文件
# （根据实际情况调整）

# 编译
go build -o new-api
```

- [ ] 代码已更新
- [ ] 编译成功

### 步骤 4: 数据库迁移

启动应用时会自动创建 `conversations` 表，无需手动操作。

- [ ] 启动应用后检查 `conversations` 表是否创建成功

### 步骤 5: 重启服务

```bash
# 方式1：直接运行
./new-api

# 方式2：Docker
docker-compose up -d

# 方式3：Systemd
systemctl start new-api
```

- [ ] 服务已启动
- [ ] 检查日志无错误

---

## 🧪 功能测试检查

### 测试 1: 基础功能

- [ ] 访问主页，确认服务正常
- [ ] 登录管理员账号
- [ ] 访问对话记录管理页面 `/conversation`（需要在前端路由中配置）

### 测试 2: 对话记录

**关闭状态测试：**
- [ ] 确认功能默认为关闭状态
- [ ] 发送一个 AI 请求
- [ ] 确认 `conversations` 表中无新记录

**开启状态测试：**
- [ ] 通过后台界面启用对话记录功能
- [ ] 发送一个非流式 AI 请求（如 GPT-4）
- [ ] 检查 `conversations` 表中是否有新记录
- [ ] 发送一个流式 AI 请求
- [ ] 检查 `conversations` 表中是否有新记录（需要集成 stream handler）

### 测试 3: 查询功能

- [ ] 访问对话记录列表
- [ ] 测试按用户名筛选
- [ ] 测试按模型名筛选
- [ ] 测试按时间范围筛选
- [ ] 测试分页功能
- [ ] 点击"查看详情"，确认能看到完整的请求和响应

### 测试 4: 删除功能

- [ ] 单条删除测试
- [ ] 批量选择删除测试（选择2-3条记录）
- [ ] 按条件批量删除测试（谨慎操作）

### 测试 5: API 测试

使用 curl 或 Postman 测试：

```bash
# 获取对话列表
curl -X GET "http://localhost:3000/api/conversation/?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# 获取对话详情
curl -X GET "http://localhost:3000/api/conversation/1" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# 获取功能状态
curl -X GET "http://localhost:3000/api/conversation/setting" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# 更新功能状态
curl -X PUT "http://localhost:3000/api/conversation/setting" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'

# 批量删除
curl -X DELETE "http://localhost:3000/api/conversation/" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"ids": [1, 2]}'
```

- [ ] 所有 API 测试通过

---

## 🔧 集成 Handler（关键步骤）

### 非流式 Handler 集成

需要修改 `relay/channel/openai/relay-openai.go` 中的 `OpenaiHandler` 函数：

```go
// 在函数开始处添加
startTime := time.Now()

// 在返回前添加（在 service.IOCopyBytesGracefully 之后）
if textReq, ok := info.Request.(*dto.GeneralOpenAIRequest); ok {
    responseContent := relay.ExtractResponseContent(&simpleResponse)
    relay.RecordConversationHelper(c, info, textReq, responseContent, &simpleResponse.Usage, startTime)
}
```

- [ ] 已集成到 `OpenaiHandler`
- [ ] 测试非流式请求记录成功

### 流式 Handler 集成

需要修改 `relay/channel/openai/relay-openai.go` 中的 `OpenaiStreamHandler` 函数：

```go
// 在函数开始处添加
startTime := time.Now()
var contentCollector *relay.StreamContentCollector
var textReq *dto.GeneralOpenAIRequest
if common.ConversationLogEnabled {
    contentCollector = &relay.StreamContentCollector{}
    if req, ok := info.Request.(*dto.GeneralOpenAIRequest); ok {
        textReq = req
    }
}

// 在流处理循环中添加（每次解析 streamResponse 后）
if contentCollector != nil && len(streamResponse.Choices) > 0 {
    contentCollector.AddChunk(&streamResponse.Choices[0].Delta)
}

// 在函数返回前添加
if contentCollector != nil && textReq != nil {
    responseContent := contentCollector.GetContent()
    relay.RecordConversationHelper(c, info, textReq, responseContent, usage, startTime)
}
```

- [ ] 已集成到 `OpenaiStreamHandler`
- [ ] 测试流式请求记录成功

### 其他模型 Handler（可选）

如需支持其他模型格式（Claude, Gemini 等），参考 `CONVERSATION_INTEGRATION_GUIDE.md` 进行集成。

- [ ] Claude Handler 已集成（如需要）
- [ ] Gemini Handler 已集成（如需要）
- [ ] 其他 Handler 已集成（如需要）

---

## 📊 性能测试

### 1. 基础性能

```bash
# 记录关闭时的基准测试
ab -n 1000 -c 10 http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -p request.json

# 记录开启后的性能测试
# 对比响应时间差异
```

- [ ] 记录关闭时的平均响应时间：______ ms
- [ ] 记录开启后的平均响应时间：______ ms
- [ ] 性能影响在可接受范围内（<5%）

### 2. 并发测试

```bash
# 100 并发，1000 请求
ab -n 1000 -c 100 http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -p request.json
```

- [ ] 无请求失败
- [ ] 对话记录数量正确
- [ ] 数据库无死锁或超时

### 3. 数据库测试

```sql
-- 查询性能测试
EXPLAIN ANALYZE SELECT * FROM conversations
WHERE user_id = 1 AND created_at > 1234567890
ORDER BY created_at DESC LIMIT 10;

-- 检查表大小
SELECT
    COUNT(*) as total_records,
    SUM(LENGTH(request_messages)) as total_request_size,
    SUM(LENGTH(response_content)) as total_response_size
FROM conversations;
```

- [ ] 查询使用了索引
- [ ] 查询响应时间 < 100ms
- [ ] 表大小在预期范围内

---

## 🔐 安全检查

### 权限控制

- [ ] 普通用户无法访问对话记录 API
- [ ] 只有管理员可以查看对话记录
- [ ] 只有管理员可以删除对话记录
- [ ] 只有管理员可以修改功能开关

### 数据保护

- [ ] 遵循用户 IP 记录设置
- [ ] 敏感数据已过滤（如 API Key）
- [ ] 数据库连接使用加密（生产环境）
- [ ] 定期清理计划已设置

### 合规性

- [ ] 已更新隐私政策说明对话记录功能
- [ ] 已告知用户对话可能被记录
- [ ] 符合当地法律法规（GDPR/CCPA 等）

---

## 📝 前端集成（可选）

如需在前端添加菜单项，需要修改前端路由配置：

```jsx
// 在 App.jsx 或路由配置文件中添加
import Conversation from './pages/Conversation';

// 添加路由
{
  path: '/conversation',
  element: <Conversation />,
  meta: { requiresAuth: true, requiresAdmin: true }
}

// 添加菜单项（在管理员菜单中）
{
  text: '对话记录',
  itemKey: 'conversation',
  icon: <IconComment />,
  path: '/conversation'
}
```

- [ ] 已添加路由配置
- [ ] 已添加菜单项
- [ ] 页面可正常访问

---

## 🎉 部署完成检查

### 最终验收

- [ ] 所有功能正常工作
- [ ] 性能影响在可接受范围内
- [ ] 安全检查全部通过
- [ ] 文档已完善
- [ ] 备份策略已制定
- [ ] 监控已设置

### 文档归档

- [ ] 部署日期：__________
- [ ] 部署人员：__________
- [ ] 版本号：__________
- [ ] 备注：__________

---

## 🆘 问题排查

### 问题 1: 对话记录表未创建

**症状**: 应用启动后 `conversations` 表不存在

**排查**:
```bash
# 检查日志
tail -f logs/new-api.log | grep conversation

# 手动检查数据库
mysql -u root -p -e "USE database_name; SHOW TABLES LIKE 'conversations';"
```

**解决**:
- 确认 `model/main.go` 中的 `migrateLOGDB` 函数包含 Conversation 迁移
- 检查数据库权限
- 尝试手动创建表（参考 `model/conversation.go` 中的结构）

### 问题 2: 对话未被记录

**症状**: 功能已启用，但数据库中无记录

**排查**:
```bash
# 检查配置
curl http://localhost:3000/api/conversation/setting \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# 检查日志
tail -f logs/new-api.log | grep "record conversation"
```

**解决**:
- 确认功能已启用（`ConversationLogEnabled = true`）
- 确认 Handler 已正确集成（参考集成指南）
- 检查是否有 panic 或错误日志

### 问题 3: 前端页面无法访问

**症状**: 访问 `/conversation` 返回 404

**排查**:
- 检查前端路由配置
- 检查文件是否正确放置在 `web/src/pages/Conversation/index.jsx`
- 检查是否重新编译前端

**解决**:
- 重新编译前端：`cd web && npm run build`
- 确认路由已添加到 App.jsx
- 清除浏览器缓存

### 问题 4: API 返回 403

**症状**: 调用对话记录 API 返回 403 Forbidden

**排查**:
- 确认使用的是管理员 Token
- 检查 Token 是否过期

**解决**:
- 使用管理员账号重新登录
- 生成新的 Access Token
- 确认账号有管理员权限

---

## 📞 联系支持

如遇到无法解决的问题，请提供以下信息：

1. 错误日志（最近 50 行）
2. 数据库版本和配置
3. 部署方式（Docker/直接运行/Systemd）
4. 已完成的检查清单项
5. 具体的错误现象和复现步骤

---

**检查清单完成时间**: __________

**签字**: __________
