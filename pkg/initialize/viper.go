package initialize

import (
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper(path string) {
	// 设置配置文件类型和路径
	viper.SetConfigType("yml")
	viper.SetConfigFile(path)
	// 读取配置信息
	err := viper.ReadInConfig()
	if err != nil {
		hlog.Fatal("initialize.viper.Viper err:", err.Error())
	}
	// 将读取到的配置信息反序列化到 Config 中
	err = viper.Unmarshal(&global.Config)
	if err != nil {
		hlog.Fatal("initialize.viper.Viper err:", err.Error())
	}
	hlog.Info("initialize.viper.Viper Config:", global.Config)
	// 监视配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		hlog.Info("initialize.viper.Viper 配置文件被修改：", e.Name)
	})
}
