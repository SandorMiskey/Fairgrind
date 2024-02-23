// region: packages

package main

import (
	"models"
	"utils"

	"github.com/gofiber/fiber/v2"
)

// endregion: packages
// region: v1_tokens_get

// @Summary		token list
// @Description	get list of registered tokens w/ type
// @Tags			/tokens
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.ApiResponse{data=[]models.ClearingToken}
// @Failure		400	{object}	models.ApiResponse{}
// @Failure		500	{object}	models.ApiResponse{}
// @Router			/tokens [get]
func v1_tokens_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var tokens []models.ClearingToken
	result := DB.Preload("ClearingTokenType").Find(&tokens)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching token types", result.Error.Error())
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta["rows"] = len(tokens)
	response.Data = tokens
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_tokens_types_get

// @Summary		token types
// @Description	get list of possible token types
// @Tags			/tokens
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.ApiResponse{data=[]models.ClearingTokenType}
// @Failure		400	{object}	models.ApiResponse{data=nil}
// @Failure		500	{object}	models.ApiResponse{data=nil}
// @Router			/tokens/types [get]
func v1_tokens_types_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var types []models.ClearingTokenType
	result := DB.Find(&types)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching token types", result.Error.Error())
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta["rows"] = len(types)
	response.Data = types
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
