package errno

import (
	"errors"
	"fmt"

	"douyin/biz/model/api"
)

type ErrNo struct {
	ErrCode int64
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("ErrCode=%d, ErrMsg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

var (
	Success                                       = NewErrNo(int64(api.ErrCode_Success), "一切正常")
	ClientError                                   = NewErrNo(int64(api.ErrCode_Client), "用户端错误")
	UserRegistrationError                         = NewErrNo(int64(api.ErrCode_UserRegistration), "用户注册错误")
	UsernameVerificationFailedError               = NewErrNo(int64(api.ErrCode_UsernameVerificationFailed), "用户名校验失败")
	UsernameAlreadyExistsError                    = NewErrNo(int64(api.ErrCode_UsernameAlreadyExists), "用户名已存在")
	PasswordVerificationFailedError               = NewErrNo(int64(api.ErrCode_PasswordVerificationFailed), "密码校验失败")
	PasswordLengthNotEnoughError                  = NewErrNo(int64(api.ErrCode_PasswordLengthNotEnough), "密码长度不够 ")
	PasswordStrengthNotEnoughError                = NewErrNo(int64(api.ErrCode_PasswordStrengthNotEnough), "密码强度不够")
	UserLoginError                                = NewErrNo(int64(api.ErrCode_UserLogin), "用户登陆异常")
	UserAccountDoesNotExistError                  = NewErrNo(int64(api.ErrCode_UserAccountDoesNotExist), "用户账户不存在")
	UserPasswordError                             = NewErrNo(int64(api.ErrCode_UserPassword), "用户密码错误")
	PasswordNumberOfTimesExceedsError             = NewErrNo(int64(api.ErrCode_PasswordNumberOfTimesExceeds), "用户输入密码次数超过限制")
	UserIdentityVerificationFailedError           = NewErrNo(int64(api.ErrCode_UserIdentityVerificationFailed), "用户身份校验失败")
	UserLoginHasExpiredError                      = NewErrNo(int64(api.ErrCode_UserLoginHasExpired), "用户登陆已过期")
	AccessPermissionError                         = NewErrNo(int64(api.ErrCode_AccessPermission), "访问权限异常")
	DeletePermissionError                         = NewErrNo(int64(api.ErrCode_DeletePermission), "删除权限异常")
	UserRequestParameterError                     = NewErrNo(int64(api.ErrCode_UserRequestParameter), "用户请求参数错误")
	RepeatOperationError                          = NewErrNo(int64(api.ErrCode_RepeatOperationError), "用户重复操作")
	IllegalUserInputError                         = NewErrNo(int64(api.ErrCode_IllegalUserInput), "用户输入内容非法")
	ContainsProhibitedSensitiveWordsError         = NewErrNo(int64(api.ErrCode_ContainsProhibitedSensitiveWords), "包含违禁敏感词")
	UserUploadFileError                           = NewErrNo(int64(api.ErrCode_UserUploadFile), "用户上传文件异常")
	FileTypeUploadedNotMatchError                 = NewErrNo(int64(api.ErrCode_FileTypeUploadedNotMatch), "用户上传文件类型不匹配")
	FileTypeUploadedNotSupportError               = NewErrNo(int64(api.ErrCode_FileTypeUploadedNotSupport), "用户上传文件类型不支持")
	VideoUploadedTooLargeError                    = NewErrNo(int64(api.ErrCode_VideoUploadedTooLarge), "用户上传视频太大")
	ServiceError                                  = NewErrNo(int64(api.ErrCode_Service), "未知错误")
	SystemExecutionError                          = NewErrNo(int64(api.ErrCode_SystemExecution), "系统执行出错")
	SystemExecutionTimeoutError                   = NewErrNo(int64(api.ErrCode_SystemExecutionTimeout), "系统执行超时")
	SystemDisasterToleranceFunctionTriggeredError = NewErrNo(int64(api.ErrCode_SystemDisasterToleranceFunctionTriggered), "系统容灾功能被触发")
	SystemResourceError                           = NewErrNo(int64(api.ErrCode_SystemResource), "系统资源异常")
	CallingThirdPartyServiceError                 = NewErrNo(int64(api.ErrCode_CallingThirdPartyService), "调用第三方服务出错")
	MiddlewareServiceError                        = NewErrNo(int64(api.ErrCode_MiddlewareService), "中间件服务出错")
	RPCServiceError                               = NewErrNo(int64(api.ErrCode_RPCService), "RPC 服务出错")
	RPCServiceNotFindError                        = NewErrNo(int64(api.ErrCode_RPCServiceNotFind), "RPC 服务未找到")
	RPCServiceNotRegisteredError                  = NewErrNo(int64(api.ErrCode_RPCServiceNotRegistered), "RPC 服务未注册")
	InterfaceNotExistError                        = NewErrNo(int64(api.ErrCode_InterfaceNotExist), "接口不存在")
	CacheServiceError                             = NewErrNo(int64(api.ErrCode_CacheService), "缓存服务出错")
	KeyLengthExceedsLimitError                    = NewErrNo(int64(api.ErrCode_KeyLengthExceedsLimit), "key 长度超过限制")
	ValueLengthExceedsLimitError                  = NewErrNo(int64(api.ErrCode_ValueLengthExceedsLimit), "value 长度超过限制")
	StorageCapacityFullError                      = NewErrNo(int64(api.ErrCode_StorageCapacityFull), "存储容量已满")
	UnsupportedDataFormatError                    = NewErrNo(int64(api.ErrCode_UnsupportedDataFormat), "不支持的数据格式")
	DatabaseServiceError                          = NewErrNo(int64(api.ErrCode_DatabaseService), "数据库服务出错")
	TableDoesNotExistError                        = NewErrNo(int64(api.ErrCode_TableDoesNotExist), "表不存在")
	ColumnDoesNotExistError                       = NewErrNo(int64(api.ErrCode_ColumnDoesNotExist), "列不存在")
	DatabaseDeadlockError                         = NewErrNo(int64(api.ErrCode_DatabaseDeadlock), "数据库死锁")
	VideoLikeLimitError                           = NewErrNo(int64(api.ErrCode_VideoLikeLimitError), "用户频繁点赞")
)

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}
	// 不属于列写出的错误码的错误的错误信息不返回客户端
	s := ServiceError
	return s
}
