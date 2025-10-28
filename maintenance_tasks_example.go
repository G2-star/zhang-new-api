package main

// 在 main.go 中添加此代码段

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// StartConversationMaintenanceTasks 启动对话记录维护任务
// 在 main 函数中调用此函数
func StartConversationMaintenanceTasks() {
	// 如果对话记录功能未启用，不启动维护任务
	if !common.ConversationLogEnabled {
		return
	}

	common.SysLog("启动对话记录维护任务...")

	// 任务1：每天自动归档旧对话（30天前的数据）
	go func() {
		// 首次执行延迟1分钟，避免启动时的资源争用
		time.Sleep(1 * time.Minute)

		ticker := time.NewTicker(24 * time.Hour) // 每天执行一次
		defer ticker.Stop()

		for {
			archiveDays := 30 // 归档30天前的数据
			if days := GetEnvInt("CONVERSATION_ARCHIVE_DAYS", 30); days > 0 {
				archiveDays = days
			}

			targetTime := time.Now().AddDate(0, 0, -archiveDays).Unix()
			batchSize := 1000

			common.SysLog("开始归档旧对话...")
			ctx := context.Background()
			archived, err := model.ArchiveOldConversations(ctx, targetTime, batchSize)
			if err != nil {
				common.SysLog("归档对话失败: " + err.Error())
			} else if archived > 0 {
				common.SysLog(common.GetJsonString(map[string]interface{}{
					"message":  "对话归档完成",
					"archived": archived,
					"days":     archiveDays,
				}))
			}

			<-ticker.C
		}
	}()

	// 任务2：每周清理归档表中的超旧数据（1年前的数据）
	go func() {
		// 首次执行延迟5分钟
		time.Sleep(5 * time.Minute)

		ticker := time.NewTicker(7 * 24 * time.Hour) // 每周执行一次
		defer ticker.Stop()

		for {
			cleanupDays := 365 // 清理1年前的归档数据
			if days := GetEnvInt("CONVERSATION_CLEANUP_DAYS", 365); days > 0 {
				cleanupDays = days
			}

			targetTime := time.Now().AddDate(0, 0, -cleanupDays).Unix()
			batchSize := 1000

			common.SysLog("开始清理超旧归档数据...")
			ctx := context.Background()
			deleted, err := model.CleanupOldArchives(ctx, targetTime, batchSize)
			if err != nil {
				common.SysLog("清理归档数据失败: " + err.Error())
			} else if deleted > 0 {
				common.SysLog(common.GetJsonString(map[string]interface{}{
					"message": "归档数据清理完成",
					"deleted": deleted,
					"days":    cleanupDays,
				}))
			}

			<-ticker.C
		}
	}()

	// 任务3：每月优化数据库表
	go func() {
		// 首次执行延迟10分钟
		time.Sleep(10 * time.Minute)

		ticker := time.NewTicker(30 * 24 * time.Hour) // 每月执行一次
		defer ticker.Stop()

		for {
			common.SysLog("开始优化对话表...")
			err := model.OptimizeConversationTable()
			if err != nil {
				common.SysLog("优化对话表失败: " + err.Error())
			} else {
				common.SysLog("对话表优化完成")
			}

			<-ticker.C
		}
	}()

	// 任务4：每小时检查表大小并记录
	go func() {
		time.Sleep(15 * time.Minute)

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			stats, err := model.GetConversationTableStats()
			if err == nil {
				// 只在主表记录数超过10万时记录日志
				if mainCount, ok := stats["main_table_count"].(int64); ok && mainCount > 100000 {
					common.SysLog(common.GetJsonString(map[string]interface{}{
						"message":            "对话表统计",
						"main_table_count":   stats["main_table_count"],
						"archive_table_count": stats["archive_table_count"],
					}))
				}
			}

			<-ticker.C
		}
	}()
}

// GetEnvInt 从环境变量获取整数值的辅助函数
// 如果环境变量不存在或无法解析，返回默认值
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

/*
使用方法：
在 main.go 的 main() 函数中添加：

func main() {
    // ... 现有的初始化代码 ...

    // 启动对话记录维护任务
    StartConversationMaintenanceTasks()

    // ... 启动 HTTP 服务器 ...
}

环境变量配置（.env 文件）：

# 对话记录自动归档天数（默认30天）
CONVERSATION_ARCHIVE_DAYS=30

# 归档数据自动清理天数（默认365天）
CONVERSATION_CLEANUP_DAYS=365

# 对话记录功能开关
CONVERSATION_LOG_ENABLED=true
*/
