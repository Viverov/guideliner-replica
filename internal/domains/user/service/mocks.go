//go:generate mockgen -destination=./mocks/mock_user_repository.go -package=mocks github.com/Viverov/guideliner/internal/domains/user/repository UserRepository
//go:generate mockgen -destination=./mocks/mock_token_provider.go -package=mocks github.com/Viverov/guideliner/internal/domains/user/token_provider TokenProvider

package service

// @DirtyHack https://github.com/golang/mock/issues/494
import _ "github.com/golang/mock/mockgen/model"
