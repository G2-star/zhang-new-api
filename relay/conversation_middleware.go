package relay

import (
	"time"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
)

// 这个文件提供了一个简化的集成方案
// 通过中间件的方式在请求/响应周期中记录对话

// ConversationRecordMiddleware 对话记录中间件
// 在请求开始时记录时间和请求信息
func ConversationRecordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		c.Set("conversation_start_time", time.Now())

		// 继续处理请求
		c.Next()
	}
}

// AfterResponseHook 响应后钩子函数
// 应该在每个 handler 返回前调用此函数来记录对话
//
// 使用方式：
// defer relay.AfterResponseHook(c, info, textReq, &simpleResponse, usage, startTime)
func AfterResponseHook(c *gin.Context, info *relaycommon.RelayInfo, textReq *dto.GeneralOpenAIRequest, response *dto.OpenAITextResponse, usage *dto.Usage, startTime time.Time) {
	if response == nil {
		return
	}

	responseContent := ExtractResponseContent(response)
	RecordConversationHelper(c, info, textReq, responseContent, usage, startTime)
}

// StreamAfterResponseHook 流式响应后钩子函数
//
// 使用方式：
// collector := &relay.StreamContentCollector{}
// ... 在流处理过程中调用 collector.AddChunk(delta) ...
// defer relay.StreamAfterResponseHook(c, info, textReq, collector, usage, startTime)
func StreamAfterResponseHook(c *gin.Context, info *relaycommon.RelayInfo, textReq *dto.GeneralOpenAIRequest, collector *StreamContentCollector, usage *dto.Usage, startTime time.Time) {
	if collector == nil {
		return
	}

	responseContent := collector.GetContent()
	RecordConversationHelper(c, info, textReq, responseContent, usage, startTime)
}
