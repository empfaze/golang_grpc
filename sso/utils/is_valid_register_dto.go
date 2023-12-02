package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsValidRegisterDto(email string, password string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "Email is required")
	}

	if password == "" {
		return status.Error(codes.InvalidArgument, "Password is required")
	}

	return nil
}
