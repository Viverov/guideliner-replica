package entity

type Guide interface {
	ID() uint
	Description() string
	SetDescription(description string)
	NodesToJSON() (string, error)
	SetNodesFromJSON(string) error
}
