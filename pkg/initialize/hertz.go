package initialize

import (
	"douyin/biz/router"
	"douyin/pkg/global"
	"github.com/cloudwego/hertz/pkg/app/server"
	"strings"
)

func Hertz() {
	var builder strings.Builder
	builder.WriteString(global.Config.HertzConfig.Host)
	builder.WriteString(":")
	builder.WriteString(global.Config.HertzConfig.Port)
	hostWithPorts := builder.String()

	h := server.Default(
		server.WithHostPorts(hostWithPorts),
		server.WithExitWaitTime(0),
	)

	router.GeneratedRegister(h)

	h.Spin()
}
