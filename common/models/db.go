package models

import (
	"time"

	"gorm.io/gorm"
)

type GORM struct {
	ID        uint           `json:"id" form:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// region: batches

type ClearingBatchType struct {
	GORM
	Label       string  `json:"label"`
	Description string  `json:"description"`
	Multiplier  float64 `json:"multiplier"`
}

type ClearingBatchStatus struct {
	GORM
	Id                     uint                 `json:"id" gorm:"type:SMALLINT UNSIGNED;"`
	Label                  string               `json:"label"`
	Description            string               `json:"description"`
	ClearingLedgerStatusId uint                 `json:"-"`
	ClearingLedgerStatus   ClearingLedgerStatus `json:"clearing_ledger_status"`
}

type ClearingBatch struct {
	GORM
	ClearingBatchTypeId   uint                `json:"-"`
	ClearingBatchType     ClearingBatchType   `json:"clearing_batch_type"`
	ClearingBatchStatusId uint                `json:"-" gorm:"type:SMALLINT UNSIGNED;"`
	ClearingBatchStatus   ClearingBatchStatus `json:"clearing_batch_status"`
	Label                 string              `json:"label"`
	ProjectId             uint                `json:"project_id" gorm:"type:int(11);"`
	Description           string              `json:"description"`
}

// endregion
// region: ledger

type ClearingLedgerStatus struct {
	GORM
	Id           uint   `json:"id" gorm:"type:SMALLINT UNSIGNED;"`
	Label        string `json:"label"`
	Description  string `json:"description"`
	Withdrawable bool   `json:"withdrawable"`
}

type ClearingLedgerLabel struct {
	GORM
	Id          uint   `json:"id" gorm:"type:SMALLINT UNSIGNED;"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// endregion
// region: task

type ClearingTaskStatus struct {
	GORM
	Id            uint   `json:"id" gorm:"type:SMALLINT UNSIGNED;"`
	Label         string `json:"label"`
	Description   string `json:"description"`
	Payable       bool   `json:"payable"`
	ParentPayable bool   `json:"parent_payable"`
}

type ClearingTaskType struct {
	GORM
	Id             uint   `json:"id" gorm:"type:SMALLINT UNSIGNED;"`
	Label          string `json:"label"`
	Description    string `json:"description"`
	TaskPayable    bool   `json:"task_payable"`
	SubtaskPayable bool   `json:"subtask_payable"`
}

type ClearingTaskFee struct {
	UserId             uint    `json:"user_id" gorm:"type:INT(11); primaryKey; autoIncrement:false"`
	ProjectId          uint    `json:"project_id" gorm:"type:INT(11); primaryKey; autoIncrement:false"`
	ClearingTaskTypeId uint    `json:"clearing_task_type_id" gorm:"type:SMALLINT UNSIGNED; primaryKey; autoIncrement:false"`
	TaskFee            float64 `json:"task_fee"`
	SubtaskFee         float64 `json:"subtask_fee"`
	ClearingTokenId    uint    `json:"clearing_token_id" gorm:"type:MEDIUMINT UNSIGNED; primaryKey; autoIncrement:false"`
}

type ClearingTask struct {
	GORM
	ClearingBatchId      uint    `json:"clearing_batch_id"`
	ClearingTaskId       uint    `json:"clearing_task_id,omitempty"`
	ClearingTaskTypeId   uint    `json:"clearing_task_type_id" gorm:"type:SMALLINT UNSIGNED;"`
	ClearingTaskStatusId uint    `json:"clearing_task_status_id" gorm:"type:SMALLINT UNSIGNED;"`
	UserId               uint    `json:"user_id" gorm:"type:INT(11);"`
	Input                *string `json:"input"`
	Output               string  `json:"output"`
}

// endregion
// region: token

type ClearingTokenType struct {
	GORM
	Id          uint   `json:"id" gorm:"type:MEDIUMINT UNSIGNED;"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type ClearingToken struct {
	GORM
	Id                  uint              `json:"id" gorm:"type:MEDIUMINT UNSIGNED;"`
	Label               string            `json:"label"`
	Symbol              string            `json:"symbol"`
	Description         string            `json:"description"`
	ClearingTokenTypeId uint              `json:"-" gorm:"type:MEDIUMINT UNSIGNED;"`
	ClearingTokenType   ClearingTokenType `json:"clearing_token_type"`
}

// endregion
// region: wallet

type ClearingWalletsSummedView struct {
	ClearingLedgerAmountSum          float64   `json:"amount"`
	ClearingLedgerUpdatedAtMax       time.Time `json:"updated_at"`
	ClearingLedgerUserId             uint      `json:"-" form:"user_id"`
	ClearingLedgerStatusWithdrawable bool      `json:"withdrawable"`
	ClearingTokenSymbol              string    `json:"token_symbol"`
}

type ClearingWalletsDetailedView struct {
	ClearingLedgerAmountSum          float64   `json:"clearing_ledger_amount_sum"`
	ClearingLedgerUpdatedAtMax       time.Time `json:"clearing_ledger_updated_at_max"`
	ClearingLedgerUserId             uint      `json:"clearing_ledger_user_id"`
	ClearingLedgerStatusWithdrawable bool      `json:"clearing_ledger_status_withdrawable"`
	ClearingTokenSymbol              string    `json:"clearing_token_symbol"`
	ClearingLedgerLabelId            uint      `json:"-"`
	ClearingLedgerLabelLabel         string    `json:"clearing_ledger_label_label"`
	ClearingTasksIdCount             uint      `json:"clearing_tasks_id_count"`
	ProjectId                        uint      `json:"-"`
	ProjectName                      string    `json:"project_name"`
}

// endregion
