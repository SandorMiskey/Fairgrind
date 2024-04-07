// packages {{{

package models

// }}}
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
	Data     map[string]interface{} `json:"data"`
	Database string                 `json:"database"`
	Old      map[string]interface{} `json:"old"`
	Table    string                 `json:"table"`
	Type     string                 `json:"type"`
	Xid      uint                   `json:"xid"`
}

// endregion }}}
// region: v1 {{{

type V1LedgerWithdrawPost struct {
	Amount          float64 `json:"amount"`
	ClearingTokenId uint    `json:"clearing_token_id"`
	Reference       string  `json:"reference,omitempty"`
	UserId          uint    `json:"user_id"`
}

type V1LedgerCreditPost struct {
	Amount                 float64 `json:"amount"`
	ClearingLedgerStatusId uint    `json:"clearing_ledger_status_id"`
	ClearingLedgerLabelId  uint    `json:"clearing_ledger_label_id"`
	ClearingTaskId         uint    `json:"clearing_task_id,omitempty"`
	ClearingTokenId        uint    `json:"clearing_token_id"`
	Reference              string  `json:"reference,omitempty"`
	UserId                 uint    `json:"user_id"`
}

// endregion }}}
