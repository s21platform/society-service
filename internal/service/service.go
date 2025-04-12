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
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("CreateSociety")
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
	}
	return out, nil
}

func (s *Server) UpdateSociety(ctx context.Context, in *society.UpdateSocietyIn) (*society.EmptySociety, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	logger.AddFuncName("UpdateSociety")

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

func (s *Server) GetSocietyWithOffset(ctx context.Context, in *society.GetSocietyWithOffsetIn) (*society.GetSocietyWithOffsetOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetSocietyWithOffset")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	if in.Limit < 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid limit: limit < 0")
	}

	if in.Offset < 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid offset: offset < 0")
	}
	withOffsetData := model.WithOffsetData{
		Limit:  in.Limit,
		Offset: in.Offset,
		Name:   in.Name,
		Uuid:   uuid,
	}

	data, err := s.dbR.GetSocietyWithOffset(ctx, &withOffsetData)
	if err != nil {
		logger.Error("failed to GetSocietyWithOffset from BD")
	}
	if len(*data) < 1 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	socityUuid := make([]string, len(*data))
	for i, j := range *data {
		socityUuid[i] = j.SocietyUUID
	}
	memberOfSociety, err := s.dbR.GetMemberOfSocieties(ctx, socityUuid)
	if err != nil {
		logger.Error("failed to GetMemberOfSocieties from BD")
	}

	for _, j := range *data {
		j.IsMember = memberOfSociety[j.SocietyUUID]
	}
	out := society.GetSocietyWithOffsetOut{
		Societies: make([]*society.Society, len(*data)),
		Total:     int64(len(*data)),
	}
	for j, i := range *data {
		level := &society.Society{
			SocietyUUID: i.SocietyUUID,
			Name:        i.Name,
			PhotoURL:    i.PhotoURL,
			IsMember:    i.IsMember,
			FormatId:    i.FormatId,
		}
		out.Societies[j] = level
	}

	return &out, nil
}
