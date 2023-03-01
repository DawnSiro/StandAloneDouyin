package initialize

import (
	"douyin/pkg/global"
)

func Global() {
	// 创建白名单类型，这个可以在 winhex64 里查看，同类型的文件貌似有的还不太一样
	global.FileTypeMap.Store("0000002066747970", ".mp4")
	global.FileTypeMap.Store("0000001c66747970", ".mp4")
	global.FileTypeMap.Store("0000001866747970", ".mp4")
	global.FileTypeMap.Store("0000001466747970", ".mp4")
}
