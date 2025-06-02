package monitor

import (
	"fmt"
	"net/http"
	"peekaping/src/modules/monitor_notification"
	"peekaping/src/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate = validator.New()

type MonitorController struct {
	monitorService             Service
	logger                     *zap.SugaredLogger
	monitorNotificationService monitor_notification.Service
}

func NewMonitorController(
	monitorService Service,
	logger *zap.SugaredLogger,
	monitorNotificationService monitor_notification.Service,
) *MonitorController {
	validate.RegisterStructValidation(CreateUpdateDtoStructLevelValidation, CreateUpdateDto{})

	return &MonitorController{
		monitorService,
		logger,
		monitorNotificationService,
	}
}

// @Router		/monitors [get]
// @Summary		Get monitors
// @Tags			Monitors
// @Produce		json
// @Security  BearerAuth
// @Param     q    query     string  false  "Search query"
// @Param     page query     int     false  "Page number" default(1)
// @Param     limit query    int     false  "Items per page" default(10)
// @Param     active query   bool    false  "Active status"
// @Param     status query   int     false  "Status"
// @Success		200	{object}	utils.ApiResponse[[]Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *MonitorController) FindAll(ctx *gin.Context) {
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

	active, err := utils.GetQueryBool(ctx, "active")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid active parameter (must be true or false)"))
		return
	}

	var statusPtr *int
	if statusStr := ctx.Query("status"); statusStr != "" {
		statusVal, err := utils.GetQueryInt(ctx, "status", 0)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid status parameter (must be int)"))
			return
		}
		statusPtr = &statusVal
	}

	response, err := ic.monitorService.FindAll(ctx, page, limit, q, active, statusPtr)
	if err != nil {
		ic.logger.Errorw("Failed to fetch monitors", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("success", response))
}

// @Router		/monitors [post]
// @Summary		Create monitor
// @Tags			Monitors
// @Produce		json
// @Accept		json
// @Security  BearerAuth
// @Param     body body   CreateUpdateDto  true  "Monitor object"
// @Success		201	{object}	utils.ApiResponse[Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *MonitorController) Create(ctx *gin.Context) {
	var monitor *CreateUpdateDto
	if err := ctx.ShouldBindJSON(&monitor); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	if err := validate.Struct(monitor); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	// Validate monitor type and config
	if err := ic.monitorService.ValidateMonitorConfig(monitor.Type, monitor.Config); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(fmt.Sprintf("Invalid monitor configuration: %v", err)))
		return
	}

	createdMonitor, err := ic.monitorService.Create(ctx, monitor)
	if err != nil {
		ic.logger.Errorw("Failed to create monitor", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}
	ic.logger.Infof("Created monitor: %+v\n", createdMonitor)

	// Handle multiple notification IDs
	if len(monitor.NotificationIds) > 0 {
		for _, notificationId := range monitor.NotificationIds {
			_, err = ic.monitorNotificationService.Create(ctx, createdMonitor.ID, notificationId)
			if err != nil {
				ic.logger.Errorw("Failed to create monitor-notification record", "error", err)
				ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
				return
			}
		}
	}

	ctx.JSON(http.StatusCreated, utils.NewSuccessResponse("Monitor created successfully", createdMonitor))
}

// @Router		/monitors/{id} [get]
// @Summary		Get monitor by ID
// @Tags			Monitors
// @Produce		json
// @Security BearerAuth
// @Param       id   path      string  true  "Monitor ID"
// @Success		200	{object}	utils.ApiResponse[MonitorResponseDto]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *MonitorController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	monitor, err := ic.monitorService.FindByID(ctx, id)
	if err != nil {
		ic.logger.Errorw("Failed to fetch monitor", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	if monitor == nil {
		ctx.JSON(http.StatusNotFound, utils.NewFailResponse("Monitor not found"))
		return
	}

	// Fetch notification_ids
	notificationRels, err := ic.monitorNotificationService.FindByMonitorID(ctx, id)
	if err != nil {
		ic.logger.Errorw("Failed to fetch monitor-notification relations", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}
	notificationIds := make([]string, 0, len(notificationRels))
	for _, rel := range notificationRels {
		notificationIds = append(notificationIds, rel.NotificationID)
	}

	// Compose response with notification_ids
	response := MonitorResponseDto{
		ID:              monitor.ID,
		Name:            monitor.Name,
		Interval:        monitor.Interval,
		Timeout:         monitor.Timeout,
		Type:            monitor.Type,
		Active:          monitor.Active,
		MaxRetries:      monitor.MaxRetries,
		RetryInterval:   monitor.RetryInterval,
		ResendInterval:  monitor.ResendInterval,
		Status:          int(monitor.Status),
		CreatedAt:       monitor.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       monitor.UpdatedAt.Format(time.RFC3339),
		NotificationIds: notificationIds,
		ProxyId:         monitor.ProxyId,
		Config:          monitor.Config,
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("success", response))
}

// @Router		/monitors/{id} [put]
// @Summary		Update monitor
// @Tags			Monitors
// @Produce		json
// @Accept		json
// @Security BearerAuth
// @Param       id   path      string  true  "Monitor ID"
// @Param       monitor body     CreateUpdateDto  true  "Monitor object"
// @Success		200	{object}	utils.ApiResponse[Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *MonitorController) UpdateFull(ctx *gin.Context) {
	id := ctx.Param("id")

	var monitor CreateUpdateDto
	if err := ctx.ShouldBindJSON(&monitor); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	// validate
	if err := validate.Struct(monitor); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	// if float64(monitor.Timeout)*0.8 >= float64(monitor.Interval) {
	// 	ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Timeout cannot be greater than 80% of interval"))
	// 	return
	// }

	// Validate monitor type and config
	if err := ic.monitorService.ValidateMonitorConfig(monitor.Type, monitor.Config); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(fmt.Sprintf("Invalid monitor configuration: %v", err)))
		return
	}

	updatedMonitor, err := ic.monitorService.UpdateFull(ctx, id, &monitor)
	if err != nil {
		ic.logger.Errorw("Failed to update monitor", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	// Delete all existing notification relations and create new ones
	err = ic.monitorNotificationService.DeleteByMonitorID(ctx, id)
	if err != nil {
		ic.logger.Errorw("Failed to delete existing monitor-notification relations", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	// Create new notification relations
	for _, notificationId := range monitor.NotificationIds {
		_, err = ic.monitorNotificationService.Create(ctx, id, notificationId)
		if err != nil {
			ic.logger.Errorw("Failed to create monitor-notification record", "error", err)
			ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
			return
		}
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Monitor updated successfully", updatedMonitor))
}

// @Router		/monitors/{id} [patch]
// @Summary		Update monitor
// @Tags			Monitors
// @Produce		json
// @Accept		json
// @Security BearerAuth
// @Param       id   path      string  true  "Monitor ID"
// @Param       monitor body     PartialUpdateDto  true  "Monitor object"
// @Success		200	{object}	utils.ApiResponse[Model]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *MonitorController) UpdatePartial(ctx *gin.Context) {
	id := ctx.Param("id")

	var monitor PartialUpdateDto
	if err := ctx.ShouldBindJSON(&monitor); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error()))
		return
	}

	// if monitor.Timeout != nil && float64(*monitor.Timeout)*0.8 >= float64(*monitor.Interval) {
	// 	ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Timeout cannot be greater than 80% of interval"))
	// 	return
	// }

	// Validate monitor type and config if they are being updated
	if monitor.Type != nil && monitor.Config != nil {
		if err := ic.monitorService.ValidateMonitorConfig(*monitor.Type, *monitor.Config); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.NewFailResponse(fmt.Sprintf("Invalid monitor configuration: %v", err)))
			return
		}
	}

	updatedMonitor, err := ic.monitorService.UpdatePartial(ctx, id, &monitor)
	if err != nil {
		ic.logger.Errorw("Failed to update monitor", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	// Handle notification IDs if they are being updated
	if len(monitor.NotificationIds) > 0 {
		// Replace all monitor-notification relations in an optimized way
		existing, err := ic.monitorNotificationService.FindByMonitorID(ctx, id)
		if err != nil {
			ic.logger.Errorw("Failed to fetch monitor-notification relations", "error", err)
			ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
			return
		}

		// Build sets for comparison
		existingMap := make(map[string]string) // notificationID -> relationID
		for _, rel := range existing {
			existingMap[rel.NotificationID] = rel.ID
		}
		newSet := make(map[string]struct{})
		for _, nid := range monitor.NotificationIds {
			newSet[nid] = struct{}{}
		}

		// Delete relations not in the new list
		for notificationID, relID := range existingMap {
			if _, found := newSet[notificationID]; !found {
				if err := ic.monitorNotificationService.Delete(ctx, relID); err != nil {
					ic.logger.Warnw("Failed to delete monitor-notification relation", "error", err)
				}
			}
		}

		// Add new relations not already present
		for _, nid := range monitor.NotificationIds {
			if _, found := existingMap[nid]; !found {
				if _, err := ic.monitorNotificationService.Create(ctx, id, nid); err != nil {
					ic.logger.Warnw("Failed to create monitor-notification relation", "error", err)
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Monitor updated successfully", updatedMonitor))
}

// @Router		/monitors/{id} [delete]
// @Summary		Delete monitor
// @Tags			Monitors
// @Produce		json
// @Security BearerAuth
// @Param       id   path      string  true  "Monitor ID"
// @Success		200	{object}	utils.ApiResponse[any]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *MonitorController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ic.monitorService.Delete(ctx, id)
	if err != nil {
		ic.logger.Errorw("Failed to delete monitor", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse[any]("Monitor deleted successfully", nil))
}

// @Router		/monitors/{id}/chartpoints [get]
// @Summary		Get monitor chart points
// @Tags			Monitors
// @Produce		json
// @Security BearerAuth
// @Param       id   path      string  true  "Monitor ID"
// @Param       period query   string  true  "Time period (30m, 3h, 6h, 24h, 1week)" Enums(30m, 3h, 6h, 24h, 1week)
// @Success		200	{object}	utils.ApiResponse[[]heartbeat.ChartPoint]
// @Failure		400	{object}	utils.APIError[any]
// @Failure		404	{object}	utils.APIError[any]
// @Failure		500	{object}	utils.APIError[any]
func (ic *MonitorController) GetMonitorChartPoints(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get and validate period parameter
	period := ctx.Query("period")
	validPeriods := map[string]bool{
		"30m":   true,
		"3h":    true,
		"6h":    true,
		"24h":   true,
		"1week": true,
	}

	if !validPeriods[period] {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid period. Must be one of: 30m, 3h, 6h, 24h, 1week"))
		return
	}

	heartbeats, err := ic.monitorService.GetMonitorChartPoints(ctx, id, period)
	if err != nil {
		ic.logger.Errorw("Failed to get monitor heartbeat", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("success", heartbeats))
}

// @Router	/monitors/{id}/heartbeats [get]
// @Summary	Get paginated heartbeats for a monitor
// @Tags		Monitors
// @Produce	json
// @Security BearerAuth
// @Param	id	path	string	true	"Monitor ID"
// @Param	limit	query	int	false	"Number of heartbeats per page (default 50)"
// @Param	page	query	int	false	"Page number (default 0)"
// @Param	important	query	bool	false	"Filter by important heartbeats only"
// @Param	reverse	query	bool	false	"Reverse the order of heartbeats"
// @Success	200	{object}	utils.ApiResponse[[]heartbeat.Model]
// @Failure	400	{object}	utils.APIError[any]
// @Failure	404	{object}	utils.APIError[any]
// @Failure	500	{object}	utils.APIError[any]
func (ic *MonitorController) FindByMonitorIDPaginated(ctx *gin.Context) {
	id := ctx.Param("id")

	limit, err := utils.GetQueryInt(ctx, "limit", 50)
	if err != nil || limit < 1 || limit > 1000 {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid limit parameter (1-1000)"))
		return
	}

	page, err := utils.GetQueryInt(ctx, "page", 0)
	if err != nil || page < 0 {
		ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid page parameter (>=0)"))
		return
	}

	var importantPtr *bool
	if ctx.Query("important") != "" {
		importantPtr, err = utils.GetQueryBool(ctx, "important")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid important parameter (must be true or false)"))
			return
		}
	}

	reverse := false
	if ctx.Query("reverse") != "" {
		reversePtr, err := utils.GetQueryBool(ctx, "reverse")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Invalid reverse parameter (must be true or false)"))
			return
		}
		if reversePtr != nil {
			reverse = *reversePtr
		}
	}

	results, err := ic.monitorService.GetHeartbeats(ctx, id, limit, page, importantPtr, reverse)
	if err != nil {
		ic.logger.Errorw("Failed to get heartbeats", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("success", results))
}

// @Router	/monitors/{id}/uptime [get]
// @Summary	Get monitor uptime stats (24h, 7d, 30d, 365d)
// @Tags		Monitors
// @Produce	json
// @Security BearerAuth
// @Param	id	path	string	true	"Monitor ID"
// @Success	200	{object}	utils.ApiResponse[UptimeStatsDto]
// @Failure	400	{object}	utils.APIError[any]
// @Failure	404	{object}	utils.APIError[any]
// @Failure	500	{object}	utils.APIError[any]
func (ic *MonitorController) GetUptimeStats(ctx *gin.Context) {
	id := ctx.Param("id")

	stats, err := ic.monitorService.GetUptimeStats(ctx, id)
	if err != nil {
		ic.logger.Errorw("Failed to get uptime stats", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.NewFailResponse("Internal server error"))
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("success", stats))
}
