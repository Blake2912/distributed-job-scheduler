package eventbus

// Will be storing all of the types of messages that would be needed to be used as part of the event bus
type TTLExpiredEvent struct {
	Key string
}
