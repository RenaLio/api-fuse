package main

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

type ChatCompletionNoStreamResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatCompletionStreamResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []ChoiceWithDelta `json:"choices"`
	Usage   Usage             `json:"usage,omitempty"`
}

type ChoiceWithDelta struct {
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
	Delta        Delta  `json:"delta"`
}

type Delta struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}
