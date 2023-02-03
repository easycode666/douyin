package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/easycode666/douyin/kitex_gen/user"
	"github.com/easycode666/douyin/kitex_gen/user/userservice"
	"github.com/easycode666/douyin/pkg/config"
	"github.com/easycode666/douyin/pkg/errno"
	"github.com/easycode666/douyin/pkg/mw"
)

var userClient userservice.Client

// User RPC 客户端初始化
func initUserRpc(viper config.Viper) {
	EtcdAddress := fmt.Sprintf("%s:%d", viper.GetString("Etcd.Address"), viper.GetInt("Etcd.Port"))
	// 服务发现
	r, err := etcd.NewEtcdResolver([]string{EtcdAddress})
	if err != nil {
		panic(err)
	}
	ServiceName := viper.GetString("Server.Name")

	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(ServiceName),
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)

	c, err := userservice.NewClient(
		ServiceName,
		client.WithMiddleware(mw.CommonMiddleware),
		client.WithInstanceMW(mw.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(30*time.Second),             // rpc timeout
		client.WithConnectTimeout(30000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(tracing.NewClientSuite()),        // tracer
		client.WithResolver(r),                            // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: ServiceName}),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

func Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	resp, err = userClient.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

func Login(ctx context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	resp, err = userClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

func Info(ctx context.Context, req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {
	resp, err = userClient.Info(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}
