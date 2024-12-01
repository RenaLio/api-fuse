package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 处理请求的函数
func handleCompletion(w http.ResponseWriter, r *http.Request) {
	var req PromptRequest

	// 解析请求体中的JSON数据
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Stream {
		handleStream(w, r, &req)
	} else {
		handleNoStream(w, r, &req)
	}
}

func handleStream(w http.ResponseWriter, r *http.Request, req *PromptRequest) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	for _ = range 20 {
		resp := generatePromptResponse(req)
		jsonStr, _ := json.Marshal(resp)
		fmt.Fprintf(w, "data: %s\n\n", jsonStr)
		w.(http.Flusher).Flush()
	}
	fmt.Fprintf(w, "data: [DONE]\n\n")
}

func handleNoStream(w http.ResponseWriter, r *http.Request, req *PromptRequest) {
	resp := generatePromptResponse(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func generatePromptResponse(req *PromptRequest) *PromptResponse {
	response := &PromptResponse{
		ID:      "cmpl-1234567890",
		Object:  "text_completion",
		Created: 1622816795,
		Model:   req.Model,
		Choices: []struct {
			Text         string `json:"text"`
			Index        int    `json:"index"`
			Logprobs     any    `json:"logprobs"`
			FinishReason string `json:"finish_reason"`
		}{
			{
				Text:         "This is a simulated response from the OpenAI API.",
				Index:        0,
				Logprobs:     nil,
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     5,
			CompletionTokens: 10,
			TotalTokens:      15,
		},
	}
	return response
}
