package auth

import ssov1 "github.com/empfaze/golang_grpc/protos/gen/go/sso"

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}
