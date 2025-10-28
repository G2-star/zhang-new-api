package model

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// StartConversationMaintenance 启动对话记录自动维护任务
// 此函数应该在 main.go 中调用
func StartConversationMaintenance() {
	// 如果对话记录功能未启用，不启动维护任务
	if !common.ConversationLogEnabled {
		common.SysLog("对话记录功能未启用，跳过维护任务")
		return
	}

	common.SysLog("启动对话记录自动维护任务")

	// 任务1：每天自动归档旧对话
	go autoArchiveTask()

	// 任务2：每周自动清理归档数据
	go autoCleanupTask()

	// 任务3：每月自动优化数据库表
	go autoOptimizeTask()
}

// autoArchiveTask 每天自动归档旧对话
func autoArchiveTask() {
	// 首次执行延迟1分钟
	time.Sleep(1 * time.Minute)

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		// 从环境变量读取归档天数，默认30天
		archiveDays := getEnvInt("CONVERSATION_ARCHIVE_DAYS", 30)
		if archiveDays <= 0 {
			archiveDays = 30
		}

		targetTime := time.Now().AddDate(0, 0, -archiveDays).Unix()
		batchSize := 1000

		common.SysLog("开始自动归档对话记录...")
		ctx := context.Background()
		archived, err := ArchiveOldConversations(ctx, targetTime, batchSize)
		if err != nil {
			common.SysError("自动归档失败: " + err.Error())
		} else if archived > 0 {
			common.SysLog(common.GetJsonString(map[string]interface{}{
				"task":     "auto_archive",
				"archived": archived,
				"days":     archiveDays,
			}))
		}

		<-ticker.C
	}
}

// autoCleanupTask 每周自动清理归档数据
func autoCleanupTask() {
	// 首次执行延迟2分钟
	time.Sleep(2 * time.Minute)

	ticker := time.NewTicker(7 * 24 * time.Hour)
	defer ticker.Stop()

	for {
		// 从环境变量读取清理天数，默认365天
		cleanupDays := getEnvInt("CONVERSATION_CLEANUP_DAYS", 365)
		if cleanupDays <= 0 {
			cleanupDays = 365
		}

		targetTime := time.Now().AddDate(0, 0, -cleanupDays).Unix()
		batchSize := 1000

		common.SysLog("开始自动清理归档数据...")
		ctx := context.Background()
		deleted, err := CleanupOldArchives(ctx, targetTime, batchSize)
		if err != nil {
			common.SysError("自动清理失败: " + err.Error())
		} else if deleted > 0 {
			common.SysLog(common.GetJsonString(map[string]interface{}{
				"task":    "auto_cleanup",
				"deleted": deleted,
				"days":    cleanupDays,
			}))
		}

		<-ticker.C
	}
}

// autoOptimizeTask 每月自动优化数据库表
func autoOptimizeTask() {
	// 首次执行延迟3分钟
	time.Sleep(3 * time.Minute)

	ticker := time.NewTicker(30 * 24 * time.Hour)
	defer ticker.Stop()

	for {
		common.SysLog("开始自动优化对话记录表...")

		// 优化主表和归档表
		if err := OptimizeConversationTable(); err != nil {
			common.SysError("优化表失败: " + err.Error())
		} else {
			common.SysLog("表优化完成")
		}

		<-ticker.C
	}
}

// getEnvInt 从环境变量获取整数值
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
