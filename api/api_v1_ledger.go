// region: packages

package main

import (
	"models"
	"utils"

	"github.com/gofiber/fiber/v2"
)

// endregion: packages
// region: v1_ledger_labels_get

//	@Summary		ledger entry labels
//	@Description	get list of possible ledger entry labels
//	@Tags			/ledger
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.ApiResponse{data=[]models.ClearingLedgerLabel}
//	@Failure		400	{object}	models.ApiResponse{data=nil}
//	@Failure		500	{object}	models.ApiResponse{data=nil}
//	@Router			/ledger/labels [get]
func v1_ledger_labels_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var labels []models.ClearingLedgerLabel
	result := DB.Find(&labels)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching ledger statuses", result.Error.Error())
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta["rows"] = len(labels)
	response.Data = labels
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_ledger_statuses_get

//	@Summary		ledger entry statuses
//	@Description	get list of possible token ledger statuses
//	@Tags			/ledger
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.ApiResponse{data=[]models.ClearingLedgerStatus}
//	@Failure		400	{object}	models.ApiResponse{data=nil}
//	@Failure		500	{object}	models.ApiResponse{data=nil}
//	@Router			/ledger/statuses [get]
func v1_ledger_statuses_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var statuses []models.ClearingLedgerStatus
	result := DB.Find(&statuses)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching ledger statuses", result.Error.Error())
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta["rows"] = len(statuses)
	response.Data = statuses
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
