package healthcheck

import (
	"context"
	"peekaping/src/modules/events"
	"peekaping/src/modules/healthcheck/executor"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/shared"
	"time"
)

// isImportantForNotification determines if a heartbeat is important for notification purposes.
func (s *HealthCheckSupervisor) isImportantForNotification(prevBeatStatus, currBeatStatus heartbeat.MonitorStatus) bool {
	up := heartbeat.MonitorStatusUp
	down := heartbeat.MonitorStatusDown
	pending := heartbeat.MonitorStatusPending
	maintenance := heartbeat.MonitorStatusMaintenance

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
	up := heartbeat.MonitorStatusUp
	down := heartbeat.MonitorStatusDown
	pending := heartbeat.MonitorStatusPending
	maintenance := heartbeat.MonitorStatusMaintenance

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

// handleHeartbeatTick processes a single monitor tick in its own goroutine.
func (s *HealthCheckSupervisor) handleMonitorTick(
	ctx context.Context,
	m *Monitor,
	executor executor.Executor,
	proxyModel *shared.Proxy,
	intervalUpdateCb func(newInterval time.Duration),
) {
	callCtx, cCancel := context.WithTimeout(
		ctx,
		time.Duration(m.Timeout)*time.Second,
	)
	defer cCancel()

	result := executor.Execute(callCtx, m, proxyModel)

	ping := int(result.EndTime.Sub(result.StartTime).Milliseconds())

	internalCtx := context.Background()
	// get the previous heartbeat
	previousBeats, err := s.heartbeatService.FindByMonitorIDPaginated(internalCtx, m.ID, 1, 0, nil, false)
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
	if result.Status == heartbeat.MonitorStatusDown {
		if !isFirstBeat && m.MaxRetries > 0 && previousBeat.Retries < m.MaxRetries {
			hb.Status = heartbeat.MonitorStatusPending
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

		if result.Status == heartbeat.MonitorStatusDown && m.ResendInterval > 0 {
			hb.DownCount += 1

			if hb.DownCount >= m.ResendInterval {
				shouldNotify = true
				hb.Notified = true
				hb.DownCount = 0
			}
		}
	}

	if result.Status == heartbeat.MonitorStatusUp {
		s.logger.Debugf("%s successful response %d ms | interval %d seconds | type %s", m.Name, ping, m.Interval, m.Type)
	} else if result.Status == heartbeat.MonitorStatusPending {
		s.logger.Debugf("%s pending response %d ms | interval %d seconds | type %s", m.Name, ping, m.Interval, m.Type)
	} else if result.Status == heartbeat.MonitorStatusDown {
		s.logger.Debugf("%s down response %d ms | interval %d seconds | type %s", m.Name, ping, m.Interval, m.Type)
	} else if result.Status == heartbeat.MonitorStatusMaintenance {
		s.logger.Debugf("%s maintenance response %d ms | interval %d seconds | type %s", m.Name, ping, m.Interval, m.Type)
	}

	// TODO: calculate uptime

	dbHb, err := s.heartbeatService.Create(internalCtx, hb)
	if err != nil {
		s.logger.Errorf("Failed to create heartbeat", err.Error())
		return
	}

	s.eventBus.Publish(events.Event{
		Type:    events.HeartbeatEvent,
		Payload: dbHb,
	})

	if shouldNotify {
		s.eventBus.Publish(events.Event{
			Type:    events.MonitorStatusChanged,
			Payload: dbHb,
		})
	}
}
