package uservice

import "github.com/Viverov/guideliner/internal/domains/util/urepo"

// CheckDefaultRepoErrors contains default checks for repo errors
func CheckDefaultRepoErrors(err error) error {
	switch t := err.(type) {
	case *urepo.UnexpectedRepositoryError:
		return NewStorageError(t.Error())
	case *urepo.EntityNotFoundError:
		return NewNotFoundError(t.EntityName(), t.Id())
	default:
		return NewUnexpectedServiceError()
	}
}
