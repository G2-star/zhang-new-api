# 对话记录功能集成指南

本文件说明如何在现有的 relay handler 中集成对话记录功能。

## 集成步骤

### 1. 在非流式 Handler 中集成

在 `relay/channel/openai/relay-openai.go` 的 `OpenaiHandler` 函数中，需要在返回响应之前记录对话：

```go
func OpenaiHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError) {
    startTime := time.Now() // 添加：记录开始时间
    defer service.CloseResponseBodyGracefully(resp)

    // ... 现有的响应解析代码 ...

    // 解析响应
    err = common.Unmarshal(responseBody, &simpleResponse)
    if err != nil {
        return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
    }

    // ... 现有的处理逻辑 ...

    // 在写入客户端响应之后，记录对话内容
    service.IOCopyBytesGracefully(c, resp, responseBody)

    // 添加：提取请求信息并记录对话
    if textReq, ok := info.Request.(*dto.GeneralOpenAIRequest); ok {
        responseContent := relay.ExtractResponseContent(&simpleResponse)
        relay.RecordConversationHelper(c, info, textReq, responseContent, &simpleResponse.Usage, startTime)
    }

    return &simpleResponse.Usage, nil
}
```

### 2. 在流式 Handler 中集成

在 `relay/channel/openai/relay-openai.go` 的 `OpenaiStreamHandler` 函数中，需要在流式传输过程中收集内容：

```go
func OpenaiStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError) {
    startTime := time.Now() // 添加：记录开始时间

    // 添加：创建内容收集器
    var contentCollector *relay.StreamContentCollector
    var textReq *dto.GeneralOpenAIRequest
    if common.ConversationLogEnabled {
        contentCollector = &relay.StreamContentCollector{}
        if req, ok := info.Request.(*dto.GeneralOpenAIRequest); ok {
            textReq = req
        }
    }

    // ... 现有的流式处理代码 ...

    scanner := bufio.NewScanner(resp.Body)
    // ...

    for scanner.Scan() {
        data := scanner.Bytes()
        // ... 解析 data ...

        var streamResponse dto.ChatCompletionsStreamResponse
        err = json.Unmarshal(data, &streamResponse)

        // 添加：收集流式内容
        if contentCollector != nil && len(streamResponse.Choices) > 0 {
            contentCollector.AddChunk(&streamResponse.Choices[0].Delta)
        }

        // ... 发送给客户端 ...
    }

    // 添加：流式处理完成后记录对话
    if contentCollector != nil && textReq != nil {
        responseContent := contentCollector.GetContent()
        relay.RecordConversationHelper(c, info, textReq, responseContent, usage, startTime)
    }

    return usage, nil
}
```

### 3. 在 compatible_handler.go 中集成（统一入口）

更推荐的方式是在 `relay/compatible_handler.go` 的 `TextHelper` 函数中统一集成：

```go
// 在 relay/compatible_handler.go 的 DoRequest 函数中
func DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (*http.Response, *types.NewAPIError) {
    // 保存开始时间到 context
    c.Set("request_start_time", time.Now())

    // 保存原始请求
    if textReq, ok := info.Request.(*dto.GeneralOpenAIRequest); ok {
        c.Set("original_request", textReq)
    }

    // ... 现有的请求发送代码 ...
}

// 在响应处理完成后调用记录函数
```

## 需要修改的文件列表

### 核心文件（必须修改）：
1. `relay/channel/openai/relay-openai.go`
   - `OpenaiHandler` 函数
   - `OpenaiStreamHandler` 函数

### 可选文件（如需支持其他模型格式）：
2. `relay/channel/claude/main.go`
   - `ClaudeHandler` 函数
   - `ClaudeStreamHandler` 函数

3. `relay/channel/gemini/main.go`
   - `GeminiHandler` 函数
   - `GeminiStreamHandler` 函数

### 通用入口（推荐统一处理）：
4. `relay/compatible_handler.go`
   - 在 `TextHelper` 或相关函数中统一处理

## 注意事项

1. **性能考虑**：使用 goroutine 异步记录，不阻塞主请求流程
2. **错误处理**：记录失败不应影响正常的 API 响应
3. **隐私保护**：遵循用户设置中的 IP 记录选项
4. **数据大小**：对于超长对话，考虑截断或压缩存储
5. **流式处理**：需要在整个流完成后再记录，避免不完整的对话

## 测试建议

1. 测试非流式对话记录
2. 测试流式对话记录
3. 测试开关功能（启用/禁用）
4. 测试多种模型格式
5. 测试大规模并发下的性能影响
6. 测试数据库存储和查询性能

## 环境变量

在 `.env` 文件中可以添加：
```
# 对话记录功能开关（默认关闭）
CONVERSATION_LOG_ENABLED=false

# 对话内容最大长度（字符数，超过会被截断）
CONVERSATION_MAX_LENGTH=50000
```
