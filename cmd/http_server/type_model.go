package main

type Model struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	Created    int64  `json:"created,omitempty"`
	OwnedBy    string `json:"owned_by,omitempty"`
	Permission any    `json:"permission,omitempty"`
	Root       string `json:"root,omitempty"`
	Parent     string `json:"parent,omitempty"`
}

type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}
