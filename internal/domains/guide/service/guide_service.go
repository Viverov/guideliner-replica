package service

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/util"
)

type GuideService interface {
	Find(FindConditions) ([]entity.GuideDTO, error)
	FindById(id uint) (entity.GuideDTO, error)
	Count(CountConditions) (count int64, err error)
	Create(description string, nodesJson string) (entity.GuideDTO, error)
	Update(id uint, params UpdateParams) error
}

type FindConditions struct {
	util.DefaultFindConditions
	Search string
}

type CountConditions struct {
	Search string
}

type UpdateParams struct {
	Description string
	NodesJson   string
}
