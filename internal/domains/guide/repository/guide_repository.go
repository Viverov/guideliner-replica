package repository

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/util"
)

type GuideRepository interface {
	FindById(uint) (entity.Guide, error)
	Find(FindConditions) ([]entity.Guide, error)
	Count(CountConditions) (count int64, err error)
	Insert(entity.Guide) (id uint, err error)
	Update(entity.Guide) error
}

type FindConditions struct {
	util.DefaultFindConditions
	Search string
}

type CountConditions struct {
	Search string
}
