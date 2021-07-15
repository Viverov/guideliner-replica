package condition

type baseCondition struct {
	CondType `json:"cond_type,omitempty"`
}

func (b baseCondition) Type() CondType {
	return b.CondType
}
