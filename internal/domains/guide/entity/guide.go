package entity

type Guide interface {
	ID() uint
	Description() string
	SetDescription(description string)
	//CurrentMaxNodeId() uint
	//SetRootNode(options NodeCreateOptions) (*node, error)
	//SetNextNodes(nodeId uint, nextNodes []NodeCreateOptions) ([]*node, error)
	NodesToJSON() (string, error)
	SetNodesFromJSON(string) error
}
