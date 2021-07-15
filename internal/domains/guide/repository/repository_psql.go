package repository

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"gorm.io/gorm"
)

type guideModel struct {
	gorm.Model
	//CurrentMaxNodeId uint   `gorm:"not null"`
	Description      string `gorm:"not null"`
	NodesJson        string `gorm:"not null"`
}

type guideRepositoryPsql struct {
	db *gorm.DB
}

func NewGuideRepositoryPsql(db *gorm.DB) *guideRepositoryPsql {
	return &guideRepositoryPsql{
		db: db,
	}
}

func (r *guideRepositoryPsql) FindById(id uint) (entity.Guide, error) {
	gm := &guideModel{
		Model: gorm.Model{
			ID: id,
		},
	}

	result := r.db.First(gm)
	if result.Error != nil {
		return nil, &CommonRepositoryError{
			action:    "FindById",
			errorText: result.Error.Error(),
		}
	}

	//return entity.NewGuideWithParams(gm.ID, gm.NodesJson, gm.Description, gm.CurrentMaxNodeId)
	return entity.NewGuideWithParams(gm.ID, gm.NodesJson, gm.Description)
}

func (r *guideRepositoryPsql) Find(condition FindConditions) ([]entity.Guide, error) {
	if condition.Search == "" {
		return nil, &InvalidFindConditionError{}
	}

	var gms []guideModel
	result := r.db.Limit(condition.ResolveLimit()).Find(gms)
	if result.Error != nil {
		return nil, &CommonRepositoryError{
			action:    "Find",
			errorText: result.Error.Error(),
		}
	}

	var guides []entity.Guide
	for _, gm := range gms {
		g, err := entity.NewGuideWithParams(gm.ID, gm.NodesJson, gm.Description)
		if err != nil {
			return nil, err
		}
		guides = append(guides, g)
	}

	return guides, nil
}

func (r *guideRepositoryPsql) Insert(guide entity.Guide) (id uint, err error) {
	nodesJson, err := guide.NodesToJSON()
	if err != nil {
		return 0, err
	}
	gm := &guideModel{
		//CurrentMaxNodeId: guide.CurrentMaxNodeId,
		Description:      guide.Description(),
		NodesJson:        nodesJson,
	}
	result := r.db.Create(gm)

	if result.Error != nil {
		return 0, &CommonRepositoryError{
			action:    "insert",
			errorText: result.Error.Error(),
		}
	}

	return gm.ID, nil
}

func (r *guideRepositoryPsql) Update(guide entity.Guide) error {
	nodesJson, err := guide.NodesToJSON()
	if err != nil {
		return err
	}
	gm := &guideModel{
		Model: gorm.Model{
			ID: guide.ID(),
		},
		//CurrentMaxNodeId: guide.CurrentMaxNodeId(),
		Description:      guide.Description(),
		NodesJson:        nodesJson,
	}

	result := r.db.Save(gm)
	if result.Error != nil {
		return &CommonRepositoryError{
			action:    "update",
			errorText: result.Error.Error(),
		}
	}

	return nil
}
