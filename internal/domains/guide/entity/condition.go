package entity

import "time"

type CondType string

const (
	TypeManual CondType = "MANUAL"
	TypeTime   CondType = "TIME"
)

type Condition struct {
	CondType `json:"type,omitempty"`
	//Duration only for TIME condition
	Duration time.Duration `json:"duration,omitempty"`
}

func NewManualCondition() *Condition {
	return &Condition{
		CondType: TypeManual,
	}
}

func NewTimeCondition(duration time.Duration) *Condition {
	return &Condition{
		Duration: duration,
		CondType: TypeTime,
	}
}
