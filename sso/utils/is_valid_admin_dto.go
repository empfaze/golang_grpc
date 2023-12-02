package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsValidAdminDto(userId int64) error {
	if userId == 0 {
		return status.Error(codes.InvalidArgument, "User ID is required")
	}

	return nil
}
