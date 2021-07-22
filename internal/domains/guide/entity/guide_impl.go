package entity

import (
	"encoding/json"
	"fmt"
)

type guideImpl struct {
	id          uint
	description string
	rootNode    *node
}

func NewGuide() *guideImpl {
	return &guideImpl{
		id:          0,
		rootNode:    nil,
		description: "",
	}
}

func NewGuideWithParams(id uint, nodesJson string, description string) (*guideImpl, error) {
	if id == 0 {
		return nil, &InvalidIdError{}
	}

	g := &guideImpl{
		id:          id,
		description: "",
		rootNode:   nil,
	}
	g.SetDescription(description)
	err := g.SetNodesFromJSON(nodesJson)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *guideImpl) ID() uint {
	return g.id
}

func (g *guideImpl) SetID(id uint) error {
	if id == 0 {
		return &InvalidIdError{}
	}
	g.id = id
	return nil}

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
		return &InvalidJsonError{}
	}
	g.rootNode = rootNode

	return nil
}
