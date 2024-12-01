package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

func HandleChat(w http.ResponseWriter, r *http.Request) {
	var req ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 检查是否需要流式返回
	if req.Stream {
		handleStreamResponse(w, &req)
	} else {
		handleNonStreamResponse(w, &req)
	}
}

const responseText = "This is a simulated response from the OpenAI Chat API."

func handleNonStreamResponse(w http.ResponseWriter, req *ChatCompletionRequest) {
	resp := generateNoStreamResponse(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleStreamResponse(w http.ResponseWriter, req *ChatCompletionRequest) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	req.N = 18
	for i := 0; i < req.N; i++ {
		resp := generateStreamResponse(req)
		resp.Choices[0].Index = i
		jsonResp, _ := json.Marshal(resp)
		fmt.Fprintf(w, "data: %s\n\n", jsonResp)
		w.(http.Flusher).Flush()
		time.Sleep(500 * time.Millisecond) // 模拟生成时间
	}
	fmt.Fprintf(w, "data: [DONE]\n\n")
}

func generateNoStreamResponse(req *ChatCompletionRequest) *ChatCompletionNoStreamResponse {
	return &ChatCompletionNoStreamResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", rand.Intn(1000000)),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []Choice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: responseText,
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}
}

func generateStreamResponse(req *ChatCompletionRequest) *ChatCompletionStreamResponse {
	return &ChatCompletionStreamResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", rand.Intn(1000000)),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []ChoiceWithDelta{
			{
				Index: 0,
				Delta: Delta{
					Role:    "assistant",
					Content: responseText,
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 10,
			TotalTokens:      20,
		},
	}
}
