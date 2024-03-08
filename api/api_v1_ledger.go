// region: packages

package main

import (
	"models"
	"utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// endregion: packages
// region: globals

const (
	CLEARING_LEDGER_STATUS_PENDING  uint = 1
	CLEARING_LEDGER_STATUS_WITHDRAW uint = 2
	CLEARING_LEDGER_LABEL_WITHDRAW  uint = 7
)

// endregion
// region: v1_ledger_labels_get

// @Summary		ledger entry labels
// @Description	get list of possible ledger entry labels
// @Tags			/ledger
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.ApiResponse{data=[]models.ClearingLedgerLabel}
// @Failure		400	{object}	models.ApiResponse{data=nil}
// @Failure		500	{object}	models.ApiResponse{data=nil}
// @Router			/ledger/labels [get]
func v1_ledger_labels_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var labels []models.ClearingLedgerLabel
	result := DB.Find(&labels)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching ledger statuses", result.Error.Error())
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta.Rows = len(labels)
	response.Data = labels
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_ledger_statuses_get

// @Summary		ledger entry statuses
// @Description	get list of possible token ledger statuses
// @Tags			/ledger
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.ApiResponse{data=[]models.ClearingLedgerStatus}
// @Failure		400	{object}	models.ApiResponse{data=nil}
// @Failure		500	{object}	models.ApiResponse{data=nil}
// @Router			/ledger/statuses [get]
func v1_ledger_statuses_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var statuses []models.ClearingLedgerStatus
	result := DB.Find(&statuses)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching ledger statuses", result.Error.Error())
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta.Rows = len(statuses)
	response.Data = statuses
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_ledger_get

// @Summary		get ledger
// @Description	get filtered list of ledger entries
// @Tags			/ledger
// @Accept			json
// @Produce		json
// @Param			id							query		int		false	"id"
// @Param			user_id						query		int		false	"user/grinder id"
// @Param			clearing_task_id			query		int		false	"task id"
// @Param			clearing_ledger_status_id	query		int		false	"ledger entry status id"
// @Param			clearing_ledger_label_id	query		int		false	"ledger entry label id"
// @Param			clearing_token_id			query		int		false	"token id"
// @Param			orm_order_by				query		string	false	"order by <param> <direction>, as in 'clearing_ledger_status_id asc, clearing_ledger_label_id desc'"
// @Param			orm_page					query		int		false	"which page"
// @Param			orm_limit					query		int		false	"page size (aka # of results)"
// @Success		200							{object}	models.ApiResponse{data=[]models.ClearingLedger}
// @Failure		400							{object}	models.ApiResponse{data=nil}
// @Failure		500							{object}	models.ApiResponse{data=nil}
// @Router			/ledger [get]
func v1_ledger_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	filters := &models.ClearingLedger{
		ClearingLedgerLabelId:  uint(c.QueryInt("clearing_ledger_label_id", 0)),
		ClearingLedgerStatusId: uint(c.QueryInt("clearing_ledger_status_id", 0)),
		ClearingTaskId:         uint(c.QueryInt("clearing_task_id", 0)),
		ClearingTokenId:        uint(c.QueryInt("clearing_token_id", 0)),
		GORM: models.GORM{
			ID: uint(c.QueryInt("id", 0)),
		},
		UserId: uint(c.QueryInt("user_id", 0)),
	}

	order := c.Query("orm_order_by")
	page := c.QueryInt("orm_page", 0)
	limit := c.QueryInt("orm_limit", 0)

	// endregion
	// region: data

	var ledger []models.ClearingLedger
	var count int64
	var result *gorm.DB

	result = DB.Model(&models.ClearingLedger{}).Where(&filters).Count(&count)
	if result.Error != nil {
		Logger(LOG_ERR, "error while counting ledger entries", result.Error.Error())
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	result = DB.Scopes(utils.Paginate(page, limit)).Where(&filters).Order(order).Find(&ledger)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching ledger entries", result.Error.Error())
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta.Rows = len(ledger)
	response.Meta.Count = count
	response.Data = ledger
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_ledger_delete

// @Summary		withdraw
// @Description	withdraw from grinder's wallet
// @Tags			/ledger
// @Accept			json
// @Produce		json
// @Param			amount				query		number	true	"amount to withdraw"
// @Param			clearing_token_id	query		int		true	"token id"
// @Param			user_id				query		int		true	"user/grinder id"
// @Param			reference			query		string	false	"optional reference"
// @Success		200					{object}	models.ApiResponse{data=[]models.ClearingLedger}
// @Failure		400					{object}	models.ApiResponse{data=nil}
// @Failure		500					{object}	models.ApiResponse{data=nil}
// @Router			/ledger [delete]
func v1_ledger_delete(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	withdraw := models.ClearingLedger{
		Amount:                 c.QueryFloat("amount", 0) * -1,
		ClearingLedgerLabelId:  CLEARING_LEDGER_LABEL_WITHDRAW,
		ClearingLedgerStatusId: CLEARING_LEDGER_STATUS_WITHDRAW,
		ClearingTokenId:        uint(c.QueryInt("clearing_token_id", 0)),
		Reference:              c.Query("reference"),
		UserId:                 uint(c.QueryInt("user_id", 0)),
	}

	if withdraw.Amount > 0 {
		response.Data = "amount cannot be less than 0"
		return c.Status(400).JSON(response)
	}
	if len(withdraw.Reference) > 1024 {
		response.Data = "reference cannot be longer than 1024 bytes"
		return c.Status(400).JSON(response)
	}

	// endregion
	// region: data

	// TODO: IN TX, check if new balance is under 0

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		response.Data = err.Error()
		return c.Status(500).JSON(response)
	}

	result := tx.Create(&withdraw)
	if result.Error != nil {
		tx.Rollback()
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: availability

	var available models.ClearingWalletsSummedView
	filter := models.ClearingWalletsSummedView{
		ClearingLedgerUserId:             withdraw.UserId,
		ClearingTokenId:                  withdraw.ClearingTokenId,
		ClearingLedgerStatusWithdrawable: true,
	}

	availability := tx.Table(CLEARING_WALLETS_SUMMED_VIEW).Where(&filter).Find(&available)
	if availability.Error != nil {
		Logger(LOG_ERR, "error while fetching wallet", availability.Error.Error())
		response.Data = availability.Error.Error()
		return c.Status(500).JSON(response)
	}
	if available.ClearingLedgerAmountSum < 0 {
		response.Data = "no sufficient funds"
		return c.Status(400).JSON(response)
	}

	// endregion
	// region: response

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Data = err.Error()
		return c.Status(500).JSON(response)
	}

	response.Meta.Rows = 1
	response.Data = withdraw
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
