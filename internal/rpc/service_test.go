package rpc_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	society "github.com/s21platform/society-proto/society-proto"
	"github.com/s21platform/society-service/internal/config"
	"github.com/s21platform/society-service/internal/model"
	"github.com/s21platform/society-service/internal/rpc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestServer_GetSocietiesForUser(t *testing.T) {
	t.Parallel()

	uuid := "test-uuid"
	ctx := context.WithValue(context.Background(), config.KeyUUID, uuid)
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := rpc.NewMockDbRepo(controller)

	userUuid := "user-uuid"
	expectedSocieties := []model.SocietyWithOffsetData{
		{
			Name:       "Society 1",
			AvatarLink: "https://example.com/avatar1",
			SocietyId:  1,
			IsMember:   true,
			IsPrivate:  false,
		},
		{
			Name:       "Society 2",
			AvatarLink: "https://example.com/avatar2",
			SocietyId:  2,
			IsMember:   false,
			IsPrivate:  true,
		},
	}
	mockRepo.EXPECT().GetSocietiesForUser(uuid, userUuid).Return(&expectedSocieties, nil)

	s := rpc.New(mockRepo)

	t.Run("get_societies_for_user_ok", func(t *testing.T) {
		in := &society.GetSocietiesForUserIn{UserUuid: userUuid}
		out, err := s.GetSocietiesForUser(ctx, in)
		assert.NoError(t, err)
		expectedOut := &society.GetSocietiesForUserOut{
			Society: []*society.Society{
				{
					Name:       "Society 1",
					AvatarLink: "https://example.com/avatar1",
					SocietyId:  1,
					IsMember:   true,
					IsPrivate:  false,
				},
				{
					Name:       "Society 2",
					AvatarLink: "https://example.com/avatar2",
					SocietyId:  2,
					IsMember:   false,
					IsPrivate:  true,
				},
			},
		}
		assert.Equal(t, expectedOut, out)
	})

	t.Run("uuid_not_in_context", func(t *testing.T) {
		ctx := context.Background()
		in := &society.GetSocietiesForUserIn{UserUuid: userUuid}
		_, err := s.GetSocietiesForUser(ctx, in)
		assert.EqualError(t, err, "uuid not found in context")
	})

	t.Run("repository_error", func(t *testing.T) {
		mockRepo.EXPECT().GetSocietiesForUser(uuid, userUuid).Return(nil, fmt.Errorf("db error"))
		in := &society.GetSocietiesForUserIn{UserUuid: userUuid}
		_, err := s.GetSocietiesForUser(ctx, in)
		assert.EqualError(t, err, "failed to get society for user: db error")
	})
}

func TestServer_GetSocietyInfo(t *testing.T) {
	t.Parallel()

	uuid := "test-uuid"
	ctx := context.WithValue(context.Background(), config.KeyUUID, uuid)
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := rpc.NewMockDbRepo(controller)

	//userUuid := "user-uuid"
	expectedSocieties := model.SocietyInfo{
		Name:             "Society 1",
		Description:      "This is a test society 1",
		OwnerId:          "user-uuid",
		PhotoUrl:         "https://example.com/avatar1",
		IsPrivate:        true,
		CountSubscribers: 5,
	}
	mockRepo.EXPECT().GetSocietyInfo(int64(1)).Return(&expectedSocieties, nil)

	s := rpc.New(mockRepo)

	t.Run("get_societies_info_ok", func(t *testing.T) {
		in := &society.GetSocietyInfoIn{Id: 1}
		out, err := s.GetSocietyInfo(ctx, in)
		assert.NoError(t, err)
		expectedOut := &society.GetSocietyInfoOut{
			Name:             "Society 1",
			Description:      "This is a test society 1",
			OwnerUUID:        "user-uuid",
			PhotoUrl:         "https://example.com/avatar1",
			IsPrivate:        true,
			CountSubscribers: 5,
		}
		assert.Equal(t, expectedOut, out)
	})

	t.Run("repository_error", func(t *testing.T) {
		mockRepo.EXPECT().GetSocietyInfo(int64(1)).Return(nil, fmt.Errorf("db error"))
		in := &society.GetSocietyInfoIn{Id: 1}
		_, err := s.GetSocietyInfo(ctx, in)
		assert.EqualError(t, err, "failed to get society info: db error")
	})
}
