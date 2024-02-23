// region: packages

package main

import (
	"models"
	"utils"

	"github.com/gofiber/fiber/v2"
)

// endregion: packages
// region: globals

const (
	CLEARING_WALLETS_DETAILED_VIEW string = "clearing_wallets_detailed_view"
	CLEARING_WALLETS_SUMMED_VIEW   string = "clearing_wallets_summed_view"
)

// endregion
// region: v1_wallets_summed_get

// @Summary		get summed wallet of user
// @Description	get wallet of user summed by token type and status
// @Tags			/wallets
// @Accept			json
// @Produce		json
// @Param			user_id	query		int	true	"user/grinder id"
// @Success		200		{object}	models.ApiResponse{data=[]models.ClearingWalletsSummedView}
// @Failure		400		{object}	models.ApiResponse{}
// @Failure		500		{object}	models.ApiResponse{}
// @Router			/wallet/summed [get]
func v1_wallets_summed_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	user_id := uint(c.QueryInt("user_id", 0))
	if user_id < 1 {
		response.Message = "invalid user_id"
		return c.Status(400).JSON(response)
	}

	// endregion
	// region: data

	var wallet []models.ClearingWalletsSummedView
	result := DB.Table(CLEARING_WALLETS_SUMMED_VIEW).Where(&models.ClearingWalletsSummedView{ClearingLedgerUserId: user_id}).Find(&wallet)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching wallet", result.Error.Error())
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta["rows"] = len(wallet)
	response.Data = wallet
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_wallets_detailed_get

// @Summary		get detailed wallet of user
// @Description	get wallet of user summed by project, token type and status
// @Tags			/wallets
// @Accept			json
// @Produce		json
// @Param			user_id	query		int	true	"user/grinder id"
// @Success		200		{object}	models.ApiResponse{data=[]models.ClearingWalletsDetailedView}
// @Failure		400		{object}	models.ApiResponse{}
// @Failure		500		{object}	models.ApiResponse{}
// @Router			/wallets/detailed [get]
func v1_wallets_detailed_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	user_id := uint(c.QueryInt("user_id", 0))
	if user_id < 1 {
		response.Message = "invalid user_id"
		return c.Status(400).JSON(response)
	}

	// endregion
	// region: data

	var wallet []models.ClearingWalletsDetailedView
	result := DB.Table(CLEARING_WALLETS_DETAILED_VIEW).Where(&models.ClearingWalletsDetailedView{ClearingLedgerUserId: user_id}).Find(&wallet)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching wallet", result.Error.Error())
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta["rows"] = len(wallet)
	response.Data = wallet
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
