package util

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gofrs/uuid"
)

const (
	regionID        = "" // bucket所在位置，可查看oss对象储存控制台的概况获取
	accessKeyID     = ""
	accessKeySecret = ""
	roleArn         = ""
	roleSessionName = ""
	endpoint        = ""
	bucketName      = ""
)

func GetSTS() *sts.AssumeRoleResponse {
	//构建一个阿里云客户端, 用于发起请求。
	//设置调用者（RAM用户或RAM角色）的AccessKey ID和AccessKey Secret。
	//第一个参数就是bucket所在位置，可查看oss对象储存控制台的概况获取
	//第二个参数就是步骤一获取的AccessKey ID
	//第三个参数就是步骤一获取的AccessKey Secret
	client, err := sts.NewClientWithAccessKey(regionID,
		accessKeyID, accessKeySecret)

	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	//设置参数。关于参数含义和设置方法，请参见《API参考》。
	request.RoleArn = roleArn                 //步骤三获取的角色ARN
	request.RoleSessionName = roleSessionName //步骤三中的RAM角色名称

	//发起请求，并得到响应。
	response, err := client.AssumeRole(request)
	if err != nil {
		hlog.Error("util.oss.GetSTS err:", err)
	}
	hlog.Info("util.oss.GetSTS response:", response)
	return response
}

func UploadVideo(file *io.Reader, fileName string) (videoURL, coverURL string, err error) {
	// 从STS服务获取的安全令牌（SecurityToken）。
	response := GetSTS()
	securityToken := response.Credentials.SecurityToken //上面获取的临时授权的数据里的Credentials.SecurityToken
	// 从STS服务获取的临时访问密钥（AccessKey ID和AccessKey Secret）。
	// 从STS服务获取临时访问凭证后，可以通过临时访问密钥和安全令牌生成OSSClient。
	// 创建OSSClient实例。
	// 第一个参数就是bucket的Endpoint，可以在对象储存oss控制台的bucket的概览得到，例如http://oss-cn-beijing.aliyuncs.com
	// 第二个参数就是上面获取的临时授权的数据里的Credentials.AccessKeyId
	// 第三个参数就是上面获取的临时授权的数据里的Credentials.AccessKeySecret
	client, err := oss.New(endpoint,
		response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, oss.SecurityToken(securityToken))
	if err != nil {
		hlog.Error("util.oss.UploadVideo err:", err)
		return
	}
	// 填写Bucket名称，例如examplebucket。
	bucketName := bucketName
	// 填写Object的完整路径，完整路径中不能包含Bucket名称，例如exampledir/exampleobject.txt。
	objectName := "DouYin/video/" + fileName

	bucket, err := client.Bucket(bucketName)
	// ObjectID 即为 OSS 中的文件路径
	err = bucket.PutObject(objectName, *file)
	if err != nil {
		hlog.Error("util.oss.UploadVideo err:", err)
		return "", "", err
	}

	// 使用阿里云OSS生成视频封面，默认取第一秒的截图
	style := "video/snapshot,t_1000,f_jpg,w_800,h_600"
	// 指定过期时间为10分钟。
	//expiration := time.Now().Add(time.Minute * 10)
	// 指定原图名称。如果图片不在Bucket根目录，需携带文件完整路径，例如example/example.jpg。
	u1, err := uuid.NewV4()
	if err != nil {
		hlog.Error("util.oss.UploadVideo err:", err)
		return "", "", err
	}

	// 指定用于存放处理后图片的Bucket名称，该Bucket需与原图所在Bucket在同一地域。
	targetBucketName := bucketName
	// 指定处理后图片名称。如果图片不在Bucket根目录，需携带文件完整访问路径，例如exampledir/example.jpg。
	coverName := "DouYin/cover/" + u1.String() + ".jpg"
	// 将图片缩放为固定宽高100 px后转存到指定存储空间。
	style = "video/snapshot,t_1000,f_jpg,w_800,h_600"
	process := fmt.Sprintf("%s|sys/saveas,o_%v,b_%v", style, base64.URLEncoding.EncodeToString([]byte(coverName)),
		base64.URLEncoding.EncodeToString([]byte(targetBucketName)))
	result, err := bucket.ProcessObject(objectName, process)
	if err != nil {
		hlog.Error("util.oss.UploadVideo err:", err)
		return "", "", err
	}

	hlog.Info("util.oss.GetSTS result:", result)

	// 填写Object的完整路径，完整路径中不能包含Bucket名称，例如exampledir/exampleobject.txt。
	//coverObjectName := result.Object

	videoURL = "https://" + bucketName + "." + endpoint + "/" + objectName
	coverURL = "https://" + bucketName + "." + endpoint + "/" + result.Object

	hlog.Info("util.oss.GetSTS videoURL:", videoURL)
	hlog.Info("util.oss.GetSTS coverURL:", coverURL)
	return
}
