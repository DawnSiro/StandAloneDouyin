package constant

const (
	CommentTableName            = "comment"
	MessageTableName            = "message"
	RelationTableName           = "relation"
	UserFavoriteVideosTableName = "user_favorite_video"
	VideoTableName              = "video"
	UserTableName               = "user"
	SecretKey                   = "CloudWeRun"
	IdentityKey                 = "id"
	ApiServiceName              = "api"
	CommentServiceName          = "comment"
	FavoriteServiceName         = "Favorite"
	FeedServiceName             = "feed"
	MessageServiceName          = "message"
	PublishServiceName          = "publish"
	RelationServiceName         = "relation"
	UserServiceName             = "user"
	MySQLDefaultDSN             = "douyin:BS5sp3K4yZTiEJ4S@tcp(119.29.27.252:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	//MySQLDefaultDSN = "douyin:!wwTF5VK)vPglY@-@tcp(172.17.0.1:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	TCP             = "tcp"
	UserServiceAddr = ":30110"
	NoteServiceAddr = ":30120"
	ExportEndpoint  = ":4317"
	ETCDAddress     = ":2379"
	DefaultLimit    = 10
	MaxVideoNum     = 30
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

// 数据库层面的枚举
// 0-未删除，1-已删除
const (
	DataNotDeleted = 0
	DataDeleted    = 1
)
