package model

import (
	"context"
	"encoding/json"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/gin-gonic/gin"
)

// Conversation 对话记录表 - 用于记录完整的AI对话内容
type Conversation struct {
	Id               int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId           int    `json:"user_id" gorm:"index:idx_user_model_time;index;not null"`
	Username         string `json:"username" gorm:"index;not null;default:''"`
	ModelName        string `json:"model_name" gorm:"index:idx_user_model_time;index;not null;default:''"`
	TokenId          int    `json:"token_id" gorm:"index;default:0"`
	TokenName        string `json:"token_name" gorm:"default:''"`
	ChannelId        int    `json:"channel_id" gorm:"index;default:0"`
	RequestMessages  string `json:"request_messages" gorm:"type:text"` // JSON格式的请求消息
	ResponseContent  string `json:"response_content" gorm:"type:text"` // AI响应内容
	PromptTokens     int    `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens int    `json:"completion_tokens" gorm:"default:0"`
	TotalTokens      int    `json:"total_tokens" gorm:"default:0"`
	IsStream         bool   `json:"is_stream" gorm:"default:false"`
	CreatedAt        int64  `json:"created_at" gorm:"bigint;index:idx_user_model_time;index;not null"` // Unix 时间戳
	UseTime          int    `json:"use_time" gorm:"default:0"`                                          // 响应时间（毫秒）
	Ip               string `json:"ip" gorm:"index;default:''"`
	Group            string `json:"group" gorm:"index;default:''"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// RecordConversationParams 记录对话的参数
type RecordConversationParams struct {
	UserId           int
	Username         string
	ModelName        string
	TokenId          int
	TokenName        string
	ChannelId        int
	RequestMessages  interface{} // 支持 []dto.Message 或其他格式
	ResponseContent  string
	PromptTokens     int
	CompletionTokens int
	IsStream         bool
	CreatedAt        int64 // 添加创建时间字段
	UseTime          int
	Ip               string
	Group            string
}

// RecordConversation 记录对话内容
func RecordConversation(c *gin.Context, params RecordConversationParams) {
	// 检查是否启用对话记录功能
	if !common.ConversationLogEnabled {
		return
	}

	// 将 RequestMessages 转为 JSON 字符串
	requestJSON, err := json.Marshal(params.RequestMessages)
	if err != nil {
		logger.LogError(c, "failed to marshal request messages: "+err.Error())
		return
	}

	conversation := &Conversation{
		UserId:           params.UserId,
		Username:         params.Username,
		ModelName:        params.ModelName,
		TokenId:          params.TokenId,
		TokenName:        params.TokenName,
		ChannelId:        params.ChannelId,
		RequestMessages:  string(requestJSON),
		ResponseContent:  params.ResponseContent,
		PromptTokens:     params.PromptTokens,
		CompletionTokens: params.CompletionTokens,
		TotalTokens:      params.PromptTokens + params.CompletionTokens,
		IsStream:         params.IsStream,
		CreatedAt:        common.GetTimestamp(),
		UseTime:          params.UseTime,
		Ip:               params.Ip,
		Group:            params.Group,
	}

	err = LOG_DB.Create(conversation).Error
	if err != nil {
		logger.LogError(c, "failed to record conversation: "+err.Error())
	}
}

// GetConversations 查询对话记录（支持分页和筛选）
func GetConversations(userId int, modelName string, username string, startTime int64, endTime int64, startIdx int, num int) ([]*Conversation, int64, error) {
	var conversations []*Conversation
	var total int64

	tx := LOG_DB.Model(&Conversation{})

	// 筛选条件
	if userId > 0 {
		tx = tx.Where("user_id = ?", userId)
	}
	if modelName != "" {
		tx = tx.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if username != "" {
		tx = tx.Where("username = ?", username)
	}
	if startTime > 0 {
		tx = tx.Where("created_at >= ?", startTime)
	}
	if endTime > 0 {
		tx = tx.Where("created_at <= ?", endTime)
	}

	// 获取总数
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = tx.Order("created_at DESC").Limit(num).Offset(startIdx).Find(&conversations).Error
	return conversations, total, err
}

// GetConversationById 根据ID查询单条对话记录
func GetConversationById(id int) (*Conversation, error) {
	var conversation Conversation
	err := LOG_DB.Where("id = ?", id).First(&conversation).Error
	return &conversation, err
}

// DeleteConversations 批量删除对话记录
func DeleteConversations(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	return LOG_DB.Where("id IN ?", ids).Delete(&Conversation{}).Error
}

// DeleteConversationsByCondition 按条件批量删除对话记录
func DeleteConversationsByCondition(userId int, modelName string, username string, startTime int64, endTime int64) (int64, error) {
	tx := LOG_DB.Model(&Conversation{})

	// 筛选条件
	if userId > 0 {
		tx = tx.Where("user_id = ?", userId)
	}
	if modelName != "" {
		tx = tx.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if username != "" {
		tx = tx.Where("username = ?", username)
	}
	if startTime > 0 {
		tx = tx.Where("created_at >= ?", startTime)
	}
	if endTime > 0 {
		tx = tx.Where("created_at <= ?", endTime)
	}

	result := tx.Delete(&Conversation{})
	return result.RowsAffected, result.Error
}

// DeleteOldConversations 删除旧的对话记录（用于定时清理）
func DeleteOldConversations(ctx context.Context, targetTimestamp int64, limit int) (int64, error) {
	var total int64 = 0

	for {
		if ctx.Err() != nil {
			return total, ctx.Err()
		}

		result := LOG_DB.Where("created_at < ?", targetTimestamp).Limit(limit).Delete(&Conversation{})
		if result.Error != nil {
			return total, result.Error
		}

		total += result.RowsAffected

		if result.RowsAffected < int64(limit) {
			break
		}
	}

	return total, nil
}

// GetConversationStats 获取对话统计信息
func GetConversationStats(userId int, modelName string, startTime int64, endTime int64) (map[string]interface{}, error) {
	var stats struct {
		TotalCount       int64 `json:"total_count"`
		TotalTokens      int64 `json:"total_tokens"`
		TotalPromptTokens int64 `json:"total_prompt_tokens"`
		TotalCompTokens  int64 `json:"total_comp_tokens"`
	}

	tx := LOG_DB.Model(&Conversation{}).
		Select("COUNT(*) as total_count, SUM(total_tokens) as total_tokens, SUM(prompt_tokens) as total_prompt_tokens, SUM(completion_tokens) as total_comp_tokens")

	if userId > 0 {
		tx = tx.Where("user_id = ?", userId)
	}
	if modelName != "" {
		tx = tx.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if startTime > 0 {
		tx = tx.Where("created_at >= ?", startTime)
	}
	if endTime > 0 {
		tx = tx.Where("created_at <= ?", endTime)
	}

	err := tx.Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"total_count":        stats.TotalCount,
		"total_tokens":       stats.TotalTokens,
		"total_prompt_tokens": stats.TotalPromptTokens,
		"total_comp_tokens":  stats.TotalCompTokens,
	}

	return result, nil
}
