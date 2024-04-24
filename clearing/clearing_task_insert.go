// packages {{{

package main

import (
	"models"
	// "utils"
)

/// }}}

func ClearTask(id uint) {

	// tx {{{

	tx := Db.Begin()
	if err := tx.Error; err != nil {
		Logger(LOG_ERR, err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			Logger(LOG_ERR, r)
		}
	}()

	// }}}
	// fetch task {{{

	task := models.ClearingTask{}
	qt := tx.
		Joins("ClearingBatch").
		Joins("ClearingBatch.ClearingBatchType").
		Joins("ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus").
		Joins("ClearingTaskStatus").
		Joins("ClearingTaskType").
		Find(&task, id)
	if qt.Error != nil {
		tx.Rollback()
		Logger(LOG_ERR, qt.Error)
		return
	}
	Logger(LOG_INFO, MSG_TASK_PROCESSING, task.ID)
	Logger(LOG_DEBUG, MSG_TASK_UNCLEARED, task)
	// Logger(LOG_DEBUG, MSG_TASK_UNCLEARED, JsonPP(task))

	// }}}
	// checks {{{

	// __ClearingBarchStatus and ClearingBatchType__
	// According to the DDL, clearing_batches.clearing_batch_type_id and
	// clearing_batches.clearing_batch_status_id cannot be NULL, and these
	// fields also act as foreign keys pointing to the appropriate tables,
	// therefore, it is not necessary to validate these.

	// __Multiplier__
	// task.ClearingBatch.ClearingBatchType.Multiplier is "FLOAT NOT NULL DEFAULT 1"
	// in the db, so it should be safe to assume it is always set. It still
	// can be less than 0, but this option is left open for the sake of
	// flexibility, and the frontend should, hopefully, take care of this
	// anyway.
	multiplier := task.ClearingBatch.ClearingBatchType.Multiplier

	// __ClearingLedgerStatus__
	// task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus might be
	// nil, therefore it needs to be checked, as without this, it is undefined
	// what status the record would be entered into the ledger with.
	if task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id == 0 {
		tx.Rollback()
		Logger(LOG_INFO, MSG_TASK_NONCREDITABLE, task.ID)
		Logger(LOG_DEBUG, MSG_TASK_NONCREDITABLE, "task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus", JsonPP(task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus))
		return
	}
	status := task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id

	// TODO: check task.ClearingTaskStatus.Payable and task.ClearingTaskStatus.ParentPayable

	Logger(LOG_DEBUG, task.ID, "task.ClearingBatch.ClearingBatchType.Multiplier", multiplier)
	Logger(LOG_DEBUG, task.ID, "task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id", status)

	// }}}
	// grinder fees per project and task type {{{

	fees := []models.ClearingTaskFee{}
	qf := tx.Find(&fees, "project_id = ? AND clearing_task_type_id = ?", task.ClearingBatch.ProjectId, task.ClearingTaskTypeId)
	if qf.Error != nil {
		tx.Rollback()
		Logger(LOG_ERR, qf.Error)
		return
	}
	if len(fees) == 0 {
		tx.Rollback()
		Logger(LOG_INFO, MSG_TASK_NONCREDITABLE, task.ID)
		Logger(LOG_DEBUG, MSG_TASK_NONCREDITABLE, "ClearingTaskFee", fees)
		return
	}
	Logger(LOG_DEBUG, task.ID, JsonPP(fees))

	// }}}

	// TODO: update cleared_at
	// TODO: update ledger & task
	// TODO: process parant task

	// commit {{{

	if err := tx.Commit().Error; err != nil {
		Logger(LOG_ERR, err)
	}

	// }}}

}
