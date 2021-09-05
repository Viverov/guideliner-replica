package utils

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/domains/util"
	"strings"
)

// ParseOrderQuery parse orders from http requests.
// orderStr: source string from query. Ex: "+id,-created_at".
// mappers: contains replacers for queryString into service format. Ex: created_at -> createdAt
func ParseOrderQuery(orderStr string, mappers map[string]string) ([]util.Order, error) {
	if orderStr == "" {
		return []util.Order{}, nil
	}

	parts := strings.Split(orderStr, ",")

	var result []util.Order
	for _, part := range parts {
		direction, field := part[:1], part[1:]
		var o util.Order

		switch direction {
		case "+":
			o.Direction = util.ASC
		case "-":
			o.Direction = util.DESC
		default:
			return nil, fmt.Errorf("can't parse order, expected string like \"+field1,-field2\", received %s", part)
		}

		resolvedField, ok := mappers[field]
		if !ok {
			return nil, fmt.Errorf("undefined field: %s", field)
		}
		o.Field = resolvedField

		result = append(result, o)
	}

	return result, nil
}
