package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsValidLoginDto(email string, password string, appID int32) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "Email is required")
	}

	if password == "" {
		return status.Error(codes.InvalidArgument, "Password is required")
	}

	if appID == 0 {
		return status.Error(codes.InvalidArgument, "App ID is required")
	}

	return nil
}
