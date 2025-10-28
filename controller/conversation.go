package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

// GetConversations 获取对话记录列表（支持筛选）
func GetConversations(c *gin.Context) {
	// 权限检查：只有管理员可以查看
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

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	startIdx := (page - 1) * pageSize

	// 查询数据
	conversations, total, err := model.GetConversations(userId, modelName, username, startTime, endTime, startIdx, pageSize)
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
			"data":      conversations,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetConversationDetail 获取单条对话详情
func GetConversationDetail(c *gin.Context) {
	// 权限检查：只有管理员可以查看
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	conversation, err := model.GetConversationById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "对话记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    conversation,
	})
}

// DeleteConversations 批量删除对话记录
func DeleteConversations(c *gin.Context) {
	// 权限检查：只有管理员可以删除
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	var req struct {
		Ids []int `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	if len(req.Ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "未选择任何记录",
		})
		return
	}

	err := model.DeleteConversations(req.Ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
		"data":    gin.H{"deleted": len(req.Ids)},
	})
}

// DeleteConversationsByCondition 按条件批量删除对话记录
func DeleteConversationsByCondition(c *gin.Context) {
	// 权限检查：只有管理员可以删除
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	var req struct {
		UserId    int    `json:"user_id"`
		ModelName string `json:"model_name"`
		Username  string `json:"username"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	// 至少需要一个筛选条件
	if req.UserId == 0 && req.ModelName == "" && req.Username == "" && req.StartTime == 0 && req.EndTime == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "至少需要一个筛选条件",
		})
		return
	}

	deleted, err := model.DeleteConversationsByCondition(req.UserId, req.ModelName, req.Username, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
		"data":    gin.H{"deleted": deleted},
	})
}

// GetConversationStats 获取对话统计信息
func GetConversationStats(c *gin.Context) {
	// 权限检查：只有管理员可以查看
	currentUserId := c.GetInt("id")
	if !model.IsAdmin(currentUserId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	userId, _ := strconv.Atoi(c.Query("user_id"))
	modelName := c.Query("model_name")
	startTime, _ := strconv.ParseInt(c.Query("start_time"), 10, 64)
	endTime, _ := strconv.ParseInt(c.Query("end_time"), 10, 64)

	stats, err := model.GetConversationStats(userId, modelName, startTime, endTime)
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

// UpdateConversationLogSetting 更新对话记录功能开关
func UpdateConversationLogSetting(c *gin.Context) {
	// 权限检查：只有管理员可以修改
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}

	common.ConversationLogEnabled = req.Enabled

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "设置已更新",
		"data":    gin.H{"enabled": req.Enabled},
	})
}

// GetConversationLogSetting 获取对话记录功能开关状态
func GetConversationLogSetting(c *gin.Context) {
	// 权限检查：只有管理员可以查看
	userId := c.GetInt("id")
	if !model.IsAdmin(userId) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限访问",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    gin.H{"enabled": common.ConversationLogEnabled},
	})
}
