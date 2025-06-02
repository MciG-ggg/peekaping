package cleanup

import (
	"context"
	"time"

	"peekaping/src/modules/heartbeat"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// StartCleanupCron starts the general cleanup cron job(s).
func StartCleanupCron(heartbeatService heartbeat.Service, logger *zap.SugaredLogger) {
	c := cron.New()

	// Heartbeat cleanup task
	c.AddFunc("14 03 * * *", func() {
		// setting, err := settingService.FindByKey("keepDataPeriodDays")
		// if err != nil || setting == nil {
		// 	logger.Errorw("Failed to fetch keepDataPeriodDays setting", "error", err)
		// 	return
		// }
		// keepDays, err := strconv.Atoi(setting.Value)
		// if err != nil {
		// 	logger.Errorw("Invalid keepDataPeriodDays value", "value", setting.Value, "error", err)
		// 	return
		// }

		keepDays := 365 // TODO: make it configurable
		cutoff := time.Now().UTC().AddDate(0, 0, -keepDays)
		deleted, err := heartbeatService.DeleteOlderThan(context.Background(), cutoff)
		if err != nil {
			logger.Errorw("Failed to delete old heartbeats", "error", err)
			return
		}
		logger.Infow("Deleted old heartbeats", "count", deleted, "cutoff", cutoff)
	})

	// Future cleanup tasks can be added here

	c.Start()
}
