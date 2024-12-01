package siliconflow

import (
	"awesomeGoProject/pkg/adapter"
	"awesomeGoProject/pkg/adapter/openai"
)

type SiliconflowClient struct {
	*openai.Client
}

func NewAdapter(id uint64, name, url, apiKey string, models []string) adapter.Provider {
	client := openai.NewClient(id, name, url, apiKey, models)
	c := client.(*openai.Client)
	return &SiliconflowClient{
		Client: c,
	}
}
