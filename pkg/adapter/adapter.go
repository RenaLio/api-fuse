package adapter

import (
	v1 "awesomeGoProject/api/v1"
	"net/http"
)

type Provider interface {
	GetId() uint64
	GetName() string
	GetURL() string
	GetAPIKey() string
	GetModels() []string
	GetActivateModels() []string
	IsModelActive(modelId string) bool
	// HasModel 传入的是真正的模型Id，实际的模型id
	HasModel(modelId string) bool
	SetModelActive(modelId string, status bool)
	FetchModels() error
	HealthCheck(modelId string) bool

	Chat(w http.ResponseWriter, r *http.Request, req *v1.ChatCompletionRequest) error
}

type Adapter interface {
	Provider
}

type adapter struct {
	Provider
}

// NewAdapter returns a new Adapter
func NewAdapter(provider Provider) Adapter {
	return &adapter{Provider: provider}
}

type (
	ModelRespInfo = v1.ModelResp
	ModelInfo     = v1.Model
)

//func CheckModelActive(provider Provider, modelId string) bool {
//	return provider.IsModelActive(modelId)
//}
