package v1

type ChatCompletionRequest struct {
	Model            string             `json:"model"`                       // 模型名称
	Messages         []ChatMessage      `json:"messages"`                    // 聊天消息
	Temperature      float64            `json:"temperature,omitempty"`       // 温度
	TopP             float64            `json:"top_p,omitempty"`             // top_p
	Stream           bool               `json:"stream,omitempty"`            // 是否流式返回
	Stop             []string           `json:"stop,omitempty"`              // 停止标志
	MaxTokens        int                `json:"max_tokens,omitempty"`        // 最大 token 数
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`  // 存在惩罚
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"` // 频率惩罚
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
	N                int                `json:"n,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
	Name    string `json:"name,omitempty"`
}
