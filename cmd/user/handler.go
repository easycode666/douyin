package main

import (
	"context"
	user "github.com/easycode666/douyin/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	// TODO: Your code here...
	return
}

// Info implements the UserServiceImpl interface.
func (s *UserServiceImpl) Info(ctx context.Context, req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {
	// TODO: Your code here...
	return
}
