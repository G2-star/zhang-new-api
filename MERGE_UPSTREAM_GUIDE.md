# 与官方版本合并指南 - 保持二开功能的同时更新上游代码

## 🎯 目标

- ✅ 从官方仓库获取最新更新
- ✅ 保留你的二开功能（对话记录）
- ✅ 解决可能的冲突
- ✅ 互不影响，各自独立

---

## 📋 方案概述

### 策略：使用 Git 的多远程仓库功能

```
你的仓库 (origin)
    ↓ fork
官方仓库 (upstream) → 定期合并 → 你的仓库 (origin)
                                      ↓
                                你的二开功能始终保留
```

---

## 🚀 完整操作步骤

### 第一步：Push 你的项目到 GitHub

#### 1.1 创建 GitHub 仓库

访问 GitHub，创建新仓库：
- 仓库名：`new-api` 或 `new-api-conversation`
- 描述：`New API with conversation logging feature`
- 可见性：Public 或 Private

#### 1.2 初始化 Git 并 Push

```bash
# 进入项目目录
cd D:\Users\Zhang\Desktop\new-api

# 初始化 Git（如果还没有）
git init

# 添加你的 GitHub 仓库为 origin
git remote add origin https://github.com/zhang/new-api.git

# 或使用 SSH
git remote add origin git@github.com:zhang/new-api.git

# 创建 .gitignore 文件
cat > .gitignore << 'EOF'
# 数据文件
/data/
/logs/
*.db
*.sqlite

# 环境配置
.env
.env.local

# 编译产物
/bin/
/dist/
*.exe
*.dll
*.so
*.dylib

# 前端
/web/node_modules/
/web/dist/
/web/.vite/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# 系统文件
.DS_Store
Thumbs.db

# Docker
docker-compose.override.yml

# 日志
*.log

# 临时文件
*.tmp
*.temp
EOF

# 添加所有文件
git add .

# 第一次提交
git commit -m "feat: Add conversation logging feature

- Add complete conversation recording (request + response)
- Add backend management interface (filter by user/model/time)
- Add batch delete and conditional delete
- Add auto-archive and performance optimization
- Query performance improved 5-10x
- Storage space saved 70-90%

Version: v1.0.0-conversation
Based on: calciumion/new-api:latest
"

# 创建主分支（如果需要）
git branch -M main

# 推送到 GitHub
git push -u origin main
```

---

### 第二步：添加官方仓库为上游（Upstream）

```bash
# 添加官方仓库为 upstream
git remote add upstream https://github.com/QuantumNous/new-api.git

# 验证远程仓库
git remote -v

# 应该看到：
# origin    https://github.com/zhang/new-api.git (fetch)
# origin    https://github.com/zhang/new-api.git (push)
# upstream  https://github.com/QuantumNous/new-api.git (fetch)
# upstream  https://github.com/QuantumNous/new-api.git (push)
```

---

### 第三步：创建分支策略（重要！）

为了更好地管理合并，建议使用分支策略：

```bash
# 当前在 main 分支

# 创建开发分支
git checkout -b develop

# 创建功能分支（保存你的二开功能）
git checkout -b feature/conversation-logging

# 提交你的二开功能
git add .
git commit -m "feat: conversation logging feature"

# 推送功能分支
git push -u origin feature/conversation-logging

# 回到 main 分支
git checkout main
```

**分支说明：**
- `main` - 稳定版本，用于生产部署
- `develop` - 开发分支，用于日常开发
- `feature/conversation-logging` - 功能分支，保存对话记录功能
- `upstream-sync` - 同步分支，用于合并官方更新

---

### 第四步：首次同步官方代码

```bash
# 1. 获取官方最新代码
git fetch upstream

# 2. 查看官方的分支
git branch -r

# 3. 创建同步分支（基于官方 main 分支）
git checkout -b upstream-sync upstream/main

# 4. 查看官方最新的提交
git log --oneline -10

# 5. 回到你的 main 分支
git checkout main

# 6. 合并官方代码到你的 main 分支
git merge upstream/main

# 如果有冲突，继续看下一步
```

---

### 第五步：解决冲突（关键步骤）

#### 5.1 识别冲突文件

```bash
# 查看冲突文件
git status

# 可能的冲突文件：
# - common/constants.go
# - model/main.go
# - router/api-router.go
# - docker-compose.yml
# - go.mod
```

#### 5.2 解决冲突策略

**对于你修改过的文件，采用以下策略：**

##### 策略 1：保留你的修改 + 官方的新增（推荐）

```bash
# 对于 common/constants.go
git checkout --ours common/constants.go
# 然后手动添加官方的新增配置项

# 对于 model/main.go
git checkout --ours model/main.go
# 然后手动合并官方的新增迁移

# 对于 router/api-router.go
git checkout --ours router/api-router.go
# 然后手动添加官方的新增路由
```

##### 策略 2：使用合并工具

```bash
# 使用 VS Code 打开冲突文件
code common/constants.go

# VS Code 会显示冲突标记：
# <<<<<<< HEAD (你的修改)
# var ConversationLogEnabled = false
# =======
# var NewOfficialFeature = true (官方的修改)
# >>>>>>> upstream/main

# 选择：
# - Accept Current Change（保留你的）
# - Accept Incoming Change（使用官方的）
# - Accept Both Changes（都保留）
# - Compare Changes（对比后手动编辑）
```

##### 策略 3：手动合并（最安全）

```bash
# 1. 查看你的修改
git diff HEAD common/constants.go

# 2. 查看官方的修改
git diff upstream/main common/constants.go

# 3. 手动编辑文件，保留双方的修改
nano common/constants.go
# 或
code common/constants.go

# 4. 标记为已解决
git add common/constants.go
```

#### 5.3 具体文件的冲突解决示例

**文件 1: `common/constants.go`**

```go
// 冲突前（你的版本）
var LogConsumeEnabled = true
var ConversationLogEnabled = false // 你添加的

// 冲突（官方也在这里添加了新功能）
<<<<<<< HEAD
var LogConsumeEnabled = true
var ConversationLogEnabled = false // 对话记录功能开关
=======
var LogConsumeEnabled = true
var NewOfficialFeature = true // 官方新功能
>>>>>>> upstream/main

// 解决后（保留双方）
var LogConsumeEnabled = true
var ConversationLogEnabled = false // 对话记录功能开关（你的）
var NewOfficialFeature = true // 官方新功能（官方的）
```

**文件 2: `model/main.go`**

```go
// 你的版本
if err = LOG_DB.AutoMigrate(&Log{}); err != nil {
    return err
}
if err = LOG_DB.AutoMigrate(&Conversation{}); err != nil {
    return err
}

// 官方版本（假设官方也添加了新表）
if err = LOG_DB.AutoMigrate(&Log{}); err != nil {
    return err
}
if err = LOG_DB.AutoMigrate(&OfficialNewTable{}); err != nil {
    return err
}

// 解决后（保留双方）
if err = LOG_DB.AutoMigrate(&Log{}); err != nil {
    return err
}
if err = LOG_DB.AutoMigrate(&Conversation{}); err != nil {
    return err
}
if err = LOG_DB.AutoMigrate(&ConversationArchive{}); err != nil {
    return err
}
if err = LOG_DB.AutoMigrate(&OfficialNewTable{}); err != nil {
    return err
}
```

**文件 3: `router/api-router.go`**

```go
// 保留你的对话记录路由
conversationRoute := apiRouter.Group("/conversation")
conversationRoute.Use(middleware.AdminAuth())
{
    // ... 你的路由
}

// 添加官方的新路由
officialNewRoute := apiRouter.Group("/official-new")
officialNewRoute.Use(middleware.AdminAuth())
{
    // ... 官方的路由
}
```

#### 5.4 完成合并

```bash
# 1. 解决所有冲突后，添加文件
git add .

# 2. 查看状态
git status

# 3. 完成合并提交
git commit -m "Merge upstream/main into main

- Merged official updates
- Preserved conversation logging feature
- Resolved conflicts in:
  - common/constants.go
  - model/main.go
  - router/api-router.go
"

# 4. 推送到你的仓库
git push origin main
```

---

### 第六步：定期同步官方更新（日常维护）

#### 6.1 创建同步脚本

创建 `sync-upstream.sh`：

```bash
#!/bin/bash

# 官方仓库同步脚本
# 用法: ./sync-upstream.sh

set -e

echo "🔄 开始同步官方仓库..."

# 1. 获取官方最新代码
echo "📥 获取官方最新代码..."
git fetch upstream

# 2. 查看官方更新
echo "📋 官方最新提交："
git log --oneline upstream/main ^main -10

# 3. 询问是否继续
read -p "是否继续合并？(y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "❌ 取消合并"
    exit 1
fi

# 4. 备份当前分支
BACKUP_BRANCH="backup-$(date +%Y%m%d-%H%M%S)"
git branch $BACKUP_BRANCH
echo "✅ 已创建备份分支: $BACKUP_BRANCH"

# 5. 合并官方代码
echo "🔀 开始合并..."
if git merge upstream/main --no-edit; then
    echo "✅ 合并成功！"

    # 6. 推送到远程
    read -p "是否推送到远程仓库？(y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]
    then
        git push origin main
        echo "✅ 已推送到远程仓库"
    fi
else
    echo "⚠️  检测到冲突，请手动解决"
    echo "冲突文件："
    git status --short | grep "^UU"
    echo ""
    echo "解决冲突后执行："
    echo "  git add ."
    echo "  git commit"
    echo "  git push origin main"
    echo ""
    echo "如需放弃合并："
    echo "  git merge --abort"
    echo "  git checkout $BACKUP_BRANCH"
fi

echo "✅ 完成！"
```

#### 6.2 使用同步脚本

```bash
# 赋予执行权限
chmod +x sync-upstream.sh

# 执行同步
./sync-upstream.sh
```

#### 6.3 手动同步步骤

```bash
# 每次官方更新后执行：

# 1. 切换到 main 分支
git checkout main

# 2. 获取官方最新代码
git fetch upstream

# 3. 查看官方更新内容
git log upstream/main ^main --oneline

# 4. 合并官方代码
git merge upstream/main

# 5. 如有冲突，解决冲突
# （参考第五步）

# 6. 推送到你的仓库
git push origin main
```

---

## 🛡️ 防止冲突的最佳实践

### 1. 遵循最小修改原则（已做到✅）

你的二开项目已经遵循了最小修改原则：
- ✅ 只修改了 3 个文件（19 行）
- ✅ 新增功能都在独立的文件中
- ✅ 不修改官方的核心逻辑

### 2. 使用独立的命名空间

```go
// ✅ 好的做法（已采用）
var ConversationLogEnabled = false
type Conversation struct {}
func RecordConversation() {}

// ❌ 不好的做法（会冲突）
var LogEnabled = false  // 名字太通用
type Record struct {}   // 名字太通用
```

### 3. 模块化设计

```
你的功能在独立的文件中：
model/
  ├── conversation.go            # 你的
  ├── conversation_archive.go    # 你的
  └── log.go                     # 官方的 ✅ 不冲突

controller/
  ├── conversation.go            # 你的
  └── log.go                     # 官方的 ✅ 不冲突
```

### 4. 使用 Git 属性标记（可选）

创建 `.gitattributes` 文件：

```bash
# 对于你修改过的文件，使用自定义合并策略
common/constants.go merge=union
model/main.go merge=union
router/api-router.go merge=union
```

---

## 📊 冲突矩阵（可能冲突的文件）

| 文件 | 冲突概率 | 处理策略 |
|------|---------|---------|
| `common/constants.go` | 中 | 保留双方修改，手动合并 |
| `model/main.go` | 低 | 保留双方的 AutoMigrate |
| `router/api-router.go` | 低 | 保留双方的路由组 |
| `docker-compose.yml` | 高 | 使用你的版本，手动添加官方新配置 |
| `go.mod` | 中 | 使用官方版本，保留你的依赖 |
| `model/conversation.go` | 无 | 你的独立文件 |
| `controller/conversation.go` | 无 | 你的独立文件 |
| 其他你的新文件 | 无 | 完全独立 |

---

## 🔍 测试合并后的代码

### 1. 编译测试

```bash
# 1. 编译后端
go mod tidy
go build -o new-api.exe

# 2. 编译前端
cd web
npm install
npm run build
cd ..

# 3. 运行测试
go test ./...
```

### 2. 功能测试

```bash
# 1. 启动服务
./new-api.exe

# 2. 测试官方功能
curl http://localhost:3000/api/status

# 3. 测试你的功能
curl http://localhost:3000/api/conversation/setting \
  -H "Authorization: Bearer TOKEN"

# 4. 测试对话记录
# 发送一个 AI 请求，检查是否正常记录
```

### 3. Docker 测试

```bash
# 1. 构建镜像
docker build -t zhang/new-api:test .

# 2. 启动测试
docker run -d -p 3000:3000 zhang/new-api:test

# 3. 测试功能
curl http://localhost:3000/api/status
```

---

## 📝 合并检查清单

合并官方代码后，确保：

- [ ] 代码能正常编译（`go build`）
- [ ] 所有测试通过（`go test ./...`）
- [ ] 官方的新功能正常工作
- [ ] 你的对话记录功能正常工作
- [ ] 数据库迁移正常（conversations 表存在）
- [ ] Docker 镜像能正常构建
- [ ] docker-compose 能正常启动
- [ ] 前端页面正常访问
- [ ] API 接口都能正常调用
- [ ] 没有破坏性的修改

---

## 🆘 常见问题处理

### Q1: 合并时出现大量冲突怎么办？

**方案 1：重新应用你的修改（推荐）**

```bash
# 1. 放弃当前合并
git merge --abort

# 2. 创建新分支基于官方最新代码
git checkout -b rebase-conversation upstream/main

# 3. 手动应用你的修改
# 从你的备份分支复制文件
cp ../backup/model/conversation.go model/
cp ../backup/controller/conversation.go controller/
# ... 复制其他你的新增文件

# 4. 手动修改冲突的文件
# 编辑 common/constants.go, model/main.go, router/api-router.go

# 5. 提交
git add .
git commit -m "Re-apply conversation logging feature on latest upstream"

# 6. 测试无误后，替换 main 分支
git checkout main
git reset --hard rebase-conversation
```

**方案 2：使用 cherry-pick（适合少量提交）**

```bash
# 1. 查看你的功能提交
git log --oneline main ^upstream/main

# 2. 切换到官方最新代码
git checkout -b new-main upstream/main

# 3. 挑选你的提交
git cherry-pick <commit-hash-1>
git cherry-pick <commit-hash-2>

# 4. 解决冲突（如果有）
git add .
git cherry-pick --continue

# 5. 替换 main 分支
git checkout main
git reset --hard new-main
```

### Q2: 如何确保不会丢失我的修改？

```bash
# 1. 合并前创建备份分支
git branch backup-before-merge-$(date +%Y%m%d)

# 2. 或者创建备份标签
git tag backup-v1.0.0-before-merge

# 3. 如果出错，恢复备份
git checkout backup-before-merge-20250128
git checkout -b recover-branch

# 4. 或者恢复到标签
git checkout backup-v1.0.0-before-merge
```

### Q3: 官方删除了我依赖的文件怎么办？

```bash
# 1. 查看官方删除了什么
git diff upstream/main --name-status | grep "^D"

# 2. 如果是你依赖的文件，从旧版本恢复
git checkout HEAD -- path/to/deleted/file

# 3. 或者修改你的代码，使用官方的新方式
```

### Q4: 如何跟踪官方的 release 版本？

```bash
# 1. 获取官方的所有标签
git fetch upstream --tags

# 2. 查看官方的 release 版本
git tag -l | grep -E "^v[0-9]"

# 3. 切换到特定版本
git checkout -b sync-v0.7.0 v0.7.0

# 4. 合并到你的代码
git checkout main
git merge sync-v0.7.0
```

---

## 📦 完整的工作流程图

```
官方仓库 (QuantumNous/new-api)
        │
        │ git fetch upstream
        ▼
    upstream/main
        │
        │ git merge upstream/main
        ▼
    你的 main 分支
        │
        │ 解决冲突（如果有）
        ▼
    提交合并
        │
        │ git push origin main
        ▼
    你的远程仓库 (zhang/new-api)
        │
        │ 部署
        ▼
    生产环境
```

---

## 🎯 最佳实践总结

### ✅ 做的好的地方（你已经做到了）

1. **最小修改原则** - 只修改了 19 行代码
2. **独立文件** - 新功能都在独立文件中
3. **命名空间** - 使用 Conversation 前缀，避免冲突
4. **模块化设计** - 功能独立，易于维护

### 🔄 日常维护建议

1. **定期同步**（建议每月一次）
   ```bash
   ./sync-upstream.sh
   ```

2. **关注官方更新**
   - 订阅官方仓库的 Release
   - 关注 CHANGELOG

3. **测试后再合并**
   - 在测试分支先合并
   - 测试通过后再合并到 main

4. **保持文档更新**
   - 记录每次合并的变化
   - 更新你的 README

---

## 📚 相关资源

- [Git 官方文档 - 分支管理](https://git-scm.com/book/zh/v2/Git-%E5%88%86%E6%94%AF-%E5%88%86%E6%94%AF%E7%9A%84%E6%96%B0%E5%BB%BA%E4%B8%8E%E5%90%88%E5%B9%B6)
- [GitHub - Fork 同步](https://docs.github.com/zh/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork)
- [Pro Git 中文版](https://git-scm.com/book/zh/v2)

---

**祝你顺利合并！** 🎉

有任何问题随时查阅本文档。
