package condition

import (
	"time"
)

type TimeCondition struct {
	Condition `json:"condition,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
}

func (t *TimeCondition) SetDuration(duration time.Duration) {
	t.Duration = duration
}

func NewTimeCondition() *TimeCondition {
	return &TimeCondition{
		Condition: baseCondition{CondType: TypeTime},
	}
}
