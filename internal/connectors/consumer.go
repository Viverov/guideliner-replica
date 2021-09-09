package connectors

type Consumer interface {
	GetChannel() chan<- Event
}
