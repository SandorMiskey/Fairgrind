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
	result := tx.
		Joins("ClearingBatch").
		Joins("ClearingBatch.ClearingBatchType").
		Joins("ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus").
		Joins("ClearingTaskStatus").
		Joins("ClearingTaskType").
		Find(&tasks, "cleared_at IS NULL")
	if result.Error != nil {
		tx.Rollback()
		Logger(LOG_ERR, result.Error)
		return
	}
	// Logger(LOG_DEBUG, JsonPP(tasks))

	// }}}
	// tasks loop {{{

	for _, task := range tasks {
		tx.SavePoint("sp")
		Logger(LOG_INFO, MSG_TASK_UNCLEARED, task.ID)
		Logger(LOG_DEBUG, JsonPP(task))

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

		// __ClearingLedgerStatus__
		// task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus might be
		// nil, therefore it needs to be checked, as without this, it is undefined
		// what status the record would be entered into the ledger with.
		if task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus.Id == 0 {
			Logger(LOG_INFO, MSG_TASK_NONCREDITABLE, task.ID)
			Logger(LOG_DEBUG, JsonPP(task.ClearingBatch.ClearingBatchStatus.ClearingLedgerStatus))
			tx.RollbackTo("sp")
			continue
		}

		// TODO: check/get fees

		// }}}

		// TODO: get fees

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
