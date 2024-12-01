package openai

import (
	v1 "awesomeGoProject/api/v1"
	"awesomeGoProject/pkg/adapter"
	"awesomeGoProject/pkg/tools"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type Client struct {
	Id          uint64
	Name        string
	URL         string
	ApiKey      string
	Models      []string
	ModelStatus map[string]bool
	mu          sync.RWMutex
	// httpClient  *http.Client
}

func NewClient(id uint64, name, url, apiKey string, models []string) adapter.Provider {
	return &Client{
		Id:          id,
		Name:        name,
		URL:         url,
		ApiKey:      apiKey,
		Models:      models,
		ModelStatus: make(map[string]bool),
	}
}

// GetId implements adapter.Provider.
func (c *Client) GetId() uint64 {
	return c.Id
}

// GetName implements adapter.Provider.
func (c *Client) GetName() string {
	return c.Name
}

// GetURL implements adapter.Provider.
func (c *Client) GetURL() string {
	return c.URL
}

// GetAPIKey implements adapter.Provider.
func (c *Client) GetAPIKey() string {
	return c.ApiKey
}

// GetModels implements adapter.Provider.
func (c *Client) GetModels() []string {
	result := make([]string, 0)
	for modelId := range c.ModelStatus {
		result = append(result, modelId)
	}
	return result
}

// GetActivateModels implements adapter.Provider.
func (c *Client) GetActivateModels() []string {
	result := make([]string, 0)
	for modelId, status := range c.ModelStatus {
		if status {
			result = append(result, modelId)
		}
	}
	return result
}

// IsModelActive implements adapter.Provider.
func (c *Client) IsModelActive(modelId string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ModelStatus[modelId]
}

func (c *Client) HasModel(modelId string) bool {
	_, ok := c.ModelStatus[modelId]
	return ok
}

// SetModelActive implements adapter.Provider.
func (c *Client) SetModelActive(modelId string, status bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ModelStatus[modelId] = status
}

// FetchModels implements adapter.Provider.
func (c *Client) FetchModels() error {
	// 请求接口
	// -> 失败 -> 默认写的models
	// -> 成功 -> Set集合
	rowUrl := fmt.Sprintf("%s/v1/models", c.URL)
	u, err := url.Parse(rowUrl)
	if err != nil {
		return err
	}
	url := u.String()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()
	// TODO 定时获取model更新
	if resp.StatusCode != http.StatusOK {
		c.mu.Lock()
		defer c.mu.Unlock()
		for _, val := range c.Models {
			c.ModelStatus[val] = true
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var modelsResponse adapter.ModelRespInfo
		err = json.Unmarshal(body, &modelsResponse)
		if err != nil {
			return err
		}
		var modelSet = make(map[string]bool)
		for _, model := range modelsResponse.Data {
			modelSet[model.ID] = true
		}
		for _, val := range c.Models {
			modelSet[val] = true
		}
		c.mu.Lock()
		for key := range modelSet {
			c.ModelStatus[key] = true
		}
		c.mu.Unlock()
	}
	return nil
}

// Chat implements adapter.Provider.
func (c *Client) Chat(w http.ResponseWriter, r *http.Request, req *v1.ChatCompletionRequest) error {
	var err error
	if req.Stream {
		err = c.handleStream(w, r, req)
	} else {
		err = c.handleNonStream(w, r, req)
	}
	if err != nil {
		c.SetModelActive(req.Model, false)
		return err
	}
	return nil
}

func (c *Client) handleStream(w http.ResponseWriter, r *http.Request, req *v1.ChatCompletionRequest) error {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	body, err := json.Marshal(req)
	if err != nil {
		return errors.New("failed to serialize request body")
	}
	resp, err := c.DoChatRequest(body)
	if err != nil {
		return errors.New("failed to send request." + err.Error())
	}
	defer resp.Close()
	scanner := bufio.NewScanner(resp)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		fmt.Fprintf(w, "%s\n\n", line)
		//n, err := w.Write([]byte(line + "\n\n"))
		//if err != nil {
		//	fmt.Println(n, err)
		//}
		w.(http.Flusher).Flush()
	}
	// 检查流读取错误
	if err := scanner.Err(); err != nil {
		return errors.New("failed to read stream response." + err.Error())
	}
	return nil
}

func (c *Client) handleNonStream(w http.ResponseWriter, r *http.Request, req *v1.ChatCompletionRequest) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	body, err := json.Marshal(req)
	resp, err := c.DoChatRequest(body)
	if err != nil {
		return errors.New("failed to send request." + err.Error())
	}
	defer resp.Close()
	_, err = io.Copy(w, resp)
	if err != nil {
		return errors.New("failed to write response." + err.Error())
	}
	return nil
}

func (c *Client) DoChatRequest(body []byte) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/v1/chat/completions", c.URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, respBody)
	}
	return resp.Body, nil
}

func (c *Client) CheckModelStatus(modelId string) bool {
	msg := tools.GetRandMsg()
	req := &v1.ChatCompletionRequest{
		Model: modelId,
		Messages: []v1.ChatMessage{
			{
				Role:    "user",
				Content: msg,
			},
		},
		MaxTokens: 12,
	}
	body, _ := json.Marshal(req)
	_, err := c.DoChatRequest(body)
	return err == nil
}

func (c *Client) HealthCheck(modelId string) bool {
	msg := tools.GetRandMsg()
	req := &v1.ChatCompletionRequest{
		Model: modelId,
		Messages: []v1.ChatMessage{
			{
				Role:    "user",
				Content: msg,
			},
		},
		MaxTokens: 12,
		Stream:    true,
	}
	body, _ := json.Marshal(req)
	_, err := c.DoChatRequest(body)
	return err == nil
}
