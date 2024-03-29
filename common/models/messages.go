package models

import (
	"encoding/json"
)

// region: api {{{

type ApiRequest struct {
	BodyRaw    string            `json:"body_raw,omitempty"`
	Queries    map[string]string `json:"queries,omitempty"`
	ReqHeaders interface{}       `json:"req_headers,omitempty"`
}

type ApiResponseMeta struct {
	Count int64 `json:"count,omitempty"`
	Rows  int   `json:"rows,omitempty"`
}

type ApiResponse struct {
	Data    interface{}     `json:"data"`
	Meta    ApiResponseMeta `json:"meta,omitempty"`
	Request ApiRequest      `json:"request"`
	Success bool            `json:"success"`
}

// endregion }}}
// region: mq {{{

type MqMsg struct {
	Database string          `json:"database"`
	Table    string          `json:"table"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
}

// endregion }}}
// region: v1 {{{

type V1LedgerWithdrawPost struct {
	Amount          float64 `json:"amount"`
	ClearingTokenId uint    `json:"clearing_token_id"`
	Reference       string  `json:"reference,omitempty"`
	UserId          uint    `json:"user_id"`
}

// endregion }}}
