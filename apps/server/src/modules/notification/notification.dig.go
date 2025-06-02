package notification

import (
	"peekaping/src/modules/monitor_notification"

	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewRepository)
	container.Provide(monitor_notification.NewService)
	container.Provide(func(repo Repository, monitorNotificationService monitor_notification.Service) Service {
		return NewService(repo, monitorNotificationService)
	})
	container.Provide(NewController)
	container.Provide(NewRoute)
	container.Provide(NewNotificationEventListener)
}
