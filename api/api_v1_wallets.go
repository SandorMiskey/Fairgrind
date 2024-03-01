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
		response.Data = "invalid user_id"
		return c.Status(400).JSON(response)
	}

	// endregion
	// region: data

	var wallet []models.ClearingWalletsSummedView
	result := DB.Table(CLEARING_WALLETS_SUMMED_VIEW).Where(&models.ClearingWalletsSummedView{ClearingLedgerUserId: user_id}).Find(&wallet)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching wallet", result.Error.Error())
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta.Rows = len(wallet)
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
// @Param			clearing_ledger_user_id				query		int		false	"user/grinder id"
// @Param			clearing_ledger_status_withdrawable	query		bool	false	"is withdrawable?"
// @Param			clearing_token_symbol				query		string	false	"token symbol"
// @Param			clearing_ledger_label_id			query		int		false	"ledger label id"
// @Param			clearing_ledger_label_label			query		string	false	"ledger label"
// @Param			project_id							query		int		false	"project id"
// @Param			project_name						query		string	false	"project name"
// @Param			orm_order_by						query		string	false	"order by field asc/desc, fields are: filter fields plus clearing_ledger_amount_sum, clearing_ledger_updated_at_max and clearing_tasks_id_count"
// @Param			orm_page							query		int		false	"which page"
// @Param			orm_size							query		int		false	"page size (aka # of results)"
// @Success		200									{object}	models.ApiResponse{data=[]models.ClearingWalletsDetailedView}
// @Failure		400									{object}	models.ApiResponse{}
// @Failure		500									{object}	models.ApiResponse{}
// @Router			/wallets/detailed [get]
func v1_wallets_detailed_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	filters := &models.ClearingWalletsDetailedView{
		ClearingLedgerUserId:             uint(c.QueryInt("clearing_ledger_user_id", 0)),
		ClearingLedgerStatusWithdrawable: c.QueryBool("clearing_ledger_status_withdrawable"),
		ClearingTokenSymbol:              c.Query("clearing_token_symbol"),
		ClearingLedgerLabelId:            uint(c.QueryInt("clearing_ledger_label_id", 0)),
		ClearingLedgerLabelLabel:         c.Query("clearing_ledger_label_label"),
		ProjectId:                        uint(c.QueryInt("project_id", 0)),
		ProjectName:                      c.Query("project_name"),
	}

	order := c.Query("orm_order_by")
	page := c.QueryInt("orm_page", 0)
	limit := c.QueryInt("orm_limit", 0)

	// endregion
	// region: data

	var wallet []models.ClearingWalletsDetailedView
	var count int64
	var result *gorm.DB

	result = DB.Table(CLEARING_WALLETS_DETAILED_VIEW).Where(&filters).Count(&count)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching wallet", result.Error.Error())
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	result = DB.Table(CLEARING_WALLETS_DETAILED_VIEW).Scopes(utils.Paginate(page, limit)).Where(&filters).Order(order).Find(&wallet)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching wallet", result.Error.Error())
		response.Data = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta.Rows = len(wallet)
	response.Meta.Count = count
	response.Data = wallet
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
