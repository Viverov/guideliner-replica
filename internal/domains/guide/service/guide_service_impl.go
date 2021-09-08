package service

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/guide/repository"
	"github.com/Viverov/guideliner/internal/domains/util/uservice"
)

type guideServiceImpl struct {
	repository repository.GuideRepository
}

func NewGuideService(rep repository.GuideRepository) *guideServiceImpl {
	return &guideServiceImpl{repository: rep}
}

func (s *guideServiceImpl) Find(cond FindConditions) ([]entity.GuideDTO, error) {
	guides, err := s.repository.Find(repository.FindConditions{
		DefaultFindConditions: cond.DefaultFindConditions,
		Search:                cond.Search,
	})
	if err != nil {
		return nil, processRepositoryError(err)
	}

	var dtos []entity.GuideDTO
	for _, guide := range guides {
		dto, err := entity.NewGuideDTOFromEntity(guide)
		if err != nil {
			return nil, uservice.NewUnexpectedServiceError()
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

func (s *guideServiceImpl) FindById(id uint) (entity.GuideDTO, error) {
	guide, err := s.repository.FindById(id)
	if err != nil {
		return nil, processRepositoryError(err)
	}
	if guide == nil {
		return nil, nil
	}

	dto, err := entity.NewGuideDTOFromEntity(guide)
	if err != nil {
		return nil, uservice.NewUnexpectedServiceError()
	}

	return dto, nil
}

func (s *guideServiceImpl) Count(cond CountConditions) (int64, error) {
	count, err := s.repository.Count(repository.CountConditions{Search: cond.Search})
	if err != nil {
		return 0, processRepositoryError(err)
	}

	return count, nil
}

func (s *guideServiceImpl) Create(description string, nodesJson string, creatorID uint) (entity.GuideDTO, error) {
	guide, err := entity.NewGuide(0, nodesJson, description, creatorID)
	if err != nil {
		return nil, NewInvalidNodesJsonError()
	}

	id, err := s.repository.Insert(guide)
	if err != nil {
		return nil, processRepositoryError(err)
	}

	err = guide.SetID(id)
	if err != nil {
		return nil, err
	}

	dto, err := entity.NewGuideDTOFromEntity(guide)
	if err != nil {
		return nil, uservice.NewUnexpectedServiceError()
	}

	return dto, nil
}

func (s *guideServiceImpl) Update(id uint, params UpdateParams) (entity.GuideDTO, error) {
	guide, err := s.repository.FindById(id)
	if err != nil {
		return nil, processRepositoryError(err)
	}
	if guide == nil {
		return nil, uservice.NewNotFoundError("Guide", id)
	}

	if params.Description != "" {
		guide.SetDescription(params.Description)
	}
	if params.NodesJson != "" {
		err := guide.SetNodesFromJSON(params.NodesJson)
		if err != nil {
			return nil, NewInvalidNodesJsonError()
		}
	}

	err = s.repository.Update(guide)
	if err != nil {
		return nil, processRepositoryError(err)
	}

	dto, err := entity.NewGuideDTOFromEntity(guide)
	if err != nil {
		return nil, uservice.NewUnexpectedServiceError()
	}

	return dto, nil
}

func (s *guideServiceImpl) CheckPermission(guideID uint, userID uint, permission Permission) (bool, error) {
	permissions, err := s.GetPermissions(guideID, userID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}
	return false, nil
}

func (s *guideServiceImpl) GetPermissions(guideID uint, userID uint) ([]Permission, error) {
	guide, err := s.repository.FindById(guideID)
	if err != nil {
		return nil, processRepositoryError(err)
	}
	if guide == nil {
		return nil, uservice.NewNotFoundError("Guide", guideID)
	}

	permissions := []Permission{}
	if guide.CreatorID() == userID {
		permissions = append(permissions, PermissionUpdate)
	}

	return permissions, nil
}

func processRepositoryError(err error) error {
	return uservice.CheckDefaultRepoErrors(err)
}
