package main

import (
	"awesomeGoProject/pkg/adapter"
	"awesomeGoProject/pkg/adapter/openai"
	"awesomeGoProject/pkg/adapter/siliconflow"
)

func NewAdapter(typeName string, id uint64, name, url, apiKey string, models []string) adapter.Adapter {
	var provider adapter.Provider
	switch typeName {
	case "openai":
		provider = openai.NewClient(id, name, url, apiKey, models)
	case "openaiCompatible":
		provider = openai.NewClient(id, name, url, apiKey, models)
	case "openaiCompatibleNoFetchModel":
		c := openai.NewClient(id, name, url, apiKey, models)
		client := c.(*openai.Client)
		provider = openai.NewOpenaiCompatibleNoFetchModelClient(client)
	case "siliconflow":
		provider = siliconflow.NewAdapter(id, name, url, apiKey, models)
	default:
		provider = openai.NewClient(id, name, url, apiKey, models)
	}
	return adapter.NewAdapter(provider)
}
