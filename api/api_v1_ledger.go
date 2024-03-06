// region: packages

package main

import (
	"models"
	"utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// endregion: packages
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
