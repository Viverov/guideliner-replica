package condition

type CondType string

const (
	TypeManual CondType = "MANUAL"
	TypeTime   CondType = "TIME"
)

type Condition interface {
	Type() CondType
}
