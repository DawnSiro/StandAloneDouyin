// Code generated by hertz generator.

package main

import (
	"douyin/biz/handler/api/ws"
	"douyin/biz/mw"
	"douyin/pkg/initialize"
)

func Init() {

	initialize.Viper()
	initialize.MySQL()
	initialize.Redis()
	initialize.Global()
	mw.InitJWT()

	// initialize.Hertz() 需要保持在最下方，因为调用完后 Hertz 就启动完毕了
	go ws.MannaClient.Run()
	initialize.Hertz()
}

func main() {
	Init()
}
