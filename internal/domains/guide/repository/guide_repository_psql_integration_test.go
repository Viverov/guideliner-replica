// +build integration

package repository

import (
	"errors"
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/Viverov/guideliner/internal/connectors"
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/util"
	"github.com/Viverov/guideliner/internal/domains/util/urepo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var cfg = config.InitConfig(config.EnvTest, "./config.json")
var dbInstance = connectors.GetDB(&connectors.DBOptions{
	Host:     cfg.DB.Host,
	Port:     cfg.DB.Port,
	Login:    cfg.DB.Login,
	Password: cfg.DB.Password,
	Name:     cfg.DB.Name,
	SSLMode:  cfg.DB.SSLMode,
})

type guideData struct {
	ID          uint
	Description string
	nodesJSON   string
}

func TestNewGuideRepositoryPsql(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should return new repository",
			args: args{
				db: &gorm.DB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGuideRepositoryPsql(tt.args.db)
			assert.NotNil(t, got)
		})
	}
}

func Test_guideRepositoryPsql_Find(t *testing.T) {
	// Setup test data (before all)
	testGuideData := prepareTestData()

	type args struct {
		condition FindConditions
	}
	tests := []struct {
		name    string
		args    args
		want    []guideData
		wantErr error
	}{
		{
			name: "Should return all records without conditions",
			args: args{
				FindConditions{
					DefaultFindConditions: util.DefaultFindConditions{
						Limit:  20,
						Offset: 0,
						Order:  []util.Order{{Field: "id", Direction: util.ASC}}},
				},
			},
			want:    testGuideData,
			wantErr: nil,
		},
		{
			name: "Should return records with limit condition",
			args: args{
				FindConditions{
					DefaultFindConditions: util.DefaultFindConditions{
						Limit:  5,
						Offset: 0,
						Order:  []util.Order{{Field: "id", Direction: util.ASC}}},
				},
			},
			want:    testGuideData[:5],
			wantErr: nil,
		},
		{
			name: "Should return records with limit and offset conditions",
			args: args{
				FindConditions{
					DefaultFindConditions: util.DefaultFindConditions{
						Limit:  5,
						Offset: 5,
						Order:  []util.Order{{Field: "id", Direction: util.ASC}}},
				},
			},
			want:    testGuideData[5:10],
			wantErr: nil,
		},
		{
			name: "Should reverse order",
			args: args{
				FindConditions{
					DefaultFindConditions: util.DefaultFindConditions{
						Limit: 20,
						Order: []util.Order{{Field: "id", Direction: util.DESC}}},
				},
			},
			want: func() []guideData {
				resultGuideData := []guideData{}
				for i := 0; i < len(testGuideData); i++ {
					resultGuideData = append([]guideData{testGuideData[i]}, resultGuideData...)
				}

				return resultGuideData
			}(),
			wantErr: nil,
		},

		{
			name: "Should return records by search field",
			args: args{
				FindConditions{
					DefaultFindConditions: util.DefaultFindConditions{},
					Search:                "est2", // We have only one record that meets the condition - "test2"
				},
			},
			want:    testGuideData[1:2],
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &guideRepositoryPsql{
				db: dbInstance,
			}
			guides, err := r.Find(tt.args.condition)

			if tt.want != nil {
				assert.Equal(t, len(tt.want), len(guides))
				for i, g := range guides {
					assert.Equal(t, tt.want[i].ID, g.ID())
					assert.Equal(t, tt.want[i].Description, g.Description())
					nj, _ := g.NodesToJSON()
					assert.Equal(t, tt.want[i].nodesJSON, nj)
				}
			} else {
				assert.Nil(t, tt.want)
			}

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}

	// Clean up (after all)
	cleanUpTestData()
}

func Test_guideRepositoryPsql_Count(t *testing.T) {
	// Setup test data (before all)
	testGuideData := prepareTestData()

	type args struct {
		condition CountConditions
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "Should return count without conditions",
			args: args{
				CountConditions{},
			},
			want:    int64(len(testGuideData)),
			wantErr: nil,
		},
		{
			name: "Should return count with 'search' condition",
			args: args{
				CountConditions{
					Search: "est1",
				},
			},
			want:    2, // We have two records that meets the condition - "test1" and "test10"
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &guideRepositoryPsql{
				db: dbInstance,
			}
			count, err := r.Count(tt.args.condition)

			assert.Equal(t, tt.want, count)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}

	// Clean up (after all)
	cleanUpTestData()
}

func Test_guideRepositoryPsql_FindById(t *testing.T) {
	// Prepare test data
	testGuideData := prepareTestData()

	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		args    args
		want    *guideData
		wantErr error
	}{
		{
			name: "Should return records by id",
			args: args{
				id: 5,
			},
			want:    &testGuideData[4],
			wantErr: nil,
		},
		{
			name: "Should return nil for undefined ID",
			args: args{
				id: 543,
			},
			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &guideRepositoryPsql{
				db: dbInstance,
			}

			got, err := r.FindById(tt.args.id)

			if tt.want != nil {
				assert.Equal(t, tt.want.ID, got.ID())
				assert.Equal(t, tt.want.Description, got.Description())
				nj, _ := got.NodesToJSON()
				assert.Equal(t, tt.want.nodesJSON, nj)
			} else {
				assert.Nil(t, tt.want)
			}

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}

	// Clean up (after all)
	cleanUpTestData()
}

func Test_guideRepositoryPsql_Insert(t *testing.T) {
	type args struct {
		guide entity.Guide
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Should create new record in db",
			args: args{
				guide: func() entity.Guide { g, _ := entity.NewGuide(0, "{}", "description", 50); return g }(),
			},
			wantErr: nil,
		},
		{
			name: "Should return error for nil entity",
			args: args{
				guide: nil,
			},
			wantErr: urepo.NewNilEntityError("Guide"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &guideRepositoryPsql{
				db: dbInstance,
			}
			gotId, err := r.Insert(tt.args.guide)

			gm := &guideModel{
				Model: gorm.Model{
					ID: gotId,
				},
			}
			resultFromDB := r.db.First(gm)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				assert.NotNil(t, resultFromDB.Error)
				assert.True(t, errors.Is(resultFromDB.Error, gorm.ErrRecordNotFound))
			} else {
				assert.Nil(t, err)
				assert.Nil(t, resultFromDB.Error)
				assert.NotNil(t, gm)
				assert.Equal(t, tt.args.guide.Description(), gm.Description)
				nj, err := tt.args.guide.NodesToJSON()
				assert.Nil(t, err)
				assert.Equal(t, nj, gm.NodesJson)
				assert.Equal(t, tt.args.guide.CreatorID(), gm.CreatorID)
			}
		})

		// Clean up (after each)
		cleanUpTestData()
	}
}

func Test_guideRepositoryPsql_Update(t *testing.T) {
	type args struct {
		guide entity.Guide
	}
	tests := []struct {
		name              string
		args              args
		updatedTestDataID uint
		wantErr           error
	}{
		{
			name: "Should update entity",
			args: args{
				guide: func() entity.Guide {
					g, _ := entity.NewGuide(5,
						"{\"condition\":{\"type\":\"MANUAL\"},\"text\":\"node_1_text\",\"next_nodes\":[{\"condition\":{\"type\":\"TIME\",\"duration\":60000000000},\"text\":\"inner_node_1_text\"}]}",
						"new description",
						50,
					)
					return g
				}(),
			},
			updatedTestDataID: 5,
			wantErr:           nil,
		},
		{
			name: "Should return error for undefined entity",
			args: args{
				guide: func() entity.Guide { g, _ := entity.NewGuide(55, "{}", "nvm", 50); return g }(),
			},
			updatedTestDataID: 0,
			wantErr:           urepo.NewEntityNotFoundError("Guide", 55),
		},
	}
	for _, tt := range tests {
		// Setup (before each)
		prepareTestData()

		t.Run(tt.name, func(t *testing.T) {
			r := &guideRepositoryPsql{
				db: dbInstance,
			}

			err := r.Update(tt.args.guide)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)

				gm := &guideModel{
					Model: gorm.Model{
						ID: tt.args.guide.ID(),
					},
				}
				resultFromDB := r.db.First(gm)
				assert.Nil(t, resultFromDB.Error)
				assert.Equal(t, tt.args.guide.Description(), gm.Description)
				nj, err := tt.args.guide.NodesToJSON()
				assert.Nil(t, err)
				assert.Equal(t, nj, gm.NodesJson)
			}
		})

		// Clean up (after each)
		cleanUpTestData()
	}
}

// prepareTestData create 10 records with id's from 1 to 10 and Descriptions from 'test1' to 'test10'
func prepareTestData() []guideData {
	var testGuideData []guideData
	for i := 1; i <= 10; i++ {
		testGuideData = append(testGuideData, guideData{
			ID:          uint(i),
			Description: "test" + fmt.Sprint(i),
			nodesJSON:   "{}",
		})
	}
	for _, gd := range testGuideData {
		dbInstance.Create(&guideModel{
			Model:       gorm.Model{ID: gd.ID},
			Description: gd.Description,
			NodesJson:   gd.nodesJSON,
		})
	}

	return testGuideData
}

// cleanUpTestData remove all guides from DB
func cleanUpTestData() {
	dbInstance.Unscoped().Where("1 = 1").Delete(&guideModel{})
}
