package entity

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity/condition"
	"time"
)

type NodeCreateOptions struct {
	text          string
	conditionType condition.CondType
	//Only for TypeTime ConditionType
	duration  time.Duration
	nextNodes []*node
}
