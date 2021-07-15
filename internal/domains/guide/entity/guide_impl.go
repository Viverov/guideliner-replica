package entity

import (
	"encoding/json"
	"fmt"
)

type guideImpl struct {
	id               uint
	description      string
	rootNode         *node
}

func NewGuide() *guideImpl {
	return &guideImpl{
		rootNode:         nil,
		description:      "",
	}
}

func NewGuideWithParams(id uint, nodesJson string, description string) (*guideImpl, error) {
	rootNode := &node{}
	err := json.Unmarshal([]byte(nodesJson), rootNode)
	if err != nil {
		return nil, &UnexpectedGuideError{info: err.Error()}
	}

	return &guideImpl{
		id:               id,
		description:      description,
		rootNode:         rootNode,
	}, nil
}

func (g *guideImpl) ID() uint {
	return g.id
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
		return "", &UnexpectedGuideError{info: fmt.Sprintf("Can't marshal guide into json with error: %s", err.Error())}
	}

	return string(b), nil
}

func (g *guideImpl) SetNodesFromJSON(nodesJson string) error {
	rootNode := &node{}
	err := json.Unmarshal([]byte(nodesJson), rootNode)
	if err != nil {
		return &UnexpectedGuideError{info: err.Error()}
	}

	return nil
}
