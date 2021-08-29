package util

import (
	"gorm.io/gorm"
	"strings"
)

const DefaultLimit = 20

type DirectionType string

const (
	ASC  DirectionType = "ASC"
	DESC DirectionType = "DESC"
)

type Order struct {
	Field     string
	Direction DirectionType
}

type DefaultFindConditions struct {
	Limit  int
	Offset int
	Order  []Order
}

func (fc DefaultFindConditions) ResolveLimit() int {
	if fc.Limit != 0 {
		return fc.Limit
	}

	return DefaultLimit
}

func (fc DefaultFindConditions) ParseOrderCondition() string {
	if fc.Order == nil || len(fc.Order) == 0 {
		return "id ASC"
	}

	var responseParts []string
	for _, o := range fc.Order {
		if o.Direction == ASC {
			responseParts = append(responseParts, o.Field+" ASC")
		} else {
			responseParts = append(responseParts, o.Field+" DESC")
		}
	}

	return strings.Join(responseParts, ",")
}

func SetDefaultConditions(db *gorm.DB, c DefaultFindConditions) *gorm.DB {
	return db.
		Order(c.ParseOrderCondition()).
		Limit(c.ResolveLimit()).
		Offset(c.Offset)
}
