package entity

type Guide interface {
	ID() uint
	SetID(id uint) error
	Description() string
	SetDescription(description string)
	NodesToJSON() (string, error)
	SetNodesFromJSON(string) error
	CreatorID() uint
	SetCreatorID(uint)
}
