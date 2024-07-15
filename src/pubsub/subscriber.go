package pubsub

type Subscriber interface {
	Notify(topic string, data interface{})
}
