package entity

type node struct {
	Condition *Condition `json:"condition,omitempty"`
	Text      string     `json:"text,omitempty"`
	NextNodes []*node    `json:"next_nodes,omitempty"`
}
