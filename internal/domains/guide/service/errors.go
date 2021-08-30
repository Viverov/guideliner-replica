package service

import "fmt"

type GuideNotFoundError struct {
	id uint
}

func (e *GuideNotFoundError) Error() string {
	return fmt.Sprintf("guide with id %d not found", e.id)
}

type InvalidNodesJsonError struct{}

func (e *InvalidNodesJsonError) Error() string {
	return "Can't parse json into correct nodes"
}
