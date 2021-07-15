package entity

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity/condition"
)

type node struct {
	//ConditionType     condition.CondType `json:"condition_type,omitempty"`
	//ConditionDuration time.Duration      `json:"condition_duration,omitempty"`
	Condition         condition.Condition `json:"condition,omitempty"`
	Text              string  `json:"text,omitempty"`
	NextNodes         []*node `json:"next_nodes,omitempty"`
}

func newNode(options NodeCreateOptions) (*node, error) {
	n := &node{}

	var cond condition.Condition
	switch options.conditionType {
	case condition.TypeManual:
		cond = condition.NewManualCondition()
	case condition.TypeTime:
		timeCond := condition.NewTimeCondition()
		timeCond.SetDuration(options.duration)
		cond = timeCond
	default:
		return nil, &InvalidConditionTypeError{}
	}
	n.Condition = cond

	if options.nextNodes != nil {
		n.NextNodes = options.nextNodes
	}

	return n, nil
}
