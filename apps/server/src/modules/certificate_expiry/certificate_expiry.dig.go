package certificate_expiry

import (
	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewService)
}