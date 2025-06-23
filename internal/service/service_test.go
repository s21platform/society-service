package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

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
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockInput := &society.SetSocietyIn{
			Name:             "Test Society",
			FormatID:         1,
			PostPermissionID: 2,
			IsSearch:         true,
		}
		expectedSocietyUUID := uuid.Generate().String()
		mockDBRepo.EXPECT().CreateSociety(ctx, gomock.Any()).Return(expectedSocietyUUID, nil)
		mockLogger.EXPECT().AddFuncName("CreateSociety")

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
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockInput := &society.SetSocietyIn{
			Name:             "Test Society",
			FormatID:         1,
			PostPermissionID: 2,
			IsSearch:         true,
		}
		expectedError := errors.New("database error")
		mockDBRepo.EXPECT().CreateSociety(ctx, gomock.Any()).Return("", expectedError)
		mockLogger.EXPECT().AddFuncName("CreateSociety")
		mockLogger.EXPECT().Error("failed to CreateSociety from BD")

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
		userUUID := uuid.Generate().String()
		ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		societyUUID := uuid.Generate().String()

		mockInput := &society.GetSocietyInfoIn{SocietyUUID: societyUUID}

		expectedSocietyInfo := &model.SocietyInfo{
			Name:           "Test Society",
			Description:    sql.NullString{String: "A test society", Valid: true},
			OwnerUUID:      uuid.Generate().String(),
			PhotoURL:       "https://example.com/photo.jpg",
			FormatID:       1,
			PostPermission: 2,
			IsSearch:       true,
			CountSubscribe: 100,
			TagsID:         []int64{1, 2},
			CanEditSociety: true,
		}

		expectedCountSubscribe := int64(150)
		expectedTags := []int64{1, 2}
		expectedCanEdit := 1

		mockDBRepo.EXPECT().GetSocietyInfo(ctx, societyUUID).Return(expectedSocietyInfo, nil)
		mockDBRepo.EXPECT().CountSubscribe(ctx, societyUUID).Return(expectedCountSubscribe, nil)
		mockDBRepo.EXPECT().IsOwnerAdminModerator(ctx, userUUID, societyUUID).Return(expectedCanEdit, nil) // <-- Добавлено
		mockDBRepo.EXPECT().GetTags(ctx, societyUUID).Return(expectedTags, nil)

		mockLogger.EXPECT().AddFuncName("GetSocietyInfo")

		result, err := s.GetSocietyInfo(ctx, mockInput)

		assert.NoError(t, err)
		assert.Equal(t, expectedSocietyInfo.Name, result.Name)
		assert.Equal(t, expectedSocietyInfo.Description.String, result.Description)
		assert.Equal(t, expectedSocietyInfo.OwnerUUID, result.OwnerUUID)
		assert.Equal(t, expectedSocietyInfo.PhotoURL, result.PhotoURL)
		assert.Equal(t, expectedSocietyInfo.FormatID, result.FormatID)
		assert.Equal(t, expectedSocietyInfo.PostPermission, result.PostPermission)
		assert.Equal(t, expectedSocietyInfo.IsSearch, result.IsSearch)

		assert.Equal(t, expectedCountSubscribe, result.CountSubscribe)

		assert.Len(t, result.TagsID, len(expectedTags))
		for i, tag := range result.TagsID {
			assert.Equal(t, expectedTags[i], tag.TagID)
		}

		assert.True(t, result.CanEditSociety) // или false, если expectedCanEdit > 3
	})

	t.Run("should_return_error_if_societyUUID_is_empty", func(t *testing.T) {
		userUUID := uuid.Generate().String()
		ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
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
		userUUID := uuid.Generate().String()
		ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		societyUUID := uuid.Generate().String()
		mockInput := &society.GetSocietyInfoIn{SocietyUUID: societyUUID}
		expectedError := errors.New("database error")

		mockDBRepo.EXPECT().GetSocietyInfo(ctx, societyUUID).Return(nil, expectedError)
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
			FormatID:       1,
			IsSearch:       true,
			PostPermission: 2,
			TagsID: []*society.TagsID{
				{TagID: 1},
				{TagID: 2},
			},
		}

		mockLogger.EXPECT().AddFuncName("UpdateSociety")
		mockDBRepo.EXPECT().IsOwnerAdminModerator(ctx, ownerUUID, societyUUID).Return(1, nil) // 1 - Owner
		mockDBRepo.EXPECT().UpdateSociety(ctx, expectedUpdateSociety).Return(nil)

		_, err := s.UpdateSociety(ctx, expectedUpdateSociety)

		assert.NoError(t, err)
	})

	t.Run("should_return_error_if_uuid_not_found_in_context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		expectedUpdateSociety := &society.UpdateSocietyIn{
			SocietyUUID:    uuid.Generate().String(),
			Name:           "Test1",
			Description:    "A test society",
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
		mockDBRepo.EXPECT().IsOwnerAdminModerator(ctx, ownerUUID, societyUUID).Return(1, nil)
		mockDBRepo.EXPECT().UpdateSociety(ctx, expectedUpdateSociety).Return(expectedError)
		mockLogger.EXPECT().Error("failed to UpdateSociety from BD")

		_, err := s.UpdateSociety(ctx, expectedUpdateSociety)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestServer_RemoveSociety(t *testing.T) {
	t.Parallel()

	db, driverMock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	s := &Server{dbR: mockDBRepo}

	societyUUID := "soc-123"
	userUUID := "user-abc"
	in := &society.RemoveSocietyIn{SocietyUUID: societyUUID}

	ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	driverMock.ExpectBegin()  // sqlxDB.Beginx()
	driverMock.ExpectCommit() // tx.Commit()

	mockLogger.EXPECT().AddFuncName("RemoveSociety")

	mockDBRepo.
		EXPECT().
		GetOwner(ctx, societyUUID).
		Return(userUUID, nil)

	mockDBRepo.
		EXPECT().
		Conn().
		Return(sqlxDB)

	mockDBRepo.
		EXPECT().
		RemoveSocietyHasTagsEntry(ctx, societyUUID, gomock.Any()).
		Return(nil)

	mockDBRepo.
		EXPECT().
		RemoveSociety(ctx, societyUUID, gomock.Any()).
		Return(nil)

	mockDBRepo.
		EXPECT().
		RemoveMembersRequestEntry(ctx, societyUUID, gomock.Any()).
		Return(nil)

	mockDBRepo.
		EXPECT().
		RemoveSocietyMembersEntry(ctx, societyUUID, gomock.Any()).
		Return(nil)

	out, err := s.RemoveSociety(ctx, in)

	assert.NoError(t, err)
	assert.Equal(t, &society.EmptySociety{}, out)

	require.NoError(t, driverMock.ExpectationsWereMet())
}

func TestServer_SubscribeToSociety(t *testing.T) {
	t.Parallel()

	db, driverMock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	s := &Server{dbR: mockDBRepo}

	societyUUID := "soc-abc"
	userUUID := "user-xyz"
	in := &society.SubscribeToSocietyIn{SocietyUUID: societyUUID}

	ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	// Успешный кейс: формат == 1 → AddSocietyMembers
	t.Run("success: format == 1", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SubscribeToSociety")
		mockDBRepo.EXPECT().GetFormatSociety(ctx, societyUUID).Return(1, nil)
		mockDBRepo.EXPECT().AddSocietyMembers(ctx, userUUID, societyUUID).Return(nil)

		out, err := s.SubscribeToSociety(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, &society.EmptySociety{}, out)
	})

	// Успешный кейс: формат != 1 → AddMembersRequests
	t.Run("success: format != 1", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SubscribeToSociety")
		mockDBRepo.EXPECT().GetFormatSociety(ctx, societyUUID).Return(2, nil)
		mockDBRepo.EXPECT().AddMembersRequests(ctx, userUUID, societyUUID).Return(nil)

		out, err := s.SubscribeToSociety(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, &society.EmptySociety{}, out)
	})

	// Ошибка: uuid отсутствует в context
	t.Run("fail: missing uuid", func(t *testing.T) {
		badCtx := context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		mockLogger.EXPECT().AddFuncName("SubscribeToSociety")
		mockLogger.EXPECT().Error("failed to not found UUID in context")

		out, err := s.SubscribeToSociety(badCtx, in)
		assert.Nil(t, out)
		assert.ErrorContains(t, err, "uuid not found in context")
	})

	// Ошибка в GetFormatSociety
	t.Run("fail: GetFormatSociety error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SubscribeToSociety")
		mockDBRepo.EXPECT().GetFormatSociety(ctx, societyUUID).Return(0, errors.New("format error"))
		mockLogger.EXPECT().Error("failed to GetFormatSociety from BD")

		out, err := s.SubscribeToSociety(ctx, in)
		assert.Nil(t, out)
		assert.ErrorContains(t, err, "format error")
	})

	// Ошибка в AddMembersRequests
	t.Run("fail: AddMembersRequests error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SubscribeToSociety")
		mockDBRepo.EXPECT().GetFormatSociety(ctx, societyUUID).Return(2, nil)
		mockDBRepo.EXPECT().AddMembersRequests(ctx, userUUID, societyUUID).Return(errors.New("add req error"))
		mockLogger.EXPECT().Error("failed to AddMembersRequests from BD")

		out, err := s.SubscribeToSociety(ctx, in)
		assert.Nil(t, out)
		assert.ErrorContains(t, err, "add req error")
	})

	// Ошибка в AddSocietyMembers
	t.Run("fail: AddSocietyMembers error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SubscribeToSociety")
		mockDBRepo.EXPECT().GetFormatSociety(ctx, societyUUID).Return(1, nil)
		mockDBRepo.EXPECT().AddSocietyMembers(ctx, userUUID, societyUUID).Return(errors.New("add mem error"))
		mockLogger.EXPECT().Error("failed to AddSocietyMembers from BD")

		out, err := s.SubscribeToSociety(ctx, in)
		assert.Nil(t, out)
		assert.ErrorContains(t, err, "add mem error")
	})
	_, _ = sqlxDB, driverMock
}

func TestServer_UnSubscribeToSociety(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDbRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	s := &Server{dbR: mockDBRepo}

	userUUID := "user-abc"
	societyUUID := "soc-123"
	ctx := context.WithValue(context.Background(), config.KeyUUID, userUUID)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	in := &society.UnSubscribeToSocietyIn{SocietyUUID: societyUUID}

	t.Run("success", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("UnSubscribeToSociety")

		mockDBRepo.
			EXPECT().
			GetRoleSocietyMembers(gomock.Any(), userUUID, societyUUID).
			Return(1, nil)

		mockDBRepo.
			EXPECT().
			UnSubscribeToSociety(gomock.Any(), userUUID, societyUUID).
			Return(nil)

		out, err := s.UnSubscribeToSociety(ctx, in)

		assert.NoError(t, err)
		assert.Equal(t, &society.EmptySociety{}, out)
	})

	t.Run("uuid not in context", func(t *testing.T) {
		ctxWithoutUUID := context.WithValue(context.Background(), config.KeyLogger, mockLogger)

		mockLogger.EXPECT().AddFuncName("UnSubscribeToSociety")
		mockLogger.EXPECT().Error("failed to not found UUID in context")

		out, err := s.UnSubscribeToSociety(ctxWithoutUUID, in)

		assert.Nil(t, out)
		assert.ErrorContains(t, err, "uuid not found in context")
	})

	t.Run("GetRoleSocietyMembers returns error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("UnSubscribeToSociety")

		mockDBRepo.
			EXPECT().
			GetRoleSocietyMembers(gomock.Any(), userUUID, societyUUID).
			Return(0, errors.New("role fetch error"))

		mockLogger.EXPECT().Error("failed to GetRoleSocietyMembers from BD")

		out, err := s.UnSubscribeToSociety(ctx, in)

		assert.Nil(t, out)
		assert.ErrorContains(t, err, "role fetch error")
	})

	t.Run("UnSubscribeToSociety returns error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("UnSubscribeToSociety")

		mockDBRepo.
			EXPECT().
			GetRoleSocietyMembers(gomock.Any(), userUUID, societyUUID).
			Return(1, nil)

		mockDBRepo.
			EXPECT().
			UnSubscribeToSociety(gomock.Any(), userUUID, societyUUID).
			Return(errors.New("unsubscribe error"))

		mockLogger.EXPECT().Error("failed to UnSubscribeToSociety from BD")

		out, err := s.UnSubscribeToSociety(ctx, in)

		assert.Nil(t, out)
		assert.ErrorContains(t, err, "unsubscribe error")
	})
}
