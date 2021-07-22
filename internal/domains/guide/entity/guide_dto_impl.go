package entity

type guideDTOImpl struct {
	id          uint
	description string
	nodesJson   string
}

func NewGuideDTOFromEntity(guide Guide) (*guideDTOImpl, error) {
	nodesJson, err := guide.NodesToJSON()
	if err != nil {
		return nil, err
	}
	return &guideDTOImpl{
		id:          guide.ID(),
		description: guide.Description(),
		nodesJson:   nodesJson,
	}, nil
}

func (g *guideDTOImpl) ID() uint {
	return g.id
}

func (g *guideDTOImpl) Description() string {
	return g.description
}

func (g *guideDTOImpl) NodesJson() string {
	return g.nodesJson
}
