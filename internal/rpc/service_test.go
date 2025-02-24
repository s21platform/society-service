package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/distribution/uuid"
	"github.com/golang/mock/gomock"
	society "github.com/s21platform/society-proto/society-proto"
	"github.com/s21platform/society-service/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)

	server := New(mockDBRepo)

	assert.NotNil(t, server)
	assert.Equal(t, mockDBRepo, server.dbR)
}

func TestServer_CreateSociety(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)

	s := &Server{dbR: mockDBRepo}
	t.Run("should_create_society_successfully", func(t *testing.T) {
		userUUID := uuid.Generate().String()
		ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)

		mockInput := &society.SetSocietyIn{
			Name:             "Test Society",
			FormatID:         1,
			PostPermissionID: 2,
			IsSearch:         true,
		}
		expectedSocietyUUID := uuid.Generate().String()
		mockDBRepo.EXPECT().CreateSociety(gomock.Any()).Return(expectedSocietyUUID, nil)

		result, err := s.CreateSociety(ctx, mockInput)
		expectedOutput := &society.SetSocietyOut{SocietyUUID: expectedSocietyUUID}

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, result)
	})
	t.Run("should_return_error_if_uuid_not_found_in_context", func(t *testing.T) {
		ctx := context.Background()
		mockInput := &society.SetSocietyIn{
			Name:             "Test Society",
			FormatID:         1,
			PostPermissionID: 2,
			IsSearch:         true,
		}
		result, err := s.CreateSociety(ctx, mockInput)
		assert.Nil(t, result)
		assert.Error(t, err)
	})
	t.Run("should_return_error_if_dbR_CreateSociety_fails", func(t *testing.T) {
		userUUID := uuid.Generate().String()
		ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)
		mockInput := &society.SetSocietyIn{
			Name:             "Test Society",
			FormatID:         1,
			PostPermissionID: 2,
			IsSearch:         true,
		}
		expectedError := errors.New("database error")
		mockDBRepo.EXPECT().CreateSociety(gomock.Any()).Return("", expectedError)
		result, err := s.CreateSociety(ctx, mockInput)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})
}
