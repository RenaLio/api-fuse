package main

type PromptRequest struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	Stream           bool     `json:"stream"`
	Echo             bool     `json:"echo"`
	FreQuencyPenalty float64  `json:"frequency_penalty"`
	Logprobs         int      `json:"logprobs"`
	PresencePenalty  float64  `json:"presence_penalty"`
	MaxTokens        int      `json:"max_tokens"`
	Stop             []string `json:"stop"`
	StreamOptions    struct {
		IncludeUsage bool `json:"include_usage"`
	}
	Suffix      string  `json:"suffix"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	User        string  `json:"user"`
}

type PromptResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// 模拟OpenAI的响应结构
type OpenAIResponse = PromptResponse
