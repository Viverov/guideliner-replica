package connectors

type Producer interface {
	WriteEvents([]Event) error
}
