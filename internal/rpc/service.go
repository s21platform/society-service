package rpc

import (
	"context"

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
	if !ok {
		return nil, status.Error(codes.Internal, "uuid not found in context")
	}

	SocietyData := model.SocietyData{
		Name:           in.Name,
		FormatID:       in.FormatID,
		PostPermission: in.PostPermissionID,
		IsSearch:       in.IsSearch,
		OwnerUUID:      uuid,
	}
	societyUUID, err := s.dbR.CreateSociety(&SocietyData)
	if err != nil {
		return nil, err
	}
	return &society.SetSocietyOut{SocietyUUID: societyUUID}, status.Error(codes.OK, "success")
}

//func (s *Server) GetAccessLevel(context.Context, *society.EmptySociety) (*society.GetAccessLevelOut, error) {
//	data, err := s.dbR.GetAccessLevel()
//	if err != nil {
//		return nil, fmt.Errorf("s.dbR.GetAccessLevel %v", err)
//	}
//
//	out := society.GetAccessLevelOut{
//		Levels: make([]*society.AccessLevel, len(*data)),
//	}
//	for j, i := range *data {
//		level := &society.AccessLevel{
//			Id:          i.Id,
//			AccessLevel: i.AccessLevel,
//		}
//		out.Levels[j] = level
//	}
//
//	return &out, err
//}
//
//func (s *Server) GetPermissions(context.Context, *society.EmptySociety) (*society.GetPermissionsOut, error) {
//	data, err := s.dbR.GetPermissions()
//	if err != nil {
//		return nil, fmt.Errorf("failed to get permission: %v", err)
//	}
//
//	out := society.GetPermissionsOut{
//		Permissions: make([]*society.Permission, len(*data)),
//	}
//
//	for a, i := range *data {
//		level := &society.Permission{
//			Id:          i.Id,
//			Name:        i.Name,
//			Description: i.Description,
//		}
//		out.Permissions[a] = level
//	}
//
//	return &out, err
//}
//
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
//func (s *Server) GetSocietyInfo(ctx context.Context, in *society.GetSocietyInfoIn) (*society.GetSocietyInfoOut, error) {
//	data, err := s.dbR.GetSocietyInfo(in.Id)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get society info: %v", err)
//	}
//
//	out := society.GetSocietyInfoOut{
//		Name:             data.Name,
//		Description:      data.Description,
//		OwnerUUID:        data.OwnerId,
//		PhotoUrl:         data.PhotoUrl,
//		IsPrivate:        data.IsPrivate,
//		CountSubscribers: data.CountSubscribers,
//	}
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
