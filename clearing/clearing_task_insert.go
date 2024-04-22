// packages {{{

package main

import (
	"models"
)

/// }}}

func TaskUncleared() {

	// lock and tx {{{

	tx := Db.Begin()
	if err := tx.Error; err != nil {
		Logger(LOG_ERR, err)
		return
	}

	Lock.Lock()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			Logger(LOG_ERR, r)
		}
		Lock.Unlock()
	}()

	// }}}
	// uncleared {{{

	tasks := []models.ClearingTask{}
	qt := tx.
		Joins("ClearingBatch").
		Joins("ClearingBatch.ClearingBatchType").
		Joins("ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus").
		Joins("ClearingTaskStatus").
		Joins("ClearingTaskType").
		Find(&tasks, "cleared_at IS NULL")
	if qt.Error != nil {
		tx.Rollback()
		Logger(LOG_ERR, qt.Error)
		return
	}

	// }}}
	// tasks loop {{{

	for _, task := range tasks {
		tx.SavePoint("sp")
		Logger(LOG_INFO, MSG_TASK_UNCLEARED, task.ID)
		Logger(LOG_DEBUG, MSG_TASK_UNCLEARED, task)
		// Logger(LOG_DEBUG, MSG_TASK_UNCLEARED, JsonPP(task))

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
			tx.RollbackTo("sp")
			Logger(LOG_INFO, MSG_TASK_NONCREDITABLE, task.ID)
			Logger(LOG_DEBUG, MSG_TASK_NONCREDITABLE, "task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus", JsonPP(task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus))
			continue
		}
		status := task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id

		// TODO: Separate the reading of tasks and the processing of individual tasks, thereby aiding the processing of parent tasks
		// TODO: check task.ClearingTaskStatus.Payable and task.ClearingTaskStatus.ParentPayable

		// __Grinder fees per project and task type__
		fees := []models.ClearingTaskFee{}
		qf := tx.Find(&fees, "project_id = ? AND clearing_task_type_id = ?", task.ClearingBatch.ProjectId, task.ClearingTaskTypeId)
		if qf.Error != nil {
			tx.RollbackTo("sp")
			Logger(LOG_ERR, qf.Error)
			continue
		}
		if len(fees) == 0 {
			tx.RollbackTo("sp")
			Logger(LOG_INFO, MSG_TASK_NONCREDITABLE, task.ID)
			Logger(LOG_DEBUG, MSG_TASK_NONCREDITABLE, "ClearingTaskFee", fees)
			continue
		}

		Logger(LOG_DEBUG, task.ID, "task.ClearingBatch.ClearingBatchType.Multiplier", multiplier)
		Logger(LOG_DEBUG, task.ID, "task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id", status)
		Logger(LOG_DEBUG, task.ID, JsonPP(fees))

		// TODO: update ledger & task

	}

	// }}}
	// commit {{{

	err := tx.Commit().Error
	if err != nil {
		Logger(LOG_ERR, err)
	}

	// }}}

}
