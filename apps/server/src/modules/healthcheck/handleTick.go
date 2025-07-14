package healthcheck

import (
	"context"
	"peekaping/src/modules/events"
	"peekaping/src/modules/healthcheck/executor"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/proxy"
	"peekaping/src/modules/shared"
	"time"
)

// isImportantForNotification determines if a heartbeat is important for notification purposes.
func (s *HealthCheckSupervisor) isImportantForNotification(prevBeatStatus, currBeatStatus heartbeat.MonitorStatus) bool {
	up := shared.MonitorStatusUp
	down := shared.MonitorStatusDown
	pending := shared.MonitorStatusPending
	maintenance := shared.MonitorStatusMaintenance

	// * ? -> ANY STATUS = important [isFirstBeat]
	// UP -> PENDING = not important
	// * UP -> DOWN = important
	// UP -> UP = not important
	// PENDING -> PENDING = not important
	// * PENDING -> DOWN = important
	// PENDING -> UP = not important
	// DOWN -> PENDING = this case not exists
	// DOWN -> DOWN = not important
	// * DOWN -> UP = important
	// MAINTENANCE -> MAINTENANCE = not important
	// MAINTENANCE -> UP = not important
	// * MAINTENANCE -> DOWN = important
	// DOWN -> MAINTENANCE = not important
	// UP -> MAINTENANCE = not important

	return (prevBeatStatus == maintenance && currBeatStatus == down) ||
		(prevBeatStatus == up && currBeatStatus == down) ||
		(prevBeatStatus == down && currBeatStatus == up) ||
		(prevBeatStatus == pending && currBeatStatus == down)
}

// isImportantBeat determines if the status of the monitor has changed in an important way since the last beat.
func (s *HealthCheckSupervisor) isImportantBeat(prevBeatStatus, currBeatStatus heartbeat.MonitorStatus) bool {
	up := shared.MonitorStatusUp
	down := shared.MonitorStatusDown
	pending := shared.MonitorStatusPending
	maintenance := shared.MonitorStatusMaintenance

	// UP -> PENDING = not important
	// * UP -> DOWN = important
	// UP -> UP = not important
	// PENDING -> PENDING = not important
	// * PENDING -> DOWN = important
	// PENDING -> UP = not important
	// DOWN -> PENDING = this case not exists
	// DOWN -> DOWN = not important
	// * DOWN -> UP = important
	// MAINTENANCE -> MAINTENANCE = not important
	// * MAINTENANCE -> UP = important
	// * MAINTENANCE -> DOWN = important
	// * DOWN -> MAINTENANCE = important
	// * UP -> MAINTENANCE = important

	return (prevBeatStatus == down && currBeatStatus == maintenance) ||
		(prevBeatStatus == up && currBeatStatus == maintenance) ||
		(prevBeatStatus == maintenance && currBeatStatus == down) ||
		(prevBeatStatus == maintenance && currBeatStatus == up) ||
		(prevBeatStatus == up && currBeatStatus == down) ||
		(prevBeatStatus == down && currBeatStatus == up) ||
		(prevBeatStatus == pending && currBeatStatus == down)
}

func (s *HealthCheckSupervisor) postProcessHeartbeat(result *executor.Result, m *Monitor, intervalUpdateCb func(newInterval time.Duration)) {
	ping := int(result.EndTime.Sub(result.StartTime).Milliseconds())

	ctx := context.Background()

	// get the previous heartbeat
	previousBeats, err := s.heartbeatService.FindByMonitorIDPaginated(ctx, m.ID, 1, 0, nil, false)
	var previousBeat *heartbeat.Model = nil
	if err != nil {
		s.logger.Errorf("Failed to get previous heartbeat for monitor %s: %v", m.ID, err)
	}
	if len(previousBeats) > 0 {
		previousBeat = previousBeats[0]
	}

	s.logger.Debugf("previousBeat %t", previousBeat != nil)

	isFirstBeat := previousBeat == nil

	hb := &heartbeat.CreateUpdateDto{
		MonitorID: m.ID,
		Status:    result.Status,
		Msg:       result.Message,
		Ping:      ping,
		Duration:  0,
		DownCount: 0,
		Retries:   0,
		Important: false,
		Time:      result.StartTime,
		EndTime:   result.EndTime,
		Notified:  false,
	}

	if !isFirstBeat {
		hb.DownCount = previousBeat.DownCount
		hb.Retries = previousBeat.Retries
	}

	// mark as pending if max retries is set and retries is less than max retries
	if result.Status == shared.MonitorStatusDown {
		if !isFirstBeat && m.MaxRetries > 0 && previousBeat.Retries < m.MaxRetries {
			hb.Status = shared.MonitorStatusPending
		}
		if intervalUpdateCb != nil {
			intervalUpdateCb(time.Duration(m.RetryInterval) * time.Second)
		}
		hb.Retries++
	} else {
		if intervalUpdateCb != nil {
			intervalUpdateCb(time.Duration(m.Interval) * time.Second)
		}
		hb.Retries = 0
	}

	s.logger.Debugf("isFirstBeat for: %s %t", m.Name, isFirstBeat)
	s.logger.Debugf("checking if important for: %s", m.Name)
	isImportant := isFirstBeat || s.isImportantBeat(previousBeat.Status, hb.Status)
	s.logger.Debugf("isImportant for %s: %t", m.Name, isImportant)

	shouldNotify := false

	// if important (beat status changed), send notification
	if isImportant {
		hb.Important = true

		// update monitor status
		// s.monitorSvc.UpdatePartial(m.ID, &monitor.UpdateDto{
		// 	Status: &hb.Status,
		// })

		if isFirstBeat || s.isImportantForNotification(previousBeat.Status, hb.Status) {
			s.logger.Debugf("sending notification %s", m.Name)
			shouldNotify = true
			hb.Notified = true
		} else {
			s.logger.Debugf("not sending notification %s", m.Name)
		}

		hb.DownCount = 0
	} else {
		hb.Important = false

		if result.Status == shared.MonitorStatusDown && m.ResendInterval > 0 {
			hb.DownCount += 1

			if hb.DownCount >= m.ResendInterval {
				shouldNotify = true
				hb.Notified = true
				hb.DownCount = 0
			}
		}
	}

	if result.Status == shared.MonitorStatusUp {
		s.logger.Debugf("%s successful response %d ms | interval %d seconds | type %s", m.Name, ping, m.Interval, m.Type)
	} else if result.Status == shared.MonitorStatusPending {
		s.logger.Debugf("%s pending response %d ms | interval %d seconds | type %s", m.Name, ping, m.Interval, m.Type)
	} else if result.Status == shared.MonitorStatusDown {
		s.logger.Debugf("%s down response %d ms | interval %d seconds | type %s", m.Name, ping, m.Interval, m.Type)
	} else if result.Status == shared.MonitorStatusMaintenance {
		s.logger.Debugf("%s maintenance response %d ms | interval %d seconds | type %s", m.Name, ping, m.Interval, m.Type)
	}

	dbHb, err := s.heartbeatService.Create(ctx, hb)
	if err != nil {
		s.logger.Errorf("Failed to create heartbeat", err.Error())
		return
	}

	// Check certificate expiry notifications if monitor is up and has TLS info
	if result.Status == shared.MonitorStatusUp && result.TLSInfo != nil {
		err := s.certificateExpiryService.CheckCertExpiryNotifications(ctx, result, m)
		if err != nil {
			s.logger.Errorf("Failed to check certificate expiry notifications for monitor %s: %v", m.ID, err)
		}
	}

	if shouldNotify {
		s.eventBus.Publish(events.Event{
			Type:    events.MonitorStatusChanged,
			Payload: dbHb,
		})
	}
}

// handleMonitorTick processes a single monitor tick in its own goroutine.
func (s *HealthCheckSupervisor) handleMonitorTick(
	ctx context.Context,
	m *Monitor,
	exec executor.Executor,
	proxyModel *proxy.Model,
	intervalUpdateCb func(newInterval time.Duration),
) {
	// Check if monitor is under maintenance
	isUnderMaintenance, err := s.isUnderMaintenance(ctx, m.ID)
	s.logger.Debugf("isUnderMaintenance for %s: %t", m.Name, isUnderMaintenance)
	if err != nil {
		s.logger.Errorf("Failed to check maintenance status for monitor %s: %v", m.ID, err)
	}

	if isUnderMaintenance {
		// If under maintenance, create a maintenance status heartbeat
		result := &executor.Result{
			Status:    shared.MonitorStatusMaintenance,
			Message:   "Monitor under maintenance",
			StartTime: time.Now(),
			EndTime:   time.Now(),
		}
		s.postProcessHeartbeat(result, m, intervalUpdateCb)
		return
	}

	callCtx, cCancel := context.WithTimeout(
		ctx,
		time.Duration(m.Timeout)*time.Second,
	)
	defer cCancel()

	// Execute the health check
	result := exec.Execute(callCtx, m, proxyModel)
	if result == nil {
		return
	}

	s.postProcessHeartbeat(result, m, intervalUpdateCb)
}
