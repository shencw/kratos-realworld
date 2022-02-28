package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/shencw/kratos-realworld/api/realworld/v1"
	"github.com/shencw/kratos-realworld/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RealWorldService struct {
	v1.UnimplementedRealWorldServer

	auth *biz.AuthUseCase
	log  *log.Helper
}

func NewRealWorldService(auth *biz.AuthUseCase, logger log.Logger) *RealWorldService {

	log.NewHelper(logger).Info("NewRealWorldService")

	return &RealWorldService{auth: auth, log: log.NewHelper(logger)}
}

func (s *RealWorldService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.UserReply, error) {
	s.log.WithContext(ctx).Infof("Auth Login: %v", in.User.GetEmail())

	return &v1.UserReply{
		User: &v1.UserReply_User{
			Email:    in.User.GetEmail(),
			Token:    "aaaa",
			Username: "wei",
			Bio:      "1",
		},
	}, nil
}

func (s *RealWorldService) Register(ctx context.Context, in *v1.RegisterRequest) (*v1.UserReply, error) {
	s.log.WithContext(ctx).Infof("Auth Register: %v", in.User.GetEmail())

	_, err := s.auth.Register(ctx, in.GetUser().Username, in.GetUser().Email, in.GetUser().Password)

	if err != nil {
		s.log.WithContext(ctx).Errorf("Auth Register: %v", err)
		return &v1.UserReply{}, err
	}

	return &v1.UserReply{User: &v1.UserReply_User{
		Email:    in.GetUser().GetEmail(),
		Token:    "token",
		Username: in.GetUser().GetUsername(),
		Bio:      "bio",
	}}, nil
}

func GetCurrentUser(ctx context.Context, in *emptypb.Empty) (*v1.UserReply, error) {
	return &v1.UserReply{}, nil
}

func UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UserReply, error) {
	return &v1.UserReply{}, nil
}
