package openai

import "awesomeGoProject/pkg/adapter"

type openaiCompatibleNoFetchModel struct {
	*Client
}

func (c *openaiCompatibleNoFetchModel) FetchModels() error {
	list := c.Models
	for _, model := range list {
		c.ModelStatus[model] = true
	}
	return nil
}

func NewOpenaiCompatibleNoFetchModelClient(c *Client) adapter.Provider {
	return &openaiCompatibleNoFetchModel{c}
}
