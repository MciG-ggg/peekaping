package stats

import (
	"peekaping/src/config"
	"peekaping/src/modules/events"

	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container, cfg *config.Config) {
	switch cfg.DBType {
	case "postgres", "postgresql", "mysql", "sqlite":
		container.Provide(NewSQLRepository)
	case "mongo":
		container.Provide(NewMongoRepository)
	}
	container.Provide(NewService)
	container.Invoke(func(s Service, bus *events.EventBus) {
		if impl, ok := s.(*ServiceImpl); ok {
			impl.RegisterEventHandlers(bus)
		}
	})
}
