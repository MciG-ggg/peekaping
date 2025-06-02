package notification

import (
	"net/http"
	"peekaping/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate = validator.New()

type Controller struct {
	service Service
	logger  *zap.SugaredLogger
}

func NewController(
	service Service,
	logger *zap.SugaredLogger,
) *Controller {
	return &Controller{
		service,
		logger,
	}
}

// @Router		/notifications [get]
// @Summary		Get notifications
// @Tags			Notifications
// @Produce		json
// @Security  BearerAuth
// @Param     q    query     string  false  "Search query"
// @Param     page query     int     false  "Page number" default(1)
// @Param     limit query    int     false  "Items per page" default(10)
// @Success		200	{object}	utils.ApiResponse[[]Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *Controller) FindAll(ctx *gin.Context) {
	// Extract query parameters for pagination and search
	page, err := utils.GetQueryInt(ctx, "page", 0)
	if err != nil || page < 0 {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid page parameter"))
		return
	}

	limit, err := utils.GetQueryInt(ctx, "limit", 10)
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid limit parameter"))
		return
	}

	q := ctx.Query("q")

	response, err := ic.service.FindAll(ctx, page, limit, q)
	if err != nil {
		ic.logger.Errorw("Failed to fetch notifications", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("success", response))
}

// @Router		/notifications [post]
// @Summary		Create notification
// @Tags			Notifications
// @Produce		json
// @Accept		json
// @Security  BearerAuth
// @Param     body body   CreateUpdateDto  true  "Notification object"
// @Success		201	{object}	utils.ApiResponse[Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *Controller) Create(ctx *gin.Context) {
	var notification *CreateUpdateDto
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ic.logger.Errorw("Invalid request body", "error", err)
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid request body"))
		return
	}

	if err := validate.Struct(notification); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid request body"))
		return
	}

	integration, ok := GetNotifier(notification.Type)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Unsupported notification type"))
		return
	}
	err := integration.Validate(notification.Config)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid config: "+err.Error()))
		return
	}

	createdNotification, err := ic.service.Create(ctx, notification)
	if err != nil {
		ic.logger.Errorw("Failed to create notification", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	ctx.JSON(http.StatusCreated, utils.NewSuccessResponse("Notification created successfully", createdNotification))
}

// @Router		/notifications/{id} [get]
// @Summary		Get notification by ID
// @Tags			Notifications
// @Produce		json
// @Security BearerAuth
// @Param       id   path      string  true  "Notification ID"
// @Success		200	{object}	utils.ApiResponse[Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *Controller) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	notification, err := ic.service.FindByID(ctx, id)
	if err != nil {
		ic.logger.Errorw("Failed to fetch notification", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	if notification == nil {
		ctx.JSON(http.StatusNotFound, utils.NewFailResponse("Notification not found"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("success", notification))
}

// @Router		/notifications/{id} [put]
// @Summary		Update notification
// @Tags			Notifications
// @Produce		json
// @Accept		json
// @Security BearerAuth
// @Param       id   path      string  true  "Notification ID"
// @Param       notification body     CreateUpdateDto  true  "Notification object"
// @Success		200	{object}	utils.ApiResponse[Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *Controller) UpdateFull(ctx *gin.Context) {
	id := ctx.Param("id")

	var notification CreateUpdateDto
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	if err := validate.Struct(notification); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	updatedNotification, err := ic.service.UpdateFull(ctx, id, &notification)
	if err != nil {
		ic.logger.Errorw("Failed to update notification", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("notification updated successfully", updatedNotification))
}

// @Router		/notifications/{id} [patch]
// @Summary		Update notification
// @Tags			Notifications
// @Produce		json
// @Accept		json
// @Security BearerAuth
// @Param       id   path      string  true  "Notification ID"
// @Param       notification body     PartialUpdateDto  true  "Notification object"
// @Success		200	{object}	utils.ApiResponse[Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *Controller) UpdatePartial(ctx *gin.Context) {
	id := ctx.Param("id")

	var notification PartialUpdateDto
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	// validate
	if err := validate.Struct(notification); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	updatedNotification, err := ic.service.UpdatePartial(ctx, id, &notification)
	if err != nil {
		ic.logger.Errorw("Failed to update notification", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("notification updated successfully", updatedNotification))
}

// @Router		/notifications/{id} [delete]
// @Summary		Delete notification
// @Tags			Notifications
// @Produce		json
// @Security BearerAuth
// @Param       id   path      string  true  "Notification ID"
// @Success		200	{object}	utils.ApiResponse[any]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ic.service.Delete(ctx, id)
	if err != nil {
		ic.logger.Errorw("Failed to delete notification", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse[any]("Notification deleted successfully", nil))
}
