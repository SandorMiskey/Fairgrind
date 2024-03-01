// region: packages

package main

import (
	"models"
	"utils"

	"github.com/gofiber/fiber/v2"
)

// endregion: packages
// region: v1_batches_types_get

// @Summary		batch types
// @Description	get list of possible batch types
// @Tags			/batches
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.ApiResponse{data=[]models.ClearingBatchType}
// @Failure		400	{object}	models.ApiResponse{data=nil}
// @Failure		500	{object}	models.ApiResponse{data=nil}
// @Router			/batches/types [get]
func v1_batches_types_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var types []models.ClearingBatchType
	result := DB.Find(&types)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching batch types", result.Error.Error())
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta.Rows = len(types)
	response.Data = types
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_batches_statuses_get

// @Summary		batch statuses
// @Description	get list of possible batch statuses
// @Tags			/batches
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.ApiResponse{data=[]models.ClearingBatchType}
// @Failure		400	{object}	models.ApiResponse{data=nil}
// @Failure		500	{object}	models.ApiResponse{data=nil}
// @Router			/batches/statuses [get]
func v1_batches_statuses_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var statuses []models.ClearingBatchStatus
	result := DB.Find(&statuses)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching batch statuses", result.Error.Error())
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
