package service

import (
	"errors"
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/guide/repository"
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
			return nil, &UnexpectedServiceError{}
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

func (s *guideServiceImpl) FindById(id uint) (entity.GuideDTO, error) {
	guide, err := s.findEntityById(id)
	if err != nil {
		return nil, err
	}

	dto, err := entity.NewGuideDTOFromEntity(guide)
	if err != nil {
		return nil, &UnexpectedServiceError{}
	}

	return dto, nil
}

func (s *guideServiceImpl) Create(description string, nodesJson string) (entity.GuideDTO, error) {
	guide := entity.NewGuide()

	guide.SetDescription(description)
	err := guide.SetNodesFromJSON(nodesJson)
	if err != nil {
		return nil, &InvalidNodesJsonError{}
	}

	id, err := s.repository.Insert(guide)
	if err != nil {
		return nil, processRepositoryError(err)
	}

	dto := entity.NewGuideDTO(id, guide.Description(), nodesJson)

	return dto, nil
}

func (s *guideServiceImpl) Update(id uint, params UpdateParams) error {
	guide, err := s.findEntityById(id)
	if err != nil {
		return err
	}

	if params.Description != "" {
		guide.SetDescription(params.Description)
	}
	if params.NodesJson != "" {
		err := guide.SetNodesFromJSON(params.NodesJson)
		if err != nil {
			return &InvalidNodesJsonError{}
		}
	}

	return nil
}

func (s *guideServiceImpl) findEntityById(id uint) (entity.Guide, error) {
	guide, err := s.repository.FindById(id)
	if err != nil {
		return nil, processRepositoryError(err)
	}
	if guide == nil {
		return nil, &GuideNotFoundError{id: id}
	}

	return guide, err
}

func processRepositoryError(err error) error {
	var cre *repository.CommonRepositoryError
	if errors.As(err, &cre) {
		return &StorageError{storageErrorText: cre.Error()}
	}
	return &UnexpectedServiceError{}
}
