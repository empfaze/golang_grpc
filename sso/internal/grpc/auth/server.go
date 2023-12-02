package auth

import (
	"context"

	ssov1 "github.com/empfaze/golang_grpc/protos/gen/go/sso"
	"github.com/empfaze/golang_grpc/sso/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	appID := req.GetAppId()

	err := utils.IsValidLoginDto(email, password, appID)
	if err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, email, password, int(appID))
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	err := utils.IsValidRegisterDto(email, password)
	if err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, email, password)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	userId := req.GetUserId()

	if err := utils.IsValidAdminDto(userId); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
