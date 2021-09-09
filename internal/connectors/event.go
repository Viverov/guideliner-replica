package connectors

type EventName string

const (
	DumbEventName EventName = "dumb"
)

type Event interface {
	Name() EventName
}

type DumbEvent struct {
	Dumb string `json:"dumb"`
}

func (d *DumbEvent) Name() EventName {
	return DumbEventName
}
