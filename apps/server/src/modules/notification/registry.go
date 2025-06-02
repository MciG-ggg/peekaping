package notification

var NotificationIntegrationRegistry = make(map[string]NotificationSender)

func RegisterNotifier(name string, notifier NotificationSender) {
	NotificationIntegrationRegistry[name] = notifier
}

func GetNotifier(name string) (NotificationSender, bool) {
	n, ok := NotificationIntegrationRegistry[name]
	return n, ok
}
