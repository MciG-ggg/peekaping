package monitor

import (
	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewMonitorRepository)
	container.Provide(NewMonitorService)
	container.Provide(NewMonitorController)
	container.Provide(NewMonitorRoute)
}
