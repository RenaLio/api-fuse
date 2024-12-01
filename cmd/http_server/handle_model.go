package main

import (
	"encoding/json"
	"net/http"
)

func handleModel(w http.ResponseWriter, r *http.Request) {
	resp := ModelsResponse{
		Object: "list",
		Data: []Model{
			{
				ID:      "gpt-3.5-turbo",
				Object:  "model",
				Created: 1709808602,
				OwnedBy: "openai",
			},
			{
				ID:      "gpt-3.5-turbo-12331231",
				Object:  "model",
				Created: 1709808602,
				OwnedBy: "openai",
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respStr, _ := json.Marshal(resp)
	w.Write(respStr)
}
