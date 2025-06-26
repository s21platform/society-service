package service

import (
	"context"
	"fmt"

	logger_lib "github.com/s21platform/logger-lib"

	society "github.com/s21platform/society-proto/society-proto"
	"github.com/s21platform/society-service/internal/config"
	"github.com/s21platform/society-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	society.UnimplementedSocietyServiceServer
	dbR DbRepo
}

func New(repo DbRepo) *Server {
	return &Server{
		dbR: repo,
	}
}

func (s *Server) CreateSociety(ctx context.Context, in *society.SetSocietyIn) (*society.SetSocietyOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("CreateSociety")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to not found UUID in context")
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	if in.Name == "" {
		logger.Error("failed to Name society is empty")
		return nil, status.Error(codes.InvalidArgument, "name not provided")
	}

	SocietyData := model.SocietyData{
		Name:           in.Name,
		FormatID:       in.FormatID,
		PostPermission: in.PostPermissionID,
		IsSearch:       in.IsSearch,
		OwnerUUID:      uuid,
	}
	societyUUID, err := s.dbR.CreateSociety(ctx, &SocietyData)
	if err != nil {
		logger.Error("failed to CreateSociety from BD")
		return nil, err
	}
	return &society.SetSocietyOut{SocietyUUID: societyUUID}, status.Error(codes.OK, "success")
}

func (s *Server) GetSocietyInfo(ctx context.Context, in *society.GetSocietyInfoIn) (*society.GetSocietyInfoOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetSocietyInfo")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to not found UUID in context")
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	if in.SocietyUUID == "" {
		logger.Error("failed to SocietyUUID is empty")
		return nil, status.Error(codes.InvalidArgument, "societyUUID not provided")
	}

	societyInfo, err := s.dbR.GetSocietyInfo(ctx, in.SocietyUUID)

	if err != nil {
		logger.Error("failed to GetSocietyInfo from BD")
		return nil, err
	}

	if !societyInfo.Description.Valid {
		societyInfo.Description.String = ""
	}

	count, err := s.dbR.CountSubscribe(ctx, in.SocietyUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get count of subscribers: %w", err)
	}
	societyInfo.CountSubscribe = count

	getTag, err := s.dbR.GetTags(ctx, in.SocietyUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	societyInfo.TagsID = getTag

	var tags []*society.TagsID
	for _, tag := range societyInfo.TagsID {
		tags = append(tags, &society.TagsID{TagID: tag})
	}

	description := ""
	if societyInfo.Description.Valid {
		description = societyInfo.Description.String
	}

	canEditSocietyInt, err := s.dbR.IsOwnerAdminModerator(ctx, uuid, in.SocietyUUID)
	canEditSociety := false
	if err != nil {
		logger.Error("failed to IsOwnerAdminModerator from BD")
		return nil, err
	}
	if canEditSocietyInt <= 3 && canEditSocietyInt >= 1 {
		canEditSociety = true
	}
	societyInfo.CanEditSociety = canEditSociety

	out := &society.GetSocietyInfoOut{
		Name:           societyInfo.Name,
		Description:    description,
		OwnerUUID:      societyInfo.OwnerUUID,
		PhotoURL:       societyInfo.PhotoURL,
		FormatID:       societyInfo.FormatID,
		PostPermission: societyInfo.PostPermission,
		IsSearch:       societyInfo.IsSearch,
		CountSubscribe: societyInfo.CountSubscribe,
		TagsID:         tags,
		CanEditSociety: canEditSociety,
	}
	return out, nil
}

func (s *Server) UpdateSociety(ctx context.Context, in *society.UpdateSocietyIn) (*society.EmptySociety, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("UpdateSociety")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to not found UUID in context")
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	if in.SocietyUUID == "" {
		logger.Error("failed to SocietyUUID is empty")
		return nil, status.Error(codes.InvalidArgument, "societyUUID not provided")
	}

	if in.Name == "" {
		logger.Error("failed to Name society is empty")
		return nil, status.Error(codes.InvalidArgument, "failed to name not provided")
	}

	isAllowed, err := s.dbR.IsOwnerAdminModerator(ctx, uuid, in.SocietyUUID)
	if err != nil {
		logger.Error("failed to IsOwnerAdminModerator from BD")
		return nil, status.Error(codes.InvalidArgument, "failed to IsOwnerAdminModerator from BD")
	}

	if isAllowed == 0 {
		logger.Error("failed to IsOwnerAdminModerator from BD")
		return nil, status.Error(codes.InvalidArgument, "failed to IsOwnerAdminModerator from BD")
	}

	if isAllowed != 1 && isAllowed != 2 && isAllowed != 3 {
		logger.Error("failed to IsOwnerAdminModerator from BD")
		return nil, status.Error(codes.InvalidArgument, "failed to peer is not Owner, Admin or Moderator")
	}

	err = s.dbR.UpdateSociety(ctx, in)

	if err != nil {
		logger.Error("failed to UpdateSociety from BD")
		return nil, err
	}

	return &society.EmptySociety{}, nil
}

func (s *Server) RemoveSociety(ctx context.Context, in *society.RemoveSocietyIn) (*society.EmptySociety, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("RemoveSociety")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to not found UUID in context")
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	role, err := s.dbR.GetOwner(ctx, in.SocietyUUID)
	if err != nil {
		logger.Error("failed to IsOwnerAdminModerator from BD")
		return nil, err
	}
	if role != uuid {
		logger.Error("failed to IsOwnerAdminModerator from BD")
		return nil, status.Error(codes.InvalidArgument, "failed to CheckRole from BD: the user does not have the rights to delete the community.")
	}

	tx, err := s.dbR.Conn().Beginx()
	if err != nil {
		logger.Error("failed to begin transaction")
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}

	err = s.dbR.RemoveSocietyHasTagsEntry(ctx, in.SocietyUUID, tx)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("failed to RemoveSocietyHasTagsEntry from BD")
		return nil, err
	}
	err = s.dbR.RemoveMembersRequestEntry(ctx, in.SocietyUUID, tx)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("failed to RemoveMembersRequestEntry from BD")
		return nil, err
	}
	err = s.dbR.RemoveSocietyMembersEntry(ctx, in.SocietyUUID, tx)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("failed to RemoveSocietyMembersEntry from BD")
		return nil, err
	}
	err = s.dbR.RemoveSociety(ctx, in.SocietyUUID, tx)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("failed to RemoveSociety from BD")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error("failed to commit transaction")
		return nil, err
	}

	return &society.EmptySociety{}, nil
}

func (s *Server) SubscribeToSociety(ctx context.Context, in *society.SubscribeToSocietyIn) (*society.EmptySociety, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SubscribeToSociety")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to not found UUID in context")
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	format, err := s.dbR.GetFormatSociety(ctx, in.SocietyUUID)
	if err != nil {
		logger.Error("failed to GetFormatSociety from BD")
		return nil, err
	}

	if format != 1 {
		err = s.dbR.AddMembersRequests(ctx, uuid, in.SocietyUUID)
		if err != nil {
			logger.Error("failed to AddMembersRequests from BD")
			return nil, err
		}
	} else {
		err = s.dbR.AddSocietyMembers(ctx, uuid, in.SocietyUUID)
		if err != nil {
			logger.Error("failed to AddSocietyMembers from BD")
			return nil, err
		}
	}

	return &society.EmptySociety{}, nil
}

func (s *Server) UnSubscribeToSociety(ctx context.Context, in *society.UnSubscribeToSocietyIn) (*society.EmptySociety, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("UnSubscribeToSociety")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to not found UUID in context")
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	_, err := s.dbR.GetRoleSocietyMembers(ctx, uuid, in.SocietyUUID)
	if err != nil {
		logger.Error("failed to GetRoleSocietyMembers from BD")
		return nil, err
	}
	err = s.dbR.UnSubscribeToSociety(ctx, uuid, in.SocietyUUID)
	if err != nil {
		logger.Error("failed to UnSubscribeToSociety from BD")
		return nil, err
	}

	return &society.EmptySociety{}, nil
}

//func (s *Server) GetSocietyWithOffset(ctx context.Context, in *society.GetSocietyWithOffsetIn) (*society.GetSocietyWithOffsetOut, error) {
//	uuid, ok := ctx.Value(config.KeyUUID).(string)
//	if !ok {
//		return nil, fmt.Errorf("uuid not found in context")
//	}
//
//	withOffsetData := model.WithOffsetData{
//		Limit:  in.Limit,
//		Offset: in.Offset,
//		Name:   in.Name,
//		Uuid:   uuid,
//	}
//	data, err := s.dbR.GetSocietyWithOffset(&withOffsetData)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get society with offset: %v", err)
//	}
//	count, err := s.dbR.GetCountSocietyWithOffset(&withOffsetData)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get count society with offset: %v", err)
//	}
//	out := society.GetSocietyWithOffsetOut{
//		Society: make([]*society.Society, len(*data)),
//		Total:   count,
//	}
//	for j, i := range *data {
//		level := &society.Society{
//			Name:       i.Name,
//			AvatarLink: i.AvatarLink,
//			SocietyId:  i.SocietyId,
//			IsMember:   i.IsMember,
//		}
//		out.Society[j] = level
//	}
//
//	return &out, err
//}
//
//func (s *Server) SubscribeToSociety(ctx context.Context, in *society.SubscribeToSocietyIn) (*society.SubscribeToSocietyOut, error) {
//	uuid, ok := ctx.Value(config.KeyUUID).(string)
//	if !ok {
//		return nil, fmt.Errorf("uuid not found in context")
//	}
//
//	data, err := s.dbR.SubscribeToSociety(in.SocietyId, uuid)
//	if err != nil {
//		return nil, fmt.Errorf("failed to subcribe to society %v", err)
//	}
//
//	out := society.SubscribeToSocietyOut{
//		Success: data,
//	}
//	return &out, err
//}
//
//func (s *Server) UnsubscribeFromSociety(ctx context.Context, in *society.UnsubscribeFromSocietyIn) (*society.UnsubscribeFromSocietyOut, error) {
//	uuid, ok := ctx.Value(config.KeyUUID).(string)
//	if !ok {
//		return nil, fmt.Errorf("uuid not found in context")
//	}
//
//	data, err := s.dbR.UnsubscribeFromSociety(in.SocietyId, uuid)
//	if err != nil {
//		return nil, fmt.Errorf("failed to unsubcribe to society %v", err)
//	}
//
//	out := society.UnsubscribeFromSocietyOut{
//		Success: data,
//	}
//	return &out, err
//}
//
//func (s *Server) GetSocietiesForUser(ctx context.Context, in *society.GetSocietiesForUserIn) (*society.GetSocietiesForUserOut, error) {
//	uuid, ok := ctx.Value(config.KeyUUID).(string)
//	if !ok {
//		return nil, fmt.Errorf("uuid not found in context")
//	}
//	data, err := s.dbR.GetSocietiesForUser(uuid, in.UserUuid)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get society for user: %v", err)
//	}
//
//	out := society.GetSocietiesForUserOut{
//		Society: make([]*society.Society, len(*data)),
//	}
//	for j, i := range *data {
//		level := &society.Society{
//			Name:       i.Name,
//			AvatarLink: i.AvatarLink,
//			SocietyId:  i.SocietyId,
//			IsMember:   i.IsMember,
//			IsPrivate:  i.IsPrivate,
//		}
//		out.Society[j] = level
//	}
//	return &out, err
//}
