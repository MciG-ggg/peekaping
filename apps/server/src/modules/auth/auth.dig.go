package auth

import (
	"peekaping/src/config"

	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container, cfg *config.Config) {
	switch cfg.DBType {
	case "postgres", "postgresql", "mysql", "sqlite":
		container.Provide(NewSQLRepository)
	case "mongo":
		container.Provide(NewMongoRepository)
	}

	container.Provide(NewRoute)
	container.Provide(NewTokenMaker)
	container.Provide(NewService)
	container.Provide(NewController)
	container.Provide(NewMiddlewareProvider)
}
