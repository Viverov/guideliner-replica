// +build unit

package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var oneMinuteDuration, _ = time.ParseDuration("1m")

func TestNewGuide(t *testing.T) {
	tests := []struct {
		name string
		want *guideImpl
	}{
		{
			name: "Should return empty guide",
			want: &guideImpl{
				id:          0,
				rootNode:    nil,
				description: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGuide()
			assert.Equal(t, tt.want.id, got.id)
			assert.Equal(t, tt.want.rootNode, got.rootNode)
			assert.Equal(t, tt.want.description, got.description)
		})
	}
}

func TestNewGuideWithParams(t *testing.T) {
	type args struct {
		id          uint
		nodesJson   string
		description string
	}
	tests := []struct {
		name    string
		args    args
		want    *guideImpl
		wantErr error
	}{
		{
			name: "Should create guide",
			args: args{
				id:          10,
				nodesJson:   "{}",
				description: "test description",
			},
			want: &guideImpl{
				id: 10,
				rootNode: &node{
					Condition: nil,
					Text:      "",
					NextNodes: nil,
				},
				description: "test description",
			},
			wantErr: nil,
		},
		{
			name: "Should throw error on invalid json",
			args: args{
				id:          10,
				nodesJson:   "asdasdasdasd",
				description: "test description",
			},
			want:    nil,
			wantErr: &InvalidJsonError{},
		},
		{
			name: "Should throw error on zero id",
			args: args{
				id:          0,
				nodesJson:   "{}",
				description: "test description",
			},
			want:    nil,
			wantErr: &InvalidIdError{},
		},
		{
			name: "Should parse nodes",
			args: args{
				id:          10,
				nodesJson:   "{\"condition\":{\"type\":\"MANUAL\"},\"text\":\"node_1_text\",\"next_nodes\":[{\"condition\":{\"type\":\"TIME\",\"duration\":60000000000},\"text\":\"inner_node_1_text\"}]}",
				description: "test description",
			},
			want: &guideImpl{
				id:          10,
				description: "test description",
				rootNode: &node{
					Condition: NewManualCondition(),
					Text:      "node_1_text",
					NextNodes: []*node{
						{
							Condition: NewTimeCondition(oneMinuteDuration),
							Text:      "inner_node_1_text",
							NextNodes: []*node{},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should parse nodes with external fields",
			args: args{
				id:          10,
				nodesJson:   "{\"some\":\"field\",\"condition\":{\"type\":\"MANUAL\"},\"text\":\"node_1_text\"}",
				description: "test description",
			},
			want: &guideImpl{
				id:          10,
				description: "test description",
				rootNode: &node{
					Condition: NewManualCondition(),
					Text:      "node_1_text",
					NextNodes: []*node{},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGuideWithParams(tt.args.id, tt.args.nodesJson, tt.args.description)

			if tt.wantErr == nil {
				assert.Nil(t, err)
				assert.NotNil(t, got)

				// Check base info
				assert.Equal(t, tt.want.id, got.id)
				assert.Equal(t, tt.want.description, got.description)

				// Check nodes
				checkNodes(t, tt.want.rootNode, got.rootNode)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func Test_guideImpl_Description(t *testing.T) {
	type fields struct {
		id          uint
		description string
		rootNode    *node
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should return description",
			fields: fields{10, "desc", &node{}},
			want:   "desc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideImpl{
				id:          tt.fields.id,
				description: tt.fields.description,
				rootNode:    tt.fields.rootNode,
			}
			assert.Equal(t, tt.fields.description, g.Description())
		})
	}
}

func Test_guideImpl_ID(t *testing.T) {
	type fields struct {
		id          uint
		description string
		rootNode    *node
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name:   "Should return ID",
			fields: fields{10, "desc", &node{}},
			want:   10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideImpl{
				id:          tt.fields.id,
				description: tt.fields.description,
				rootNode:    tt.fields.rootNode,
			}
			assert.Equal(t, tt.want, g.ID())
		})
	}
}

func Test_guideImpl_SetID(t *testing.T) {
	type fields struct {
		id          uint
		description string
		rootNode    *node
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		fields fields
		args    args
		wantId  uint
		wantErr error
	}{
		{
			name:    "Should set new ID",
			fields: fields{
				id:          1,
				description: "",
				rootNode:    nil,
			},
			args:    args{id: 50},
			wantId: 50,
			wantErr: nil,
		},
		{
			name:    "Should return InvalidIdError",
			fields:  fields{
				id:          1,
				description: "",
				rootNode:    nil,
			},
			args:    args{0},
			wantId:  1,
			wantErr: &InvalidIdError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideImpl{
				id:          tt.fields.id,
				description: tt.fields.description,
				rootNode:    tt.fields.rootNode,
			}
			err := g.SetID(tt.args.id)
			if tt.wantErr == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			assert.Equal(t, tt.wantId, g.ID())
		})
	}
}

func Test_guideImpl_NodesToJSON(t *testing.T) {
	tests := []struct {
		name    string
		guide   *guideImpl
		want    string
		wantErr error
	}{
		{
			name: "Should return empty json for emtpy nodes",
			guide: &guideImpl{
				id:          10,
				description: "",
				rootNode: &node{
					Condition: nil,
					Text:      "",
					NextNodes: nil,
				},
			},
			want:    "{}",
			wantErr: nil,
		},
		{
			name: "Should return json for difficult guide",
			guide: &guideImpl{
				id:          10,
				description: "test description",
				rootNode: &node{
					Condition: NewManualCondition(),
					Text:      "node_1_text",
					NextNodes: []*node{
						{
							Condition: NewTimeCondition(oneMinuteDuration),
							Text:      "inner_node_1_text",
							NextNodes: []*node{},
						},
					},
				},
			},
			want:    "{\"condition\":{\"type\":\"MANUAL\"},\"text\":\"node_1_text\",\"next_nodes\":[{\"condition\":{\"type\":\"TIME\",\"duration\":60000000000},\"text\":\"inner_node_1_text\"}]}",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.guide.NodesToJSON()

			if tt.wantErr == nil {
				assert.Equal(t, tt.want, got)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func Test_guideImpl_SetDescription(t *testing.T) {
	type fields struct {
		id          uint
		description string
		rootNode    *node
	}
	type args struct {
		description string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Should set new description",
			fields: fields{10, "old_desc", &node{}},
			args:   args{"new_desc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideImpl{
				id:          tt.fields.id,
				description: tt.fields.description,
				rootNode:    tt.fields.rootNode,
			}
			g.SetDescription(tt.args.description)
			assert.Equal(t, tt.args.description, g.Description())
		})
	}
}

func Test_guideImpl_SetNodesFromJSON(t *testing.T) {
	type args struct {
		nodesJson string
	}
	tests := []struct {
		name         string
		args         args
		wantRootNode *node
		wantErr      error
	}{
		{
			name: "Should set new nodes from json",
			args: args{nodesJson: "{\"condition\":{\"type\":\"MANUAL\"},\"text\":\"node_1_text\",\"next_nodes\":[{\"condition\":{\"type\":\"TIME\",\"duration\":60000000000},\"text\":\"inner_node_1_text\"}]}"},
			wantRootNode: &node{
				Condition: NewManualCondition(),
				Text:      "node_1_text",
				NextNodes: []*node{
					{
						Condition: NewTimeCondition(oneMinuteDuration),
						Text:      "inner_node_1_text",
						NextNodes: []*node{},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:         "Should return error for invalid JSON",
			args:         args{nodesJson: "abcdv"},
			wantRootNode: nil,
			wantErr:      &InvalidJsonError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideImpl{
				id:          10,
				description: "",
				rootNode:    nil,
			}
			err := g.SetNodesFromJSON(tt.args.nodesJson)

			if tt.wantErr == nil {
				assert.Nil(t, err)
				checkNodes(t, tt.wantRootNode, g.rootNode)
			} else {
				assert.Nil(t, g.rootNode)
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func checkNodes(t *testing.T, wantRootNode *node, gotRootNode *node) {
	// Check nodes
	stackWant := []*node{wantRootNode}
	stackGot := []*node{gotRootNode}
	for len(stackWant) > 0 {
		nodeWant := stackWant[len(stackWant)-1]
		stackWant = stackWant[:len(stackWant)-1]

		nodeGot := stackGot[len(stackGot)-1]
		stackGot = stackGot[:len(stackGot)-1]

		assert.Equal(t, nodeWant.Text, nodeGot.Text)
		if nodeWant.Condition != nil {
			assert.Equal(t, nodeWant.Condition.CondType, nodeGot.Condition.CondType)
			assert.Equal(t, nodeWant.Condition.Duration, nodeWant.Condition.Duration)
		}

		stackWant = append(stackWant, nodeWant.NextNodes...)
		stackGot = append(stackGot, nodeGot.NextNodes...)
	}
	assert.Equal(t, 0, len(stackGot))
}
