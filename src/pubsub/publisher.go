package pubsub

var subscribers = make([]Subscriber, 0)

func Subscribe(subscriber Subscriber) {
	subscribers = append(subscribers, subscriber)
}

func Unsubscribe(subscriber Subscriber) {
	for i, sub := range subscribers {
		if sub == subscriber {
			subscribers = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}
}

func Publish(topic string, data interface{}) {
	for _, sub := range subscribers {
		sub.Notify(topic, data)
	}
}
