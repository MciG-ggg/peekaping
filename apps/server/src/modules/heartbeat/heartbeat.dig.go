package heartbeat

import "go.uber.org/dig"

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewRepository)
	container.Provide(NewService)
}
