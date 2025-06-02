package monitor

import (
	"context"
	"fmt"
	"peekaping/src/modules/events"
	"peekaping/src/modules/healthcheck/executor"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/monitor_notification"
	"time"
)

type Service interface {
	Create(ctx context.Context, monitor *CreateUpdateDto) (*Model, error)
	FindByID(ctx context.Context, id string) (*Model, error)
	FindAll(ctx context.Context, page int, limit int, q string, active *bool, status *int) ([]*Model, error)
	FindActive(ctx context.Context) ([]*Model, error)
	UpdateFull(ctx context.Context, id string, monitor *CreateUpdateDto) (*Model, error)
	UpdatePartial(ctx context.Context, id string, monitor *PartialUpdateDto) (*Model, error)
	Delete(ctx context.Context, id string) error
	ValidateMonitorConfig(monitorType string, configJSON string) error

	GetMonitorChartPoints(ctx context.Context, id string, period string) ([]*heartbeat.ChartPoint, error)
	GetUptimeStats(ctx context.Context, id string) (*UptimeStatsDto, error)
	GetHeartbeats(ctx context.Context, id string, limit, page int, important *bool, reverse bool) ([]*heartbeat.Model, error)

	RemoveProxyReference(ctx context.Context, proxyId string) error
	FindByProxyId(ctx context.Context, proxyId string) ([]*Model, error)
}

type MonitorServiceImpl struct {
	monitorRepository          MonitorRepository
	heartbeatService           heartbeat.Service
	eventBus                   *events.EventBus
	monitorNotificationService monitor_notification.Service
	executorRegistry           *executor.ExecutorRegistry
}

func NewMonitorService(
	monitorRepository MonitorRepository,
	heartbeatService heartbeat.Service,
	eventBus *events.EventBus,
	monitorNotificationService monitor_notification.Service,
	executorRegistry *executor.ExecutorRegistry,
) Service {
	return &MonitorServiceImpl{
		monitorRepository,
		heartbeatService,
		eventBus,
		monitorNotificationService,
		executorRegistry,
	}
}

func (mr *MonitorServiceImpl) Create(ctx context.Context, monitorCreateDto *CreateUpdateDto) (*Model, error) {
	createModel := &Model{
		Type:           monitorCreateDto.Type,
		Name:           monitorCreateDto.Name,
		Interval:       monitorCreateDto.Interval,
		Timeout:        monitorCreateDto.Timeout,
		MaxRetries:     monitorCreateDto.MaxRetries,
		RetryInterval:  monitorCreateDto.RetryInterval,
		ResendInterval: monitorCreateDto.ResendInterval,
		Active:         monitorCreateDto.Active,
		Status:         heartbeat.MonitorStatusUp,
		CreatedAt:      time.Now().UTC(),
		Config:         monitorCreateDto.Config,
		ProxyId:        monitorCreateDto.ProxyId,
	}

	createdModel, err := mr.monitorRepository.Create(ctx, createModel)
	if err != nil {
		return nil, err
	}

	// Emit monitor created event
	mr.eventBus.Publish(events.Event{
		Type:    events.MonitorCreated,
		Payload: createdModel,
	})

	return createdModel, nil
}

func (mr *MonitorServiceImpl) FindByID(ctx context.Context, id string) (*Model, error) {
	return mr.monitorRepository.FindByID(ctx, id)
}

func (mr *MonitorServiceImpl) FindAll(ctx context.Context, page int, limit int, q string, active *bool, status *int) ([]*Model, error) {
	monitors, err := mr.monitorRepository.FindAll(ctx, page, limit, q, active, status)
	if err != nil {
		return nil, err
	}

	return monitors, nil
}

func (mr *MonitorServiceImpl) FindActive(ctx context.Context) ([]*Model, error) {
	return mr.monitorRepository.FindActive(ctx)
}

func (mr *MonitorServiceImpl) UpdateFull(ctx context.Context, id string, monitor *CreateUpdateDto) (*Model, error) {
	model := &Model{
		ID:             id,
		Name:           monitor.Name,
		Type:           monitor.Type,
		Interval:       monitor.Interval,
		Timeout:        monitor.Timeout,
		MaxRetries:     monitor.MaxRetries,
		RetryInterval:  monitor.RetryInterval,
		ResendInterval: monitor.ResendInterval,
		Active:         monitor.Active,
		Status:         heartbeat.MonitorStatusUp,
		UpdatedAt:      time.Now().UTC(),
		Config:         monitor.Config,
		ProxyId:        monitor.ProxyId,
	}

	err := mr.monitorRepository.UpdateFull(ctx, id, model)
	if err != nil {
		return nil, err
	}

	// Emit monitor updated event
	mr.eventBus.Publish(events.Event{
		Type:    events.MonitorUpdated,
		Payload: model,
	})

	return model, nil
}

func (mr *MonitorServiceImpl) UpdatePartial(ctx context.Context, id string, monitor *PartialUpdateDto) (*Model, error) {
	model := &UpdateModel{
		ID:             &id,
		Type:           monitor.Type,
		Name:           monitor.Name,
		Interval:       monitor.Interval,
		Timeout:        monitor.Timeout,
		MaxRetries:     monitor.MaxRetries,
		RetryInterval:  monitor.RetryInterval,
		ResendInterval: monitor.ResendInterval,
		Active:         monitor.Active,
		Status:         monitor.Status,
	}

	err := mr.monitorRepository.UpdatePartial(ctx, id, model)
	if err != nil {
		return nil, err
	}

	// Get the updated monitor
	updatedMonitor, err := mr.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Emit monitor updated event
	mr.eventBus.Publish(events.Event{
		Type:    events.MonitorUpdated,
		Payload: updatedMonitor,
	})

	return updatedMonitor, nil
}

func (mr *MonitorServiceImpl) Delete(ctx context.Context, id string) error {
	err := mr.monitorRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Cascade delete monitor_notification relations
	_ = mr.monitorNotificationService.DeleteByMonitorID(ctx, id)

	// Emit monitor deleted event
	mr.eventBus.Publish(events.Event{
		Type:    events.MonitorDeleted,
		Payload: id,
	})

	return nil
}

func (mr *MonitorServiceImpl) GetMonitorChartPoints(
	ctx context.Context,
	id string,
	period string,
) ([]*heartbeat.ChartPoint, error) {
	now := time.Now().UTC()

	periods := map[string]time.Duration{
		"30m":   30 * time.Minute,
		"3h":    3 * time.Hour,
		"6h":    6 * time.Hour,
		"24h":   24 * time.Hour,
		"1week": 7 * 24 * time.Hour,
	}

	duration, exists := periods[period]
	if !exists {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	startTime := now.Add(-duration)

	heartbeats, err := mr.heartbeatService.FindByMonitorIDAndTimeRange(ctx, id, startTime, now)
	if err != nil {
		return nil, err
	}

	return heartbeats, nil
}

// GetUptimeStats returns uptime percentages for 24h, 7d, 30d, 365d
func (mr *MonitorServiceImpl) GetUptimeStats(ctx context.Context, id string) (*UptimeStatsDto, error) {
	now := time.Now().UTC()
	periods := map[string]time.Duration{
		"24h":  24 * time.Hour,
		"7d":   7 * 24 * time.Hour,
		"30d":  30 * 24 * time.Hour,
		"365d": 365 * 24 * time.Hour,
	}

	statsMap, err := mr.heartbeatService.FindUptimeStatsByMonitorID(ctx, id, periods, now)
	if err != nil {
		return nil, err
	}

	stats := &UptimeStatsDto{
		Uptime24h:  statsMap["24h"],
		Uptime7d:   statsMap["7d"],
		Uptime30d:  statsMap["30d"],
		Uptime365d: statsMap["365d"],
	}

	return stats, nil
}

func (mr *MonitorServiceImpl) ValidateMonitorConfig(
	monitorType string,
	configJSON string,
) error {
	if mr.executorRegistry == nil {
		return fmt.Errorf("executor registry not available")
	}
	return mr.executorRegistry.ValidateConfig(monitorType, configJSON)
}

func (mr *MonitorServiceImpl) GetHeartbeats(ctx context.Context, id string, limit, page int, important *bool, reverse bool) ([]*heartbeat.Model, error) {
	return mr.heartbeatService.FindByMonitorIDPaginated(ctx, id, limit, page, important, reverse)
}

func (mr *MonitorServiceImpl) RemoveProxyReference(ctx context.Context, proxyId string) error {
	return mr.monitorRepository.RemoveProxyReference(ctx, proxyId)
}

func (mr *MonitorServiceImpl) FindByProxyId(ctx context.Context, proxyId string) ([]*Model, error) {
	return mr.monitorRepository.FindByProxyId(ctx, proxyId)
}
