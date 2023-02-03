package rpc

import "github.com/easycode666/douyin/pkg/config"

// InitRPC init rpc client
func InitRPC() {
	UserConfig := config.InitConfig("user-config")
	initUserRpc(UserConfig)

}
