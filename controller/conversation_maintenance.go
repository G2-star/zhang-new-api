package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

// ArchiveOldConversations 归档旧对话（管理员接口）
func ArchiveOldConversations(c *gin.Context) {
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	var req struct {
		Days      int `json:"days"`       // 归档多少天前的数据
		BatchSize int `json:"batch_size"` // 每批处理数量
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	// 默认值
	if req.Days <= 0 {
		req.Days = 30 // 默认归档30天前的数据
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 1000
	}

	// 计算目标时间戳
	targetTime := time.Now().AddDate(0, 0, -req.Days).Unix()

	// 异步执行归档
	go func() {
		ctx := context.Background()
		archived, err := model.ArchiveOldConversations(ctx, targetTime, req.BatchSize)
		if err != nil {
			common.SysLog("归档对话失败: " + err.Error())
		} else {
			common.SysLog(common.GetJsonString(map[string]interface{}{
				"action":   "archive_conversations",
				"archived": archived,
				"days":     req.Days,
			}))
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "归档任务已启动，将在后台执行",
		"data": gin.H{
			"target_days": req.Days,
			"batch_size":  req.BatchSize,
		},
	})
}

// CleanupOldArchives 清理归档表中的超旧数据
func CleanupOldArchives(c *gin.Context) {
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	var req struct {
		Days      int `json:"days"`       // 删除多少天前的归档数据
		BatchSize int `json:"batch_size"` // 每批处理数量
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	// 默认值
	if req.Days <= 0 {
		req.Days = 365 // 默认删除1年前的归档数据
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 1000
	}

	// 计算目标时间戳
	targetTime := time.Now().AddDate(0, 0, -req.Days).Unix()

	// 异步执行清理
	go func() {
		ctx := context.Background()
		deleted, err := model.CleanupOldArchives(ctx, targetTime, req.BatchSize)
		if err != nil {
			common.SysLog("清理归档数据失败: " + err.Error())
		} else {
			common.SysLog(common.GetJsonString(map[string]interface{}{
				"action":  "cleanup_archives",
				"deleted": deleted,
				"days":    req.Days,
			}))
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "清理任务已启动，将在后台执行",
		"data": gin.H{
			"target_days": req.Days,
			"batch_size":  req.BatchSize,
		},
	})
}

// OptimizeConversationTables 优化对话表（重建索引、回收空间）
func OptimizeConversationTables(c *gin.Context) {
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	// 异步执行优化
	go func() {
		err := model.OptimizeConversationTable()
		if err != nil {
			common.SysLog("优化对话表失败: " + err.Error())
		} else {
			common.SysLog("对话表优化完成")
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "优化任务已启动，将在后台执行",
	})
}

// GetConversationTableStats 获取对话表统计信息
func GetConversationTableStats(c *gin.Context) {
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	stats, err := model.GetConversationTableStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    stats,
	})
}

// SearchArchivedConversations 搜索归档表（扩展查询功能）
func SearchArchivedConversations(c *gin.Context) {
	currentUserId := c.GetInt("id")
	if !model.IsAdmin(currentUserId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	// 获取筛选参数
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	userId, _ := strconv.Atoi(c.Query("user_id"))
	modelName := c.Query("model_name")
	username := c.Query("username")
	startTime, _ := strconv.ParseInt(c.Query("start_time"), 10, 64)
	endTime, _ := strconv.ParseInt(c.Query("end_time"), 10, 64)
	searchArchive := c.Query("search_archive") == "true"

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	startIdx := (page - 1) * pageSize

	// 查询数据（包含归档表）
	conversations, total, err := model.GetConversationsWithArchive(userId, modelName, username, startTime, endTime, startIdx, pageSize, searchArchive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"data":           conversations,
			"total":          total,
			"page":           page,
			"page_size":      pageSize,
			"search_archive": searchArchive,
		},
	})
}
