package server

import "context"

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (e *Server) Run(_ context.Context) error {
	return nil
}

func (e *Server) Shutdown(_ context.Context) error {
	return nil
}
