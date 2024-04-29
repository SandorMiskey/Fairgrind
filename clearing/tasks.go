// packages {{{

package main

import (
	"encoding/json"
	"models"
	"time"
	"utils"

	"gorm.io/gorm"
)

/// }}}

func ClearTask(id uint, path []uint) error { // {{{

	// tx {{{

	var result *gorm.DB

	tx := Db.Begin()
	if err := tx.Error; err != nil {
		return err
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
	result = tx.
		Joins("ClearingBatch").
		Joins("ClearingBatch.ClearingBatchType").
		Joins("ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus").
		Joins("ClearingTaskStatus").
		Joins("ClearingTaskType").
		Find(&task, id)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	Logger(LOG_DEBUG, MSG_TASK_PROCESSING, task.ID)
	Logger(LOG_DEBUG, MSG_TASK_UNCLEARED, JsonPP(task))

	// }}}
	// checks {{{

	// __ClearingBatchStatus and ClearingBatchType__
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
		Logger(LOG_INFO, MSG_TASK_NONCREDITABLE, task.ID)
		Logger(LOG_DEBUG, MSG_TASK_NONCREDITABLE, "task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus", JsonPP(task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus))
		return commitTask(id, tx)
	}
	status := task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id

	// TODO: check task.ClearingTaskStatus.Payable and task.ClearingTaskStatus.ParentPayable

	Logger(LOG_DEBUG, task.ID, "task.ClearingBatch.ClearingBatchType.Multiplier", multiplier)
	Logger(LOG_DEBUG, task.ID, "task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id", status)

	// }}}
	// grinder fees per project and task type {{{

	fees := []models.ClearingTaskFee{}
	result = tx.Find(&fees, "project_id = ? AND clearing_task_type_id = ?", task.ClearingBatch.ProjectId, task.ClearingTaskTypeId)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if len(fees) == 0 {
		Logger(LOG_INFO, MSG_TASK_NONCREDITABLE, task.ID)
		Logger(LOG_DEBUG, MSG_TASK_NONCREDITABLE, "ClearingTaskFee", fees)
		return commitTask(id, tx)
	}
	Logger(LOG_DEBUG, task.ID, JsonPP(fees))

	// }}}
	// update ledger and task {{{

	// delete existing ledger entries
	result = tx.Where("clearing_task_id = ? AND deleted_at IS NULL", id).Delete(&models.ClearingLedger{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// insert new ledger entries
	for _, fee := range fees {
		Logger(LOG_DEBUG, task.ID, JsonPP(fee))

		// subtasks
		var output []interface{}
		var subtasks int
		err := json.Unmarshal([]byte(task.Output), &output)
		if err != nil {
			Logger(LOG_NOTICE, err)
			subtasks = 0
		} else {
			subtasks = len(output)
		}

		// fee
		var amount float64
		amount = fee.TaskFee * multiplier
		amount += fee.SubtaskFee * multiplier * float64(subtasks)

		// ledger entry
		ledger := models.ClearingLedger{
			Amount:                 amount,
			ClearingTaskId:         id,
			ClearingLedgerStatusId: task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id,
			ClearingLedgerLabelId:  DB_CLEARING_LEDGER_LABEL_TASK,
			ClearingTokenId:        fee.ClearingTokenId,
			UserId:                 task.UserId,
		}
		Logger(LOG_DEBUG, task.ID, JsonPP(ledger))

		// insert
		result := tx.Create(&ledger)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	// process parent task (also check for circular references)
	path = append(path, id)
	if task.ClearingTaskId != 0 && !utils.Contains(path, task.ClearingTaskId) {
		err := ClearTask(task.ClearingTaskId, path)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// }}}
	// commit {{{

	return commitTask(id, tx)

	// }}}

} // }}}

func commitTask(id uint, tx *gorm.DB) error { // {{{
	// update task.cleared_at
	result := tx.Model(&models.ClearingTask{}).Where("id = ?", id).Update("cleared_at", time.Now())
	if result.Error != nil {
		tx.Rollback()
		Logger(LOG_ERR, id, result.Error)
		return result.Error
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		Logger(LOG_ERR, id, err)
		return err
	}

	// all good
	Logger(LOG_INFO, id, MSG_TASK_CLEARED)
	return nil
} // }}}

// vim: foldmethod=marker
