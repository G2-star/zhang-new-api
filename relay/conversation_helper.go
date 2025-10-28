package relay

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
)

// RecordConversationHelper 记录对话的辅助函数
// 在响应完成后异步记录对话内容
func RecordConversationHelper(c *gin.Context, info *relaycommon.RelayInfo, textReq *dto.GeneralOpenAIRequest, responseContent string, usage *dto.Usage, startTime time.Time) {
	// 检查是否启用对话记录
	if !common.ConversationLogEnabled {
		return
	}

	// 只记录聊天对话类型
	if textReq == nil || textReq.Messages == nil || len(textReq.Messages) == 0 {
		return
	}

	// 提取响应内容
	if responseContent == "" {
		return
	}

	// 判断是否需要记录 IP
	needRecordIp := false
	if settingMap, err := model.GetUserSetting(info.UserId, false); err == nil {
		if settingMap.RecordIpLog {
			needRecordIp = true
		}
	}

	// 计算使用时间（毫秒）
	useTime := int(time.Since(startTime).Milliseconds())

	// 构造记录参数
	params := model.RecordConversationParams{
		UserId:           info.UserId,
		Username:         c.GetString("username"),
		ModelName:        info.UpstreamModelName,
		TokenId:          info.TokenId,
		TokenName:        info.TokenKey, // 使用 TokenKey 作为 TokenName
		ChannelId:        info.ChannelId,
		RequestMessages:  textReq.Messages,
		ResponseContent:  responseContent,
		PromptTokens:     0,
		CompletionTokens: 0,
		IsStream:         textReq.Stream,
		CreatedAt:        startTime.Unix(),
		UseTime:          useTime,
		Ip:               "",
		Group:            info.UsingGroup,
	}

	// 设置 IP
	if needRecordIp {
		params.Ip = c.ClientIP()
	}

	// 设置 token 数量
	if usage != nil {
		params.PromptTokens = usage.PromptTokens
		params.CompletionTokens = usage.CompletionTokens
	}

	// 异步记录到数据库
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogError(c, fmt.Sprintf("panic in RecordConversation: %v", r))
			}
		}()
		model.RecordConversation(c, params)
	}()
}

// ExtractResponseContent 从响应中提取文本内容
func ExtractResponseContent(response *dto.OpenAITextResponse) string {
	if response == nil || len(response.Choices) == 0 {
		return ""
	}

	var content string
	for _, choice := range response.Choices {
		if choice.Message.Content != nil {
			if str, ok := choice.Message.Content.(string); ok {
				content += str
			} else if arr, ok := choice.Message.Content.([]interface{}); ok {
				// 处理数组类型的 content
				contentBytes, _ := json.Marshal(arr)
				content += string(contentBytes)
			}
		}
		// 处理 reasoning_content（用于 o3 等模型）
		if choice.Message.ReasoningContent != "" {
			content += "\n[Reasoning]\n" + choice.Message.ReasoningContent
		}
		if choice.Message.Reasoning != "" {
			content += "\n[Reasoning]\n" + choice.Message.Reasoning
		}
	}

	return content
}

// ExtractStreamResponseContent 从流式响应中提取文本内容
// 注意：这个函数需要在流式处理过程中累积内容
type StreamContentCollector struct {
	Content string
}

func (s *StreamContentCollector) AddChunk(delta *dto.ChatCompletionsStreamResponseChoiceDelta) {
	if delta == nil {
		return
	}
	if delta.Content != nil {
		s.Content += *delta.Content
	}
	if delta.ReasoningContent != nil {
		s.Content += "\n[Reasoning]\n" + *delta.ReasoningContent
	}
	if delta.Reasoning != nil {
		s.Content += "\n[Reasoning]\n" + *delta.Reasoning
	}
}

func (s *StreamContentCollector) GetContent() string {
	return s.Content
}
