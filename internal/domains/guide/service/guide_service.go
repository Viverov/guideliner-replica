package service

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/util"
)

type GuideService interface {
	Find(FindConditions) ([]entity.GuideDTO, error)
	FindById(id uint) (entity.GuideDTO, error)
	Create(description string, nodesJson string) (entity.GuideDTO, error)
	Update(id uint, params UpdateParams) error
}

type FindConditions struct {
	util.DefaultFindConditions
	Search string
}

type UpdateParams struct {
	Description string
	NodesJson   string
}
