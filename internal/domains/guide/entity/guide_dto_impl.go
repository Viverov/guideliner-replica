package entity

type guideDTOImpl struct {
	id          uint
	description string
	nodesJSON   string
}

func NewGuideDTO(id uint, description string, nodesJSON string) *guideDTOImpl {
	return &guideDTOImpl{
		id:          id,
		description: description,
		nodesJSON:   nodesJSON,
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
