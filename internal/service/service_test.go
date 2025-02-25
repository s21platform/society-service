package service

import (
	"context"
	"errors"
	logger_lib "github.com/s21platform/logger-lib"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

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
		//mockLogger.EXPECT().AddFuncName("CreateSociety")
		mockDBRepo.EXPECT().CreateSociety(gomock.Any()).Return(expectedSocietyUUID, nil)
		mockLogger.EXPECT().AddFuncName("CreateSociety")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

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
		mockLogger.EXPECT().AddFuncName("CreateSociety")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		mockLogger.EXPECT().Error("failed to not found UUID in context")
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
		mockLogger.EXPECT().AddFuncName("CreateSociety")
		mockLogger.EXPECT().Error("failed to CreateSociety from BD")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		result, err := s.CreateSociety(ctx, mockInput)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})
	t.Run("should_return_error_if_name_is_empty", func(t *testing.T) {
		userUUID := uuid.Generate().String()
		ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)

		mockInput := &society.SetSocietyIn{
			Name:             "",
			FormatID:         1,
			PostPermissionID: 2,
			IsSearch:         true,
		}
		mockLogger.EXPECT().AddFuncName("CreateSociety")
		mockLogger.EXPECT().Error("failed to Name society is empty")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		result, err := s.CreateSociety(ctx, mockInput)
		statusErr, ok := status.FromError(err)

		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
		assert.Equal(t, "name not provided", statusErr.Message())

		assert.Nil(t, result)
	})
}
