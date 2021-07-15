package guide

import (
	"github.com/Viverov/guideliner/internal/domains/guide/repository"
	"github.com/Viverov/guideliner/internal/domains/guide/service"
	"gorm.io/gorm"
)

func BuildGuideService(db *gorm.DB) service.GuideService {
	return service.NewGuideService(repository.NewGuideRepositoryPsql(db))
}
