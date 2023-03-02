package mw

import (
	"bytes"
	"context"
	"encoding/hex"
	"path"
	"strconv"
	"strings"

	"douyin/biz/model/api"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func VerifyFile() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		file, err := c.FormFile("data")
		if err != nil {
			hlog.Info("mw.jwt.VerifyFile err:", err.Error())
			c.JSON(consts.StatusOK, &api.DouyinPublishActionResponse{
				StatusCode: errno.UserUploadFileError.ErrCode,
				StatusMsg:  &errno.UserUploadFileError.ErrMsg,
			})
			c.Abort()
			return
		}
		if file.Size >= constant.MaxFileSize {
			hlog.Info("mw.jwt.VerifyFile err:", errno.VideoUploadedTooLargeError.Error())
			c.JSON(consts.StatusOK, &api.DouyinPublishActionResponse{
				StatusCode: errno.VideoUploadedTooLargeError.ErrCode,
				StatusMsg:  &errno.VideoUploadedTooLargeError.ErrMsg,
			})
			c.Abort()
			return
		}

		fileSuffix := path.Ext(file.Filename)
		if _, ok := global.FileSuffixWhiteList[fileSuffix]; ok == false {
			// 文件后缀名不在白名单内
			hlog.Info("mw.jwt.VerifyFile err:", errno.FileTypeUploadedNotSupportError.Error())
			c.JSON(consts.StatusOK, &api.DouyinPublishActionResponse{
				StatusCode: errno.FileTypeUploadedNotSupportError.ErrCode,
				StatusMsg:  &errno.FileTypeUploadedNotSupportError.ErrMsg,
			})
			c.Abort()
			return
		}
		// 通过文件字节流判断文件真实类型
		f, err := file.Open()
		buffer := make([]byte, 30)
		_, err = f.Read(buffer)
		fileType := getFileType(buffer)
		if fileType == "" {
			// 文件真实类型不在白名单内
			hlog.Info("mw.jwt.VerifyFile err:", errno.FileTypeUploadedNotMatchError.Error())
			c.JSON(consts.StatusOK, &api.DouyinPublishActionResponse{
				StatusCode: errno.FileTypeUploadedNotMatchError.ErrCode,
				StatusMsg:  &errno.FileTypeUploadedNotMatchError.ErrMsg,
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// bytesToHexString 根据字节切片获取文件头，一般是16位的16进制表示的字符串
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// getFileType 根据文件头来判断文件是否为对应类型
func getFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc)
	hlog.Info("mw.file_verification.getFileType fileCode: ", fileCode)
	global.FileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)

		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}
