package cleanup

import (
	"context"
	"strconv"
	"time"

	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/setting"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// StartCleanupCron starts the general cleanup cron job(s).
func StartCleanupCron(heartbeatService heartbeat.Service, settingService setting.Service, logger *zap.SugaredLogger) {
	c := cron.New()

	// Heartbeat cleanup task
	c.AddFunc("0 * * * *", func() {
		keepDays := 365 // default fallback
		settingModel, err := settingService.GetByKey(context.Background(), "KEEP_DATA_PERIOD_DAYS")
		if err != nil {
			logger.Errorw("Failed to fetch keepDataPeriodDays setting", "error", err)
		} else if settingModel != nil {
			if v, err := strconv.Atoi(settingModel.Value); err == nil {
				keepDays = v
			} else {
				logger.Errorw("Invalid keepDataPeriodDays value", "value", settingModel.Value, "error", err)
			}
		}

		cutoff := time.Now().UTC().AddDate(0, 0, -keepDays)
		deleted, err := heartbeatService.DeleteOlderThan(context.Background(), cutoff)
		if err != nil {
			logger.Errorw("Failed to delete old heartbeats", "error", err)
			return
		}
		logger.Infow("Deleted old heartbeats", "count", deleted, "cutoff", cutoff)
	})

	c.Start()
}
