package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
	ssov1 "github.com/tty-monkey/auth-server-schemas/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyAppID = 0
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(_ context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	validate := validator.New()

	err := validate.Var(req.GetEmail(), "required,email")
	if err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) {
			switch validationErr.Tag() {
			case "required":
				return nil, status.Error(codes.InvalidArgument, "email is required")
			case "email":
				return nil, status.Error(codes.InvalidArgument, "invalid email format")
			}
		}
	}

	err = validate.Var(req.GetPassword(), "required")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyAppID {
		return nil, status.Error(codes.InvalidArgument, "app id is required")
	}

	// TODO: implement login via auth service

	return &ssov1.LoginResponse{
		Token: req.GetEmail(),
	}, nil
}

func (s *serverAPI) Register(_ context.Context, _ *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("not implemented")
}

func (s *serverAPI) IsAdmin(_ context.Context, _ *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("not implemented")
}
