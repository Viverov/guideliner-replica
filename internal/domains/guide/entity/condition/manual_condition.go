package condition

type ManualCondition struct {
	Condition `json:"condition,omitempty"`
}

func NewManualCondition() *ManualCondition {
	return &ManualCondition{
		Condition: baseCondition{
			CondType: TypeManual,
		},
	}
}
