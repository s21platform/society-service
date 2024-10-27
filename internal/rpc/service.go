package rpc

import society "github.com/s21platform/society-proto/society-proto"

type Server struct {
	society.UnimplementedSocietyServiceServer
}

func New() *Server {
	return &Server{}
}
