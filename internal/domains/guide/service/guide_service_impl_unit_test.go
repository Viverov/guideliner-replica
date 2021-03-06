// +build unit

package service

import (
	"github.com/Viverov/guideliner/internal/domains/guide/entity"
	"github.com/Viverov/guideliner/internal/domains/guide/repository"
	"github.com/Viverov/guideliner/internal/domains/guide/service/mocks"
	"github.com/Viverov/guideliner/internal/domains/util"
	"github.com/Viverov/guideliner/internal/domains/util/urepo"
	"github.com/Viverov/guideliner/internal/domains/util/uservice"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewGuideService(t *testing.T) {
	tests := []struct {
		name string
		want *guideServiceImpl
	}{
		{
			name: "Should create new guide service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, mock := prepareMocks(t)

			got := NewGuideService(mock)
			assert.NotNil(t, got)

			ctrl.Finish()
		})
	}
}

func Test_guideServiceImpl_Create(t *testing.T) {
	type args struct {
		description string
		nodesJson   string
		creatorID   uint
	}
	type resFromRepo struct {
		id  uint
		err error
	}
	tests := []struct {
		name        string
		args        args
		want        entity.GuideDTO
		wantErr     error
		resFromRepo resFromRepo
	}{
		{
			name: "Should create new guide",
			args: args{
				description: "Some guide",
				nodesJson:   "{}",
				creatorID:   10,
			},
			want:        entity.NewGuideDTO(0, "Some guide", "{}", 10),
			wantErr:     nil,
			resFromRepo: resFromRepo{id: 10, err: nil},
		},
		{
			name: "Should return error on repository error",
			args: args{
				description: "Some guide",
				nodesJson:   "{}",
				creatorID:   50,
			},
			want:        nil,
			wantErr:     uservice.NewStorageError(urepo.NewUnexpectedRepositoryError("test", "text").Error()),
			resFromRepo: resFromRepo{id: 10, err: urepo.NewUnexpectedRepositoryError("test", "text")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, rep := prepareMocks(t)

			s := &guideServiceImpl{
				repository: rep,
			}

			rep.
				EXPECT().
				Insert(gomock.Any()).
				Return(tt.resFromRepo.id, tt.resFromRepo.err)

			got, err := s.Create(tt.args.description, tt.args.nodesJson, tt.args.creatorID)

			if tt.want != nil {
				assert.Equal(t, tt.resFromRepo.id, got.ID())
				assert.Equal(t, tt.want.Description(), got.Description())
				assert.Equal(t, tt.want.NodesJson(), got.NodesJson())
				assert.Equal(t, tt.want.CreatorID(), got.CreatorID())
			} else {
				assert.Nil(t, got)
			}

			if tt.wantErr != nil {
				assert.Error(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, nil)
			}

			ctrl.Finish()
		})
	}
}

func Test_guideServiceImpl_Find(t *testing.T) {
	type args struct {
		cond FindConditions
	}
	type resFromRep struct {
		guides []entity.Guide
		err    error
	}
	tests := []struct {
		name       string
		args       args
		resFromRep resFromRep
		want       []entity.GuideDTO
		wantErr    error
	}{
		{
			name: "Should find guides without conditions",
			args: args{
				cond: FindConditions{},
			},
			resFromRep: resFromRep{
				guides: func() []entity.Guide {
					var guides []entity.Guide
					for i := 0; i < 5; i++ {
						g, _ := entity.NewGuide(uint(10+i), "{}", "test"+strconv.Itoa(i), uint(50+i))
						guides = append(guides, g)
					}
					return guides
				}(),
				err: nil,
			},
			want: func() []entity.GuideDTO {
				var dtos []entity.GuideDTO
				for i := 0; i < 5; i++ {
					dto := entity.NewGuideDTO(uint(10+i), "test"+strconv.Itoa(i), "{}", uint(50+i))
					dtos = append(dtos, dto)
				}
				return dtos
			}(),
			wantErr: nil,
		},
		{
			name: "Should pass conditions into repository",
			args: args{
				cond: FindConditions{
					DefaultFindConditions: util.DefaultFindConditions{
						Limit:  12,
						Offset: 18,
						Order: []util.Order{
							{
								Field:     "id",
								Direction: "asc",
							},
						},
					},
					Search: "testing",
				},
			},
			resFromRep: resFromRep{
				guides: []entity.Guide{},
				err:    nil,
			},
			want:    []entity.GuideDTO{},
			wantErr: nil,
		},
		{
			name: "Should process error from repository",
			args: args{
				cond: FindConditions{},
			},
			resFromRep: resFromRep{
				guides: nil,
				err:    urepo.NewUnexpectedRepositoryError("test", "text"),
			},
			want:    nil,
			wantErr: uservice.NewStorageError(urepo.NewUnexpectedRepositoryError("test", "text").Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, rep := prepareMocks(t)
			s := &guideServiceImpl{
				repository: rep,
			}

			rep.
				EXPECT().
				Find(gomock.Eq(repository.FindConditions{
					DefaultFindConditions: tt.args.cond.DefaultFindConditions,
					Search:                tt.args.cond.Search,
				})).Return(tt.resFromRep.guides, tt.resFromRep.err)

			got, err := s.Find(tt.args.cond)

			if tt.want != nil {
				assert.Equal(t, len(got), len(tt.want))
				for i, dto := range got {
					assert.Equal(t, tt.want[i].ID(), dto.ID())
					assert.Equal(t, tt.want[i].Description(), dto.Description())
					assert.Equal(t, tt.want[i].NodesJson(), dto.NodesJson())
					assert.Equal(t, tt.want[i].CreatorID(), dto.CreatorID())
				}
			} else {
				assert.Nil(t, got)
			}

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			ctrl.Finish()
		})
	}
}

func Test_guideServiceImpl_FindById(t *testing.T) {
	type args struct {
		id uint
	}
	type resFromRep struct {
		entity entity.Guide
		err    error
	}
	tests := []struct {
		name       string
		args       args
		resFromRep resFromRep
		want       entity.GuideDTO
		wantErr    error
	}{
		{
			name: "Should find user",
			args: args{
				id: 10,
			},
			resFromRep: resFromRep{
				entity: func() entity.Guide { g, _ := entity.NewGuide(10, "{}", "testdesc", 50); return g }(),
				err:    nil,
			},
			want:    entity.NewGuideDTO(10, "testdesc", "{}", 50),
			wantErr: nil,
		},
		{
			name: "Should return nil on undefined user",
			args: args{
				id: 10,
			},
			resFromRep: resFromRep{
				entity: nil,
				err:    nil,
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "Should return error on repository error",
			args: args{
				id: 10,
			},
			want:    nil,
			wantErr: uservice.NewStorageError(urepo.NewUnexpectedRepositoryError("text", "text").Error()),
			resFromRep: resFromRep{
				entity: nil,
				err:    urepo.NewUnexpectedRepositoryError("text", "text"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, rep := prepareMocks(t)
			s := &guideServiceImpl{
				repository: rep,
			}

			rep.
				EXPECT().
				FindById(gomock.Eq(tt.args.id)).
				Return(tt.resFromRep.entity, tt.resFromRep.err)

			got, err := s.FindById(tt.args.id)

			if tt.want != nil {
				assert.Equal(t, tt.want.ID(), got.ID())
				assert.Equal(t, tt.want.Description(), got.Description())
				assert.Equal(t, tt.want.NodesJson(), got.NodesJson())
				assert.Equal(t, tt.want.CreatorID(), got.CreatorID())
			} else {
				assert.Nil(t, got)
			}

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			ctrl.Finish()
		})
	}
}

func Test_guideServiceImpl_Count(t *testing.T) {
	type args struct {
		cond CountConditions
	}
	type resFromRep struct {
		count int64
		err   error
	}
	tests := []struct {
		name       string
		args       args
		resFromRep resFromRep
		want       int64
		wantErr    error
	}{
		{
			name: "Should return count without conditions",
			args: args{
				cond: CountConditions{},
			},
			resFromRep: resFromRep{
				count: 10,
				err:   nil,
			},
			want:    10,
			wantErr: nil,
		},
		{
			name: "Should pass conditions into repository",
			args: args{
				cond: CountConditions{
					Search: "testing",
				},
			},
			resFromRep: resFromRep{
				count: 10,
				err:   nil,
			},
			want:    10,
			wantErr: nil,
		},
		{
			name: "Should process error from repository",
			args: args{
				cond: CountConditions{},
			},
			resFromRep: resFromRep{
				count: 0,
				err:   urepo.NewUnexpectedRepositoryError("test", "text"),
			},
			want:    0,
			wantErr: uservice.NewStorageError(urepo.NewUnexpectedRepositoryError("test", "text").Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, rep := prepareMocks(t)
			s := &guideServiceImpl{
				repository: rep,
			}

			rep.
				EXPECT().
				Count(gomock.Eq(repository.CountConditions{
					Search: tt.args.cond.Search,
				})).Return(tt.resFromRep.count, tt.resFromRep.err)

			got, err := s.Count(tt.args.cond)

			assert.Equal(t, tt.want, got)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			ctrl.Finish()
		})
	}
}

func Test_guideServiceImpl_Update(t *testing.T) {
	type args struct {
		id     uint
		params UpdateParams
	}
	type resFromRepOnFind struct {
		guide entity.Guide
		err   error
	}
	type resFromRepOnUpdate struct {
		err error
	}
	tests := []struct {
		name               string
		args               args
		resFromRepOnFind   resFromRepOnFind
		repUpdateExpected  bool
		resFromRepOnUpdate resFromRepOnUpdate
		want               entity.GuideDTO
		wantErr            error
	}{
		{
			name: "Should update entity",
			args: args{
				id: 10,
				params: UpdateParams{
					Description: "newDesc",
					NodesJson:   "{}",
				},
			},
			resFromRepOnFind: resFromRepOnFind{
				guide: func() entity.Guide { g, _ := entity.NewGuide(10, "{}", "desc", 50); return g }(),
				err:   nil,
			},
			repUpdateExpected: true,
			resFromRepOnUpdate: resFromRepOnUpdate{
				err: nil,
			},
			want:    entity.NewGuideDTO(10, "newDesc", "{}", 50),
			wantErr: nil,
		},
		{
			name: "Should return error on undefined user",
			args: args{
				id: 10,
				params: UpdateParams{
					Description: "newDesc",
					NodesJson:   "{}",
				},
			},
			resFromRepOnFind: resFromRepOnFind{
				guide: nil,
				err:   nil,
			},
			repUpdateExpected:  false,
			resFromRepOnUpdate: resFromRepOnUpdate{},
			want:               nil,
			wantErr:            uservice.NewNotFoundError("Guide", 10),
		},
		{
			name: "Should return error on repository error",
			args: args{
				id: 10,
				params: UpdateParams{
					Description: "newDesc",
					NodesJson:   "{}",
				},
			},
			resFromRepOnFind: resFromRepOnFind{
				guide: func() entity.Guide { g, _ := entity.NewGuide(10, "{}", "desc", 50); return g }(),
				err:   nil,
			},
			repUpdateExpected: true,
			resFromRepOnUpdate: resFromRepOnUpdate{
				err: urepo.NewUnexpectedRepositoryError("test", "text"),
			},
			want:    nil,
			wantErr: uservice.NewStorageError(urepo.NewUnexpectedRepositoryError("test", "text").Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, rep := prepareMocks(t)
			s := &guideServiceImpl{
				repository: rep,
			}

			rep.
				EXPECT().
				FindById(gomock.Eq(tt.args.id)).
				Return(tt.resFromRepOnFind.guide, tt.resFromRepOnFind.err)

			if tt.repUpdateExpected {
				rep.
					EXPECT().
					Update(gomock.Eq(tt.resFromRepOnFind.guide)).
					Return(tt.resFromRepOnUpdate.err)
			}

			got, err := s.Update(tt.args.id, tt.args.params)

			if tt.want != nil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID(), got.ID())
				assert.Equal(t, tt.want.Description(), got.Description())
				assert.Equal(t, tt.want.NodesJson(), got.NodesJson())
				assert.Equal(t, tt.want.CreatorID(), got.CreatorID())
			} else {
				assert.Nil(t, got)
			}

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			ctrl.Finish()
		})
	}
}

func Test_guideServiceImpl_CheckPermission(t *testing.T) {
	type args struct {
		guideID    uint
		userID     uint
		permission Permission
	}
	type resFromRepoOnFind struct {
		guide entity.Guide
		err   error
	}
	tests := []struct {
		name              string
		args              args
		resFromRepoOnFind resFromRepoOnFind
		want              bool
		wantErr           error
	}{
		{
			name: "Should return true for valid permission",
			args: args{
				guideID:    10,
				userID:     30,
				permission: PermissionUpdate,
			},
			resFromRepoOnFind: resFromRepoOnFind{
				guide: func() entity.Guide { g, _ := entity.NewGuide(10, "{}", "desc", 30); return g }(),
				err:   nil,
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "Should return false for invalid permission",
			args: args{
				guideID:    10,
				userID:     50,
				permission: PermissionUpdate,
			},
			resFromRepoOnFind: resFromRepoOnFind{
				guide: func() entity.Guide { g, _ := entity.NewGuide(10, "{}", "desc", 30); return g }(),
				err:   nil,
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "Should return not found error for undefined guide",
			args: args{
				guideID:    10,
				userID:     50,
				permission: PermissionUpdate,
			},
			resFromRepoOnFind: resFromRepoOnFind{
				guide: nil,
				err:   nil,
			},
			want:    false,
			wantErr: uservice.NewNotFoundError("Guide", 10),
		},
	}
	for _, tt := range tests {
		ctrl, rep := prepareMocks(t)
		s := &guideServiceImpl{
			repository: rep,
		}

		rep.
			EXPECT().
			FindById(gomock.Eq(tt.args.guideID)).
			Return(tt.resFromRepoOnFind.guide, tt.resFromRepoOnFind.err)

		got, err := s.CheckPermission(tt.args.guideID, tt.args.userID, tt.args.permission)

		assert.Equal(t, tt.want, got)

		if tt.wantErr != nil {
			assert.EqualError(t, err, tt.wantErr.Error())
		} else {
			assert.Nil(t, err)
		}

		ctrl.Finish()
	}
}

func Test_guideServiceImpl_GetPermissions(t *testing.T) {
	type args struct {
		guideID uint
		userID  uint
	}
	type resFromRepoOnFind struct {
		guide entity.Guide
		err   error
	}
	tests := []struct {
		name              string
		args              args
		resFromRepoOnFind resFromRepoOnFind
		want              []Permission
		wantErr           error
	}{
		{
			name: "Should return all permissions for creator of guide",
			args: args{
				guideID: 10,
				userID:  30,
			},
			resFromRepoOnFind: resFromRepoOnFind{
				guide: func() entity.Guide { g, _ := entity.NewGuide(10, "{}", "desc", 30); return g }(),
				err:   nil,
			},
			want:    []Permission{PermissionUpdate},
			wantErr: nil,
		},
		{
			name: "Should return empty permissions for undefined user",
			args: args{
				guideID: 10,
				userID:  50,
			},
			resFromRepoOnFind: resFromRepoOnFind{
				guide: func() entity.Guide { g, _ := entity.NewGuide(10, "{}", "desc", 30); return g }(),
				err:   nil,
			},
			want:    []Permission{},
			wantErr: nil,
		},
		{
			name: "Should return not found error for undefined guide",
			args: args{
				guideID: 10,
				userID:  50,
			},
			resFromRepoOnFind: resFromRepoOnFind{
				guide: nil,
				err:   nil,
			},
			want:    nil,
			wantErr: uservice.NewNotFoundError("Guide", 10),
		},
	}
	for _, tt := range tests {
		ctrl, rep := prepareMocks(t)
		s := &guideServiceImpl{
			repository: rep,
		}

		rep.
			EXPECT().
			FindById(gomock.Eq(tt.args.guideID)).
			Return(tt.resFromRepoOnFind.guide, tt.resFromRepoOnFind.err)

		got, err := s.GetPermissions(tt.args.guideID, tt.args.userID)

		if tt.want != nil {
			assert.NotNil(t, got)
			assert.Equal(t, len(tt.want), len(got))
			for i, expected := range tt.want {
				assert.Equal(t, expected, got[i])
			}
		} else {
			assert.Nil(t, got)
		}

		if tt.wantErr != nil {
			assert.EqualError(t, err, tt.wantErr.Error())
		} else {
			assert.Nil(t, err)
		}

		ctrl.Finish()
	}
}

func prepareMocks(t *testing.T) (
	guideRepositoryCtrl *gomock.Controller,
	guideRepositoryMock *mocks.MockGuideRepository,
) {
	guideRepositoryCtrl = gomock.NewController(t)
	guideRepositoryMock = mocks.NewMockGuideRepository(guideRepositoryCtrl)

	return guideRepositoryCtrl, guideRepositoryMock
}
