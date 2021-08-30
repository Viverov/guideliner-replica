package service

type InvalidNodesJsonError struct{}

func NewInvalidNodesJsonError() *InvalidNodesJsonError {
	return &InvalidNodesJsonError{}
}

func (e *InvalidNodesJsonError) Error() string {
	return "Can't parse json into correct nodes"
}
