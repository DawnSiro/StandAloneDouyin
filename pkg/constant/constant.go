package constant

import "time"

// 中间件相关
const (
	TokenTimeOut    = 12 * time.Hour
	TokenMaxRefresh = 3 * time.Hour
)

// 点赞、登录限制相关
const (
	VideoLikeLimit     = 10
	VideoLikeLimitTime = 60 * time.Second
	UserLoginLimit     = 10
	UserLoginLimitTime = 10 * time.Minute
)

// ActionType 的枚举
// 1-发布评论，2-删除评论
// 1-点赞，2-取消点赞
// 1-关注，2-取消关注
// 1-发送消息
// 0-当前请求用户接收的消息， 1-当前请求用户发送的消息
const (
	PostComment       = 1
	DeleteComment     = 2
	Favorite          = 1
	CancelFavorite    = 2
	Follow            = 1
	CancelFollow      = 2
	SendMessageAction = 1
	ReceivedMessage   = 0
	SentMessage       = 1
)

// 数据库层面
// 0-未删除，1-已删除
const (
	DataNotDeleted              = 0
	DataDeleted                 = 1
	CommentTableName            = "`comment`"
	MessageTableName            = "`message`"
	RelationTableName           = "`relation`"
	UserFavoriteVideosTableName = "`user_favorite_video`"
	VideoTableName              = "`video`"
	UserTableName               = "`user`"
	FanListLimit                = 500 // 粉丝列表默认只展示前500个
)

// 业务相关
const (
	MaxVideoNum = 30
	MaxFileSize = 3 * 1024 * 1024 // 3MB 另外 Hertz 默认的请求体大小是 4MB
)
