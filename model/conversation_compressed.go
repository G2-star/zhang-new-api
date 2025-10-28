package model

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/QuantumNous/new-api/logger"
	"github.com/gin-gonic/gin"
)

// ConversationCompressed 压缩存储版本的对话记录
// 适用于长对话或需要节省存储空间的场景
type ConversationCompressed struct {
	Id               int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId           int    `json:"user_id" gorm:"index:idx_user_model_time;index;not null"`
	Username         string `json:"username" gorm:"index;not null;default:''"`
	ModelName        string `json:"model_name" gorm:"index:idx_user_model_time;index;not null;default:''"`
	TokenId          int    `json:"token_id" gorm:"index;default:0"`
	TokenName        string `json:"token_name" gorm:"default:''"`
	ChannelId        int    `json:"channel_id" gorm:"index;default:0"`
	RequestMessagesGz  []byte `json:"-" gorm:"type:blob"`                         // GZIP 压缩的请求消息
	ResponseContentGz  []byte `json:"-" gorm:"type:blob"`                         // GZIP 压缩的响应内容
	PromptTokens     int    `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens int    `json:"completion_tokens" gorm:"default:0"`
	TotalTokens      int    `json:"total_tokens" gorm:"default:0"`
	IsStream         bool   `json:"is_stream" gorm:"default:false"`
	CreatedAt        int64  `json:"created_at" gorm:"bigint;index:idx_user_model_time;index;not null"`
	UseTime          int    `json:"use_time" gorm:"default:0"`
	Ip               string `json:"ip" gorm:"index;default:''"`
	Group            string `json:"group" gorm:"index;default:''"`
	CompressionRatio float64 `json:"compression_ratio" gorm:"default:0"` // 压缩率
}

func (ConversationCompressed) TableName() string {
	return "conversations_compressed"
}

// CompressString 压缩字符串
func CompressString(s string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write([]byte(s))
	if err != nil {
		gz.Close()
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecompressString 解压字符串
func DecompressString(compressed []byte) (string, error) {
	if len(compressed) == 0 {
		return "", nil
	}

	buf := bytes.NewReader(compressed)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	decompressed, err := io.ReadAll(gz)
	if err != nil {
		return "", err
	}

	return string(decompressed), nil
}

// RecordConversationCompressed 记录压缩版本的对话
func RecordConversationCompressed(c *gin.Context, params RecordConversationParams) {
	// 将 RequestMessages 转为 JSON 字符串
	requestJSON, err := json.Marshal(params.RequestMessages)
	if err != nil {
		logger.LogError(c, "failed to marshal request messages: "+err.Error())
		return
	}

	// 压缩请求消息
	requestCompressed, err := CompressString(string(requestJSON))
	if err != nil {
		logger.LogError(c, "failed to compress request: "+err.Error())
		return
	}

	// 压缩响应内容
	responseCompressed, err := CompressString(params.ResponseContent)
	if err != nil {
		logger.LogError(c, "failed to compress response: "+err.Error())
		return
	}

	// 计算压缩率
	originalSize := len(requestJSON) + len(params.ResponseContent)
	compressedSize := len(requestCompressed) + len(responseCompressed)
	compressionRatio := 0.0
	if originalSize > 0 {
		compressionRatio = float64(compressedSize) / float64(originalSize)
	}

	conversation := &ConversationCompressed{
		UserId:           params.UserId,
		Username:         params.Username,
		ModelName:        params.ModelName,
		TokenId:          params.TokenId,
		TokenName:        params.TokenName,
		ChannelId:        params.ChannelId,
		RequestMessagesGz:  requestCompressed,
		ResponseContentGz:  responseCompressed,
		PromptTokens:     params.PromptTokens,
		CompletionTokens: params.CompletionTokens,
		TotalTokens:      params.PromptTokens + params.CompletionTokens,
		IsStream:         params.IsStream,
		CreatedAt:        params.CreatedAt,
		UseTime:          params.UseTime,
		Ip:               params.Ip,
		Group:            params.Group,
		CompressionRatio: compressionRatio,
	}

	err = LOG_DB.Create(conversation).Error
	if err != nil {
		logger.LogError(c, "failed to record compressed conversation: "+err.Error())
	}
}

// GetCompressedConversationById 获取压缩版本的对话并解压
func GetCompressedConversationById(id int) (*Conversation, error) {
	var compressed ConversationCompressed
	err := LOG_DB.Where("id = ?", id).First(&compressed).Error
	if err != nil {
		return nil, err
	}

	// 解压
	requestMessages, err := DecompressString(compressed.RequestMessagesGz)
	if err != nil {
		return nil, err
	}

	responseContent, err := DecompressString(compressed.ResponseContentGz)
	if err != nil {
		return nil, err
	}

	// 转换为普通对话记录
	conversation := &Conversation{
		Id:               compressed.Id,
		UserId:           compressed.UserId,
		Username:         compressed.Username,
		ModelName:        compressed.ModelName,
		TokenId:          compressed.TokenId,
		TokenName:        compressed.TokenName,
		ChannelId:        compressed.ChannelId,
		RequestMessages:  requestMessages,
		ResponseContent:  responseContent,
		PromptTokens:     compressed.PromptTokens,
		CompletionTokens: compressed.CompletionTokens,
		TotalTokens:      compressed.TotalTokens,
		IsStream:         compressed.IsStream,
		CreatedAt:        compressed.CreatedAt,
		UseTime:          compressed.UseTime,
		Ip:               compressed.Ip,
		Group:            compressed.Group,
	}

	return conversation, nil
}

// MigrateToCompressed 将普通对话迁移到压缩表
func MigrateToCompressed(oldConvId int) error {
	// 1. 读取原对话
	var conv Conversation
	err := LOG_DB.Where("id = ?", oldConvId).First(&conv).Error
	if err != nil {
		return err
	}

	// 2. 压缩数据
	requestCompressed, err := CompressString(conv.RequestMessages)
	if err != nil {
		return err
	}

	responseCompressed, err := CompressString(conv.ResponseContent)
	if err != nil {
		return err
	}

	// 计算压缩率
	originalSize := len(conv.RequestMessages) + len(conv.ResponseContent)
	compressedSize := len(requestCompressed) + len(responseCompressed)
	compressionRatio := float64(compressedSize) / float64(originalSize)

	// 3. 插入到压缩表
	compressed := &ConversationCompressed{
		UserId:           conv.UserId,
		Username:         conv.Username,
		ModelName:        conv.ModelName,
		TokenId:          conv.TokenId,
		TokenName:        conv.TokenName,
		ChannelId:        conv.ChannelId,
		RequestMessagesGz:  requestCompressed,
		ResponseContentGz:  responseCompressed,
		PromptTokens:     conv.PromptTokens,
		CompletionTokens: conv.CompletionTokens,
		TotalTokens:      conv.TotalTokens,
		IsStream:         conv.IsStream,
		CreatedAt:        conv.CreatedAt,
		UseTime:          conv.UseTime,
		Ip:               conv.Ip,
		Group:            conv.Group,
		CompressionRatio: compressionRatio,
	}

	err = LOG_DB.Create(compressed).Error
	if err != nil {
		return err
	}

	// 4. 删除原记录
	err = LOG_DB.Delete(&conv).Error
	return err
}
