package entity

type GuideDTO interface {
	ID() uint
	Description() string
	NodesJson() string
	CreatorID() uint
}
