package auth

import (
	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewRoute)
	container.Provide(NewRepository)
	container.Provide(NewTokenMaker)
	container.Provide(NewService)
	container.Provide(NewController)
	container.Provide(NewMiddlewareProvider)
}
