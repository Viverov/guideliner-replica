package entity

import (
	"encoding/json"
	"fmt"
)

type guideImpl struct {
	id          uint
	description string
	rootNode    *node
	creatorID   uint
}

func NewGuide(id uint, nodesJson string, description string, creatorID uint) (*guideImpl, error) {
	g := &guideImpl{
		id:          id,
		description: "",
		rootNode:    nil,
		creatorID:   0,
	}
	g.SetDescription(description)
	err := g.SetNodesFromJSON(nodesJson)
	if err != nil {
		return nil, err
	}
	g.SetCreatorID(creatorID)

	return g, nil
}

func (g *guideImpl) ID() uint {
	return g.id
}

func (g *guideImpl) SetID(id uint) error {
	if id == 0 {
		return NewInvalidIdError()
	}
	g.id = id
	return nil
}

func (g *guideImpl) Description() string {
	return g.description
}

func (g *guideImpl) SetDescription(description string) {
	g.description = description
}

func (g *guideImpl) NodesToJSON() (string, error) {
	b, err := json.Marshal(g.rootNode)
	if err != nil {
		return "", NewUnexpectedGuideError(fmt.Sprintf("Can't marshal guide into json with error: %s", err.Error()))
	}

	return string(b), nil
}

func (g *guideImpl) SetNodesFromJSON(nodesJson string) error {
	rootNode := &node{}
	err := json.Unmarshal([]byte(nodesJson), rootNode)
	if err != nil {
		return NewInvalidJsonError()
	}
	g.rootNode = rootNode

	return nil
}

func (g *guideImpl) CreatorID() uint {
	return g.creatorID
}

func (g *guideImpl) SetCreatorID(creatorID uint) {
	g.creatorID = creatorID
}
