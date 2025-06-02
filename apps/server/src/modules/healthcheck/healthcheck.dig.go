package healthcheck

import (
	"peekaping/src/modules/healthcheck/executor"

	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewHealthCheck)
	container.Provide(NewEventListener)
	container.Provide(executor.NewExecutorRegistry)
}
