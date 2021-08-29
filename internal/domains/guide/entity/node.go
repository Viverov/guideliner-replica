package entity

import (
	"time"
)

type node struct {
	Condition *Condition `json:"condition,omitempty"`
	Text      string     `json:"text,omitempty"`
	NextNodes []*node    `json:"next_nodes,omitempty"`
}

func newNode(options NodeCreateOptions) (*node, error) {
	n := &node{}

	var cond *Condition
	switch options.conditionType {
	case TypeManual:
		cond = NewManualCondition()
	case TypeTime:
		timeCond := NewTimeCondition(options.duration)
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

type NodeCreateOptions struct {
	text          string
	conditionType CondType
	//Only for TypeTime ConditionType
	duration  time.Duration
	nextNodes []*node
}
