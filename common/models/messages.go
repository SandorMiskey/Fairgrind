package models

import (
	"encoding/json"
)

// region: api

type ApiResponseMeta map[string]interface{}
type ApiResponse struct {
	Data    interface{}     `json:"data"`
	Message string          `json:"message"`
	Meta    ApiResponseMeta `json:"meta"`
	Success bool            `json:"success"`
}

// endregion
// region: mq

type MqMsg struct {
	Database string          `json:"database"`
	Table    string          `json:"table"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
}

// endregion
