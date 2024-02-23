// region: packages

package main

import (
	"models"
	"utils"

	"github.com/gofiber/fiber/v2"
)

// endregion: packages
// region: v1_tasks_post

// @Summary		task register
// @Description	register task
// @Tags			/tasks
// @Accept			json
// @Produce		json
// @Param			request	body		[]models.ClearingTask	true	"json request body, omit 'id' or set to 0, otherwise value will be used, 'input' and 'output' must be valid JSON if supplied"
// @Success		200		{object}	models.ApiResponse{data=[]models.ClearingTask}
// @Failure		400		{object}	models.ApiResponse{data=nil}
// @Failure		500		{object}	models.ApiResponse{data=nil}
// @Router			/tasks [post]
func v1_tasks_post(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	var tasks []models.ClearingTask
	if err := c.BodyParser(&tasks); err != nil {
		response.Message = err.Error()
		return c.Status(400).JSON(response)
	}

	// TODO: validation

	// endregion
	// region: data

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		response.Message = err.Error()
		return c.Status(500).JSON(response)
	}

	result := DB.Create(&tasks)
	if result.Error != nil {
		tx.Rollback()
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Message = err.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta["rows"] = len(tasks)
	response.Data = tasks
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_tasks_get

// //	@Summary		get tasks
// //	@Description	get filtered list of tasks
// //	@Tags			/tasks
// //	@Accept			json
// //	@Produce		json
// //	@Param			user_id		query		int	true	"user/grinder id"
// //	@Param			project_id	query		int	true	"project id"
// //	@Success		200			{object}	models.ApiResponse{data=[]models.ClearingTask}
// //	@Failure		400			{object}	models.ApiResponse{data=nil}
// //	@Failure		500			{object}	models.ApiResponse{data=nil}
// //	@Router			/tasks [get]
// func v1_tasks_get(c *fiber.Ctx) error {

// 	// region: output

// 	response := utils.GetResponse(c)

// 	// endregion
// 	// region: input

// 	// `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
// 	// `clearing_batch_id` BIGINT UNSIGNED NOT NULL,
// 	// `clearing_task_id` BIGINT UNSIGNED DEFAULT NULL,
// 	// `clearing_task_type_id` SMALLINT UNSIGNED NOT NULL,
// 	// `clearing_task_status_id` SMALLINT UNSIGNED NOT NULL,
// 	// `user_id` INT(11) NOT NULL,
// 	// `input` JSON DEFAULT NULL,
// 	// `output` JSON NOT NULL DEFAULT '[]',

// 	user_id := uint(c.QueryInt("user_id", 0))
// 	if user_id < 1 {
// 		response.Message = "invalid user_id"
// 		return c.Status(400).JSON(response)
// 	}
// 	project_id := uint(c.QueryInt("project_id", 0))
// 	if project_id < 1 {
// 		response.Message = "invalid project_id"
// 		return c.Status(400).JSON(response)
// 	}

// 	// endregion
// 	// region: data

// 	var fees []models.ClearingTaskFee
// 	result := DB.Where(&models.ClearingTaskFee{UserId: user_id, ProjectId: project_id}).Find(&fees)
// 	if result.Error != nil {
// 		Logger(LOG_ERR, "error while fetching fees for user/project", result.Error.Error())
// 		response.Message = result.Error.Error()
// 		return c.Status(500).JSON(response)
// 	}

// 	// endregion
// 	// region: response

// 	meta := response.Meta.(map[string]string)
// 	meta["rows"] = strconv.Itoa(len(fees))

// 	response.Data = fees
// 	response.Meta = meta
// 	response.Success = true

// 	return c.Status(200).JSON(response)

// 	// endregion: response

// }

// endregion
// region: v1_tasks_fees_delete

// @Summary		delete fee
// @Description	delete task/subtask fees for user per project
// @Tags			/tasks
// @Accept			json
// @Produce		json
// @Param			request	body		models.ClearingTaskFee	true	"json request body, fields not forming part of the composite index are ignored"
// @Success		200		{object}	models.ApiResponse{data=nil}
// @Failure		400		{object}	models.ApiResponse{data=nil}
// @Failure		500		{object}	models.ApiResponse{data=nil}
// @Router			/tasks/fees [delete]
func v1_tasks_fees_delete(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	fees := new(models.ClearingTaskFee)
	if err := c.BodyParser(&fees); err != nil {
		response.Message = err.Error()
		return c.Status(400).JSON(response)
	}

	// endregion
	// region: data

	result := DB.Delete(&fees)
	if result.Error != nil {
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Success = true
	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_tasks_fees_get

// @Summary		get fees
// @Description	get task/subtask fees per user per project
// @Tags			/tasks
// @Accept			json
// @Produce		json
// @Param			user_id		query		int	true	"user/grinder id"
// @Param			project_id	query		int	true	"project id"
// @Success		200			{object}	models.ApiResponse{data=[]models.ClearingTaskFee}
// @Failure		400			{object}	models.ApiResponse{data=nil}
// @Failure		500			{object}	models.ApiResponse{data=nil}
// @Router			/tasks/fees [get]
func v1_tasks_fees_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	user_id := uint(c.QueryInt("user_id", 0))
	if user_id < 1 {
		response.Message = "invalid user_id"
		return c.Status(400).JSON(response)
	}
	project_id := uint(c.QueryInt("project_id", 0))
	if project_id < 1 {
		response.Message = "invalid project_id"
		return c.Status(400).JSON(response)
	}

	// endregion
	// region: data

	var fees []models.ClearingTaskFee
	result := DB.Where(&models.ClearingTaskFee{UserId: user_id, ProjectId: project_id}).Find(&fees)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching fees for user/project", result.Error.Error())
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Meta["rows"] = len(fees)
	response.Data = fees
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_tasks_fees_post

// @Summary		set fees
// @Description	set task/subtask fees for user per project
// @Tags			/tasks
// @Accept			json
// @Produce		json
// @Param			request	body		models.ClearingTaskFee	true	"json request body"
// @Success		200		{object}	models.ApiResponse{data=[]models.ClearingTaskFee}
// @Failure		400		{object}	models.ApiResponse{data=nil}
// @Failure		500		{object}	models.ApiResponse{data=nil}
// @Router			/tasks/fees [post]
func v1_tasks_fees_post(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: input

	fees := new(models.ClearingTaskFee)
	if err := c.BodyParser(&fees); err != nil {
		response.Message = err.Error()
		return c.Status(400).JSON(response)
	}

	// TODO: validation
	//	* primary key "user/project id" (delete "id will be ignored")
	//	* user/project key already exists

	// endregion
	// region: data

	result := DB.Create(&fees)
	if result.Error != nil {
		response.Message = result.Error.Error()
		return c.Status(500).JSON(response)
	}

	// endregion
	// region: response

	response.Data = fees
	response.Success = true

	return c.Status(200).JSON(response)

	// endregion: response

}

// endregion
// region: v1_tasks_statuses_get

// @Summary		task statuses
// @Description	get list of possible task statuses
// @Tags			/tasks
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.ApiResponse{data=[]models.ClearingTaskStatus}
// @Failure		400	{object}	models.ApiResponse{data=nil}
// @Failure		500	{object}	models.ApiResponse{data=nil}
// @Router			/tasks/statuses [get]
func v1_tasks_statuses_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var statuses []models.ClearingTaskStatus
	result := DB.Find(&statuses)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching task statuses", result.Error.Error())
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
// region: v1_tasks_types_get

// @Summary		task types
// @Description	get list of possible task types
// @Tags			/tasks
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.ApiResponse{data=[]models.ClearingTaskType}
// @Failure		400	{object}	models.ApiResponse{data=nil}
// @Failure		500	{object}	models.ApiResponse{data=nil}
// @Router			/tasks/types [get]
func v1_tasks_types_get(c *fiber.Ctx) error {

	// region: output

	response := utils.GetResponse(c)

	// endregion
	// region: data

	var types []models.ClearingTaskType
	result := DB.Find(&types)
	if result.Error != nil {
		Logger(LOG_ERR, "error while fetching task types", result.Error.Error())
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
