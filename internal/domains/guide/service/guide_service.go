package service

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/util"
)

type GuideService interface {
	Find(FindConditions) ([]entity.GuideDTO, error)
	FindById(id uint) (entity.GuideDTO, error)
	Count(CountConditions) (count int64, err error)
	Create(description string, nodesJson string, creatorID uint) (entity.GuideDTO, error)
	Update(id uint, params UpdateParams) (entity.GuideDTO, error)
	CheckPermission(guideID uint, userID uint, permission Permission) (bool, error)
	GetPermissions(guideID uint, userID uint) ([]Permission, error)
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

type Permission string

const (
	PermissionUpdate Permission = "UPDATE"
)
