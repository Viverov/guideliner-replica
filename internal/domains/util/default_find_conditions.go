package util

const DefaultLimit = 20

type DefaultFindConditions struct {
	Limit  int
	Offset int
	//Order in format "+atr1-atr2" (that means atr1 ASC atr2 DESC)
	Order string
}

func (fc DefaultFindConditions) ResolveLimit() int {
	if fc.Limit != 0 {
		return fc.Limit
	}

	return DefaultLimit
}
