package service

import (
	"context"
	"errors"
	"testing"

	"github.com/s21platform/society-service/internal/model"

	logger_lib "github.com/s21platform/logger-lib"

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

func TestServer_GetSocietyInfo(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	s := &Server{dbR: mockDBRepo}

	t.Run("should_get_society_info_successfully", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		societyUUID := uuid.Generate().String()
		mockInput := &society.GetSocietyInfoIn{SocietyUUID: societyUUID}

		expectedSocietyInfo := &model.SocietyInfo{
			Name:           "Test Society",
			Description:    "A test society",
			OwnerUUID:      uuid.Generate().String(),
			PhotoURL:       "https://example.com/photo.jpg",
			FormatID:       1,
			PostPermission: 2,
			IsSearch:       true,
			CountSubscribe: 100,
			TagsID:         []int64{1, 2},
		}

		mockDBRepo.EXPECT().GetSocietyInfo(societyUUID).Return(expectedSocietyInfo, nil)
		mockLogger.EXPECT().AddFuncName("GetSocietyInfo")

		result, err := s.GetSocietyInfo(ctx, mockInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedSocietyInfo.Name, result.Name)
		assert.Equal(t, expectedSocietyInfo.Description, result.Description)
		assert.Equal(t, expectedSocietyInfo.OwnerUUID, result.OwnerUUID)
		assert.Equal(t, expectedSocietyInfo.PhotoURL, result.PhotoURL)
		assert.Equal(t, expectedSocietyInfo.FormatID, result.FormatID)
		assert.Equal(t, expectedSocietyInfo.PostPermission, result.PostPermission)
		assert.Equal(t, expectedSocietyInfo.IsSearch, result.IsSearch)
		assert.Equal(t, expectedSocietyInfo.CountSubscribe, result.CountSubscribe)
		assert.Len(t, result.TagsID, len(expectedSocietyInfo.TagsID))
	})

	t.Run("should_return_error_if_societyUUID_is_empty", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		mockInput := &society.GetSocietyInfoIn{SocietyUUID: ""}

		mockLogger.EXPECT().AddFuncName("GetSocietyInfo")
		mockLogger.EXPECT().Error("failed to SocietyUUID is empty")

		result, err := s.GetSocietyInfo(ctx, mockInput)

		assert.Nil(t, result)
		assert.Error(t, err)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
		assert.Equal(t, "societyUUID not provided", statusErr.Message())
	})

	t.Run("should_return_error_if_dbR_GetSocietyInfo_fails", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		societyUUID := uuid.Generate().String()
		mockInput := &society.GetSocietyInfoIn{SocietyUUID: societyUUID}
		expectedError := errors.New("database error")

		mockDBRepo.EXPECT().GetSocietyInfo(societyUUID).Return(nil, expectedError)
		mockLogger.EXPECT().AddFuncName("GetSocietyInfo")
		mockLogger.EXPECT().Error("failed to GetSocietyInfo from BD")

		result, err := s.GetSocietyInfo(ctx, mockInput)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestServer_UpdateSociety(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	s := &Server{dbR: mockDBRepo}
	t.Run("should_update_society_info_successfully", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		societyUUID := uuid.Generate().String()
		ownerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, ownerUUID)

		expectedUpdateSociety := &society.UpdateSocietyIn{
			SocietyUUID:    societyUUID,
			Name:           "Test1",
			Description:    "A test society",
			PhotoURL:       "https://example.com/photo.jpg",
			FormatID:       1,
			IsSearch:       true,
			PostPermission: 2,
			TagsID: []*society.TagsID{
				{TagID: 1},
				{TagID: 2},
			},
		}

		mockLogger.EXPECT().AddFuncName("UpdateSociety")
		mockDBRepo.EXPECT().UpdateSociety(expectedUpdateSociety, ownerUUID).Return(nil) // Исправлен порядок аргументов

		_, err := s.UpdateSociety(ctx, expectedUpdateSociety)

		assert.NoError(t, err)
	})
	t.Run("should_return_error_if_uuid_not_found_in_context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		expectedUpdateSociety := &society.UpdateSocietyIn{
			SocietyUUID:    uuid.Generate().String(),
			Name:           "Test1",
			Description:    "A test society",
			PhotoURL:       "https://example.com/photo.jpg",
			FormatID:       1,
			IsSearch:       true,
			PostPermission: 2,
			TagsID: []*society.TagsID{
				{TagID: 1},
				{TagID: 2},
			},
		}

		mockLogger.EXPECT().AddFuncName("UpdateSociety")
		mockLogger.EXPECT().Error("failed to not found UUID in context")

		_, err := s.UpdateSociety(ctx, expectedUpdateSociety)

		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("should_return_error_if_societyUUID_is_empty", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		ownerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, ownerUUID)

		expectedUpdateSociety := &society.UpdateSocietyIn{
			SocietyUUID:    "",
			Name:           "Test1",
			Description:    "A test society",
			PhotoURL:       "https://example.com/photo.jpg",
			FormatID:       1,
			IsSearch:       true,
			PostPermission: 2,
			TagsID: []*society.TagsID{
				{TagID: 1},
				{TagID: 2},
			},
		}

		mockLogger.EXPECT().AddFuncName("UpdateSociety")
		mockLogger.EXPECT().Error("failed to SocietyUUID is empty")

		_, err := s.UpdateSociety(ctx, expectedUpdateSociety)

		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("should_return_error_if_name_is_empty", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		ownerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, ownerUUID)

		expectedUpdateSociety := &society.UpdateSocietyIn{
			SocietyUUID:    uuid.Generate().String(),
			Name:           "",
			Description:    "A test society",
			PhotoURL:       "https://example.com/photo.jpg",
			FormatID:       1,
			IsSearch:       true,
			PostPermission: 2,
			TagsID: []*society.TagsID{
				{TagID: 1},
				{TagID: 2},
			},
		}

		mockLogger.EXPECT().AddFuncName("UpdateSociety")
		mockLogger.EXPECT().Error("failed to Name society is empty")

		_, err := s.UpdateSociety(ctx, expectedUpdateSociety)

		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("should_return_error_if_repo_update_fails", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		societyUUID := uuid.Generate().String()
		ownerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, ownerUUID)

		expectedUpdateSociety := &society.UpdateSocietyIn{
			SocietyUUID:    societyUUID,
			Name:           "Test1",
			Description:    "A test society",
			PhotoURL:       "https://example.com/photo.jpg",
			FormatID:       1,
			IsSearch:       true,
			PostPermission: 2,
			TagsID: []*society.TagsID{
				{TagID: 1},
				{TagID: 2},
			},
		}

		expectedError := errors.New("database error")
		mockLogger.EXPECT().AddFuncName("UpdateSociety")
		mockLogger.EXPECT().Error("failed to UpdateSociety from BD")
		mockDBRepo.EXPECT().UpdateSociety(expectedUpdateSociety, ownerUUID).Return(expectedError)

		_, err := s.UpdateSociety(ctx, expectedUpdateSociety)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}
