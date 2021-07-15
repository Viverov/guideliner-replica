package entity

import "fmt"

type RootNodeAlreadySetError struct{}

func (e *RootNodeAlreadySetError) Error() string {
	return fmt.Sprintf("Root nodeImpl already set")
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
	return fmt.Sprintf("ConditionType not found")
}

type UnexpectedGuideError struct {
	info string
}

func (u *UnexpectedGuideError) Error() string {
	return fmt.Sprintf("Unexpected guide error: %s", u.info)
}
