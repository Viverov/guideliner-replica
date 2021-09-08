package repository

import (
	"errors"
	"fmt"
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/util"
	"github.com/Viverov/guideliner/internal/domains/util/urepo"
	"gorm.io/gorm"
)

type guideRepositoryPsql struct {
	db *gorm.DB
}

type guideModel struct {
	gorm.Model
	Description string `gorm:"not null"`
	NodesJson   string `gorm:"not null"`
	CreatorID   uint
}

func (g guideModel) TableName() string {
	return "guides"
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, urepo.NewUnexpectedRepositoryError("FindById", result.Error.Error())
	}

	return entity.NewGuide(gm.ID, gm.NodesJson, gm.Description, gm.CreatorID)
}

func (r *guideRepositoryPsql) Find(condition FindConditions) ([]entity.Guide, error) {
	tx := util.SetDefaultConditions(r.db, condition.DefaultFindConditions)
	if condition.Search != "" {
		tx = tx.Where("Description ILIKE ?", fmt.Sprint("%", condition.Search, "%"))
	}

	var gms []*guideModel
	result := tx.Find(&gms)
	if result.Error != nil {
		return nil, urepo.NewUnexpectedRepositoryError("Find", result.Error.Error())
	}

	var guides []entity.Guide
	for _, gm := range gms {
		g, err := entity.NewGuide(gm.ID, gm.NodesJson, gm.Description, gm.CreatorID)
		if err != nil {
			return nil, err
		}
		guides = append(guides, g)
	}

	return guides, nil
}

func (r *guideRepositoryPsql) Count(cond CountConditions) (int64, error) {
	tx := r.db.Model(&guideModel{})
	if cond.Search != "" {
		tx = tx.Where("Description ILIKE ?", fmt.Sprint("%", cond.Search, "%"))
	}

	var count int64
	result := tx.Count(&count)

	if result.Error != nil {
		return 0, urepo.NewUnexpectedRepositoryError("Count", result.Error.Error())
	}

	return count, nil
}

func (r *guideRepositoryPsql) Insert(guide entity.Guide) (id uint, err error) {
	if guide == nil {
		return 0, urepo.NewNilEntityError("Guide")
	}
	nodesJson, err := guide.NodesToJSON()
	if err != nil {
		return 0, err
	}
	gm := &guideModel{
		Description: guide.Description(),
		NodesJson:   nodesJson,
		CreatorID:   guide.CreatorID(),
	}
	result := r.db.Create(gm)

	if result.Error != nil {
		return 0, urepo.NewUnexpectedRepositoryError("insert", result.Error.Error())
	}

	return gm.ID, nil
}

func (r *guideRepositoryPsql) Update(guide entity.Guide) error {
	if guide == nil {
		return urepo.NewNilEntityError("Guide")
	}

	g, err := r.FindById(guide.ID())
	if err != nil {
		return err
	}
	if g == nil {
		return urepo.NewEntityNotFoundError("Guide", guide.ID())
	}

	nodesJson, err := guide.NodesToJSON()
	if err != nil {
		return err
	}
	gm := &guideModel{
		Model: gorm.Model{
			ID: guide.ID(),
		},
		Description: guide.Description(),
		NodesJson:   nodesJson,
	}

	result := r.db.Save(gm)
	if result.Error != nil {
		return urepo.NewUnexpectedRepositoryError("update", result.Error.Error())
	}

	return nil
}
