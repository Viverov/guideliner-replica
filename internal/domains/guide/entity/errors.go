package entity

import "fmt"

type RootNodeAlreadySetError struct{}

func (e *RootNodeAlreadySetError) Error() string {
	return fmt.Sprintf("root nodeImpl already set")
}

type NodeNotFoundError struct {
	nodeId uint
}

func (n *NodeNotFoundError) Error() string {
	return fmt.Sprintf("node with %d not found.", n.nodeId)
}

type InvalidConditionTypeError struct {
}

func (i *InvalidConditionTypeError) Error() string {
	return fmt.Sprintf("conditionType not found")
}

type UnexpectedGuideError struct {
	info string
}

func (u *UnexpectedGuideError) Error() string {
	return fmt.Sprintf("unexpected guide error: %s", u.info)
}

type InvalidJsonError struct{}

func (e *InvalidJsonError) Error() string {
	return fmt.Sprintf("invalid json passed")
}

type InvalidIdError struct{}

func (e *InvalidIdError) Error() string {
	return fmt.Sprintf("id must be above zero")
}
