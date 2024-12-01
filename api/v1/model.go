package v1

type Model struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	Created    int64  `json:"created"`
	OwnedBy    string `json:"owned_by"`
	Permission any    `json:"permission,omitempty"`
}

type ModelResp struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}
