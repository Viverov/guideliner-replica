// +build unit

package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGuideDTOFromEntity(t *testing.T) {
	type args struct {
		guide Guide
	}
	tests := []struct {
		name    string
		args    args
		want    *guideDTOImpl
		wantErr error
	}{
		{
			name: "Should return DTO",
			args: args{&guideImpl{
				id:          10,
				description: "test",
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
				creatorID: 10,
			}},
			want: &guideDTOImpl{
				id:          10,
				description: "test",
				nodesJSON:   "{\"condition\":{\"type\":\"MANUAL\"},\"text\":\"node_1_text\",\"next_nodes\":[{\"condition\":{\"type\":\"TIME\",\"duration\":60000000000},\"text\":\"inner_node_1_text\"}]}",
				creatorID:   10,
			},
			wantErr: nil,
		},
		{
			name: "Should return DTO for empty nodes",
			args: args{&guideImpl{
				id:          10,
				description: "test",
				rootNode:    &node{},
			}},
			want: &guideDTOImpl{
				id:          10,
				description: "test",
				nodesJSON:   "{}",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGuideDTOFromEntity(tt.args.guide)
			if tt.wantErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, tt.want.ID(), got.ID())
				assert.Equal(t, tt.want.Description(), got.Description())
				assert.Equal(t, tt.want.NodesJson(), got.NodesJson())
				assert.Equal(t, tt.want.CreatorID(), got.CreatorID())
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func Test_guideDTOImpl_Description(t *testing.T) {
	type fields struct {
		id          uint
		description string
		nodesJson   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return description",
			fields: fields{
				id:          10,
				description: "test",
				nodesJson:   "{}",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideDTOImpl{
				id:          tt.fields.id,
				description: tt.fields.description,
				nodesJSON:   tt.fields.nodesJson,
			}
			assert.Equal(t, tt.want, g.Description())
		})
	}
}

func Test_guideDTOImpl_ID(t *testing.T) {
	type fields struct {
		id          uint
		description string
		nodesJson   string
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name: "Should return ID",
			fields: fields{
				id:          10,
				description: "test",
				nodesJson:   "{}",
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideDTOImpl{
				id:          tt.fields.id,
				description: tt.fields.description,
				nodesJSON:   tt.fields.nodesJson,
			}
			assert.Equal(t, tt.want, g.ID())
		})
	}
}

func Test_guideDTOImpl_NodesJson(t *testing.T) {
	type fields struct {
		id          uint
		description string
		nodesJson   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return nodes json",
			fields: fields{
				id:          10,
				description: "",
				nodesJson:   "{a:2, b:3}",
			},
			want: "{a:2, b:3}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideDTOImpl{
				id:          tt.fields.id,
				description: tt.fields.description,
				nodesJSON:   tt.fields.nodesJson,
			}
			assert.Equal(t, tt.want, g.NodesJson())
		})
	}
}

func Test_guideDTOImpl_CreatorID(t *testing.T) {
	type fields struct {
		id          uint
		description string
		nodesJson   string
		creatorID   uint
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name: "Should return creatorID",
			fields: fields{
				id:          10,
				description: "",
				nodesJson:   "{a:2, b:3}",
				creatorID:   50,
			},
			want: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &guideDTOImpl{
				id:          tt.fields.id,
				description: tt.fields.description,
				nodesJSON:   tt.fields.nodesJson,
				creatorID:   tt.fields.creatorID,
			}
			assert.Equal(t, tt.want, g.CreatorID())
		})
	}
}
