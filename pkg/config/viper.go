package config

import (
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Viper = *viper.Viper

func InitConfig(cfgName string) Viper {
	v := viper.New()
	v.SetDefault("Server.Name", "Server")
	v.SetConfigName(cfgName)
	v.SetConfigType("yml")
	v.AddConfigPath("../../config")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		klog.Fatalf("error reading config: %s", err)
	}
	klog.Infof("Using configuration file '%s'", v.ConfigFileUsed())

	// 监听配置文件变更
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return v
}
