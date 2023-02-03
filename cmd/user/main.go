package main

import (
	"fmt"
	"github.com/easycode666/douyin/pkg/config"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/easycode666/douyin/dal"
	"github.com/easycode666/douyin/kitex_gen/user/userservice"
	"github.com/easycode666/douyin/pkg/consts"
	"github.com/easycode666/douyin/pkg/mw"
)

// 获取配置
var (
	viper          = config.InitConfig("user-config")
	EtcdAddress    = fmt.Sprintf("%s:%d", viper.GetString("Etcd.Address"), viper.GetInt("Etcd.Port"))
	ServiceAddr    = fmt.Sprintf("%s:%d", viper.GetString("Server.Address"), viper.GetInt("Server.Port"))
	ServiceName    = viper.GetString("Server.Name")
	ExportEndpoint = viper.GetString("Otel.ExportEndpoint")
)

func Init() {
	dal.Init()
	// klog init
	klog.SetLogger(kitexlogrus.NewLogger())
	klog.SetLevel(klog.LevelInfo)
}

func main() {
	// 服务注册
	r, err := etcd.NewEtcdRegistry([]string{EtcdAddress})
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr(consts.TCP, ServiceAddr)
	if err != nil {
		panic(err)
	}

	Init()

	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(ServiceName),
		provider.WithExportEndpoint(ExportEndpoint),
		provider.WithInsecure(),
	)
	svr := userservice.NewServer(new(UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithMiddleware(mw.CommonMiddleware),
		server.WithMiddleware(mw.ServerMiddleware),
		server.WithSuite(tracing.NewServerSuite()),
		// Please keep the same as provider.WithServiceName
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: ServiceName}),
	)

	if err := svr.Run(); err != nil {
		klog.Fatalf("%s stopped with error:", ServiceName, err)
	}
}
