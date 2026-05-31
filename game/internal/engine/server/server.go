package server

import "context"

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (e *Server) Run(ctx context.Context) error {
	return nil
}

func (e *Server) Shutdown(ctx context.Context) error {
	return nil
}
