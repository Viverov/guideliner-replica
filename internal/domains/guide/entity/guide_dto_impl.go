package entity

type guideDTOImpl struct {
	id          uint
	description string
	nodesJSON   string
	creatorID   uint
}

func NewGuideDTO(id uint, description string, nodesJSON string, creatorID uint) *guideDTOImpl {
	return &guideDTOImpl{
		id:          id,
		description: description,
		nodesJSON:   nodesJSON,
		creatorID:   creatorID,
	}
}

func NewGuideDTOFromEntity(guide Guide) (*guideDTOImpl, error) {
	nodesJson, err := guide.NodesToJSON()
	if err != nil {
		return nil, err
	}
	return &guideDTOImpl{
		id:          guide.ID(),
		description: guide.Description(),
		nodesJSON:   nodesJson,
		creatorID:   guide.CreatorID(),
	}, nil
}

func (g *guideDTOImpl) ID() uint {
	return g.id
}

func (g *guideDTOImpl) Description() string {
	return g.description
}

func (g *guideDTOImpl) NodesJson() string {
	return g.nodesJSON
}

func (g *guideDTOImpl) CreatorID() uint {
	return g.creatorID
}
