//go:generate mockgen -destination=./mocks/mock_guide_repository.go -package=mocks github.com/Viverov/guideliner/internal/domains/guide/repository GuideRepository

package service

// @DirtyHack https://github.com/golang/mock/issues/494
import _ "github.com/golang/mock/mockgen/model"
