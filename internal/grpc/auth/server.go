package auth

import (
	"context"

	ssov1 "github.com/tty-monkey/auth-server-schemas/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(_ context.Context, _ *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	panic("not implemented")
}

func (s *serverAPI) Register(_ context.Context, _ *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("not implemented")
}

func (s *serverAPI) IsAdmin(_ context.Context, _ *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("not implemented")
}
