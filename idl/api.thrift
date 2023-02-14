namespace go api

enum ErrCode {
    SuccessCode                = 0
    ServiceErrCode             = 10001
    ParamErrCode               = 10002
    UserAlreadyExistErrCode    = 10003
    AuthorizationFailedErrCode = 10004
}

struct douyin_comment_action_request {
  1: required string token       // 用户鉴权token
  2: required i64 video_id (vt.gt = "0", api.vd="$>0")      // 视频id
  3: required i32 action_type (vt.in = "1", vt.in = "2")   // 1-发布评论，2-删除评论
  4: optional string comment_text // 用户填写的评论内容，在action_type=1的时候使用
  5: optional i64 comment_id   // 要删除的评论id，在action_type=2的时候使用
}

struct douyin_comment_action_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: optional Comment comment // 评论成功返回评论内容，不需要重新拉取整个列表
}

struct douyin_comment_list_request {
  1: required string token // 用户鉴权token
  2: required i64 video_id (vt.gt = "0", api.vd="$>0") // 视频id
}

struct douyin_comment_list_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: list<Comment> comment_list // 评论列表
}

struct Comment {
  1: required i64 id (vt.gt = "0", api.vd="$>0") // 视频评论id
  2: required User user // 评论用户信息
  3: required string content // 评论内容
  4: required string create_date // 评论发布日期，格式 mm-dd
}

struct User {
  1: required i64 id (vt.gt = "0", api.vd="$>0") // 用户id
  2: required string name  // 用户名称
  3: optional i64 follow_count (vt.gt = "0")  // 关注总数
  4: optional i64 follower_count (vt.gt = "0")  // 粉丝总数
  5: required bool is_follow  // true-已关注，false-未关注
  6: required string avatar  // 用户头像Url
}

struct douyin_favorite_action_request {
  1: required string token  // 用户鉴权token
  2: required i64 video_id (vt.gt = "0", api.vd="$>0")  // 视频id
  3: required i32 action_type (vt.in = "1", vt.in = "2", api.vd = "$==1||$==2") // 1-点赞，2-取消点赞
}

struct douyin_favorite_action_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
}

struct douyin_favorite_list_request {
  1: required i64 user_id (vt.gt = "0", api.vd="$>0") // 用户id
  2: required string token // 用户鉴权token
}

struct douyin_favorite_list_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: list<Video> video_list // 用户点赞视频列表
}

struct Video {
  1: required i64 id (vt.gt = "0") // 视频唯一标识
  2: required User author // 视频作者信息
  3: required string play_url // 视频播放地址
  4: required string cover_url // 视频封面地址
  5: required i64 favorite_count // 视频的点赞总数
  6: required i64 comment_count // 视频的评论总数
  7: required bool is_favorite // true-已点赞，false-未点赞
  8: required string title (vt.max_size = "60") // 视频标题
}

struct douyin_feed_request {
  1: optional i64 latest_time // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
  2: optional string token // 可选参数，登录用户设置
}

struct douyin_feed_response {
  1: required i64 status_code   // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: list<Video> video_list     // 视频列表
  4: optional i64 next_time     // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

struct douyin_message_chat_request {
  1: required string token // 用户鉴权token
  2: required i64 to_user_id (vt.gt = "0", api.vd="$>0")  // 对方用户id
}

struct douyin_message_chat_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: list<Message> message_list // 消息列表
}

struct Message {
  1: required i64 id (vt.gt = "0") // 消息id
  2: required i64 to_user_id (vt.gt = "0") // 该消息接收者的id
  3: required i64 from_user_id (vt.gt = "0") // 该消息发送者的id
  4: required string content // 消息内容
  5: optional i64 create_time // 消息创建时间
}

struct douyin_message_action_request {
  1: required string token // 用户鉴权token
  2: required i64 to_user_id (vt.gt = "0", api.vd="$>0") // 对方用户id
  3: required i32 action_type (vt.in = "1", api.vd="$==1") // 1-发送消息
  4: required string content // 消息内容
}

struct douyin_message_action_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
}

struct douyin_publish_action_request {
  1: required string token // 用户鉴权token
//  2: optional binary data // 视频数据
  2: required string title // 视频标题
}

struct douyin_publish_action_response {
  1: required i64 status_code  // 状态码，0-成功，其他值-失败
  2: optional string status_msg  // 返回状态描述
}

struct douyin_publish_list_request {
  1: required i64 user_id (vt.gt = "0", api.vd="$>0")  // 用户id
  2: required string token  // 用户鉴权token
}

struct douyin_publish_list_response {
  1: required i64 status_code  // 状态码，0-成功，其他值-失败
  2: optional string status_msg  // 返回状态描述
  3: list<Video> video_list  // 用户发布的视频列表
}


struct douyin_relation_action_request {
  1: required string token // 用户鉴权token
  2: required i64 to_user_id (vt.gt = "0", api.vd="$>0") // 对方用户id
  3: required i32 action_type (vt.in = "1", vt.in = "2") // 1-关注，2-取消关注
}

struct douyin_relation_action_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
}

struct douyin_relation_follow_list_request {
  1: required i64 user_id (vt.gt = "0", api.vd="$>0") // 用户id
  2: required string token // 用户鉴权token
}

struct douyin_relation_follow_list_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: list<User> user_list // 用户信息列表
}

struct douyin_relation_follower_list_request {
  1: required i64 user_id (vt.gt = "0", api.vd="$>0")  // 用户id
  2: required string token // 用户鉴权token
}

struct douyin_relation_follower_list_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg  // 返回状态描述
  3: list<User> user_list  // 用户列表
}


struct douyin_relation_friend_list_request {
  1: required i64 user_id (vt.gt = "0", api.vd="$>0")  // 用户id
  2: required string token  // 用户鉴权token
}

struct douyin_relation_friend_list_response {
  1: required i64 status_code  // 状态码，0-成功，其他值-失败
  2: optional string status_msg  // 返回状态描述
  3: list<FriendUser> user_list // 用户列表
}


struct FriendUser {
  1: required i64 id (vt.gt = "0") // 用户id
  2: required string name  // 用户名称
  3: optional i64 follow_count (vt.gt = "0")  // 关注总数
  4: optional i64 follower_count (vt.gt = "0")  // 粉丝总数
  5: required bool is_follow  // true-已关注，false-未关注
  6: required string avatar  // 用户头像Url
  7: optional string message // 和该好友的最新聊天消息
  8: required i64 msgType (vt.in = "0", vt.in = "1") // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
}


struct douyin_user_register_request {
  1: required string username (vt.min_size = "2", vt.max_size = "32", api.vd = "len($)>2 && len($)<32") // 注册用户名，最长32个字符
  2: required string password (vt.min_size = "6", vt.max_size = "32", vt.pattern = "[0-9A-Za-z]+", api.vd = "len($)>2 && len($)<32") // 密码，最长32个字符
}

struct douyin_user_register_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: required i64 user_id (vt.gt = "0") // 用户id
  4: required string token // 用户鉴权token
}

struct douyin_user_login_request {
  1: required string username (vt.min_size = "2", vt.max_size = "32", api.vd = "len($)>2 && len($)<32") // 登录用户名
  2: required string password (vt.min_size = "6", vt.max_size = "32", api.vd = "len($)>2 && len($)<32") // 登录密码
}

struct douyin_user_login_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: required i64 user_id (vt.gt = "0") // 用户id
  4: required string token // 用户鉴权token
}

struct douyin_user_request {
  1: required i64 user_id (vt.gt = "0", api.vd = "$>0") // 用户id
  2: required string token // 用户鉴权token
}

struct douyin_user_response {
  1: required i64 status_code // 状态码，0-成功，其他值-失败
  2: optional string status_msg // 返回状态描述
  3: required User user // 用户信息
}


// 基础接口
service FeedService {
    douyin_feed_response GetFeed(1: douyin_feed_request req) (api.get="/douyin/feed/")
}

service UserService {
    douyin_user_register_response Register(1: douyin_user_register_request req) (api.post="/douyin/user/register/")
    douyin_user_login_response Login(1: douyin_user_login_request req) (api.post="/douyin/user/login/")
    douyin_user_response GetUserInfo(1: douyin_user_request req) (api.get="/douyin/user/")
}

service PublishService {
    douyin_publish_action_response PublishAction(1: douyin_publish_action_request req) (api.post="/douyin/publish/action/")
    douyin_publish_list_response GetPublishVideos(1: douyin_publish_list_request req) (api.get="/douyin/publish/list/")
}

// 互动接口
service FavoriteService {
    douyin_favorite_action_response FavoriteVideo(1: douyin_favorite_action_request req) (api.post="/douyin/favorite/action/")
    douyin_favorite_list_response GetFavoriteList(1: douyin_favorite_list_request req) (api.get="/douyin/favorite/list/")
}

service CommentService {
    douyin_comment_action_response CommentAction(1: douyin_comment_action_request req) (api.post="/douyin/comment/action/")
    douyin_comment_list_response GetCommentList(1: douyin_comment_list_request req) (api.get="/douyin/comment/list/")
}

// 社交接口
service RelationService {
    douyin_relation_action_response Follow(1: douyin_relation_action_request req) (api.post="/douyin/relation/action/")
    douyin_relation_follow_list_response GetFollowList(1: douyin_relation_follow_list_request req) (api.get="/douyin/relation/follow/list/")
    douyin_relation_follower_list_response GetFollowerList(1: douyin_relation_follower_list_request req) (api.get="/douyin/relation/follower/list/")
    douyin_relation_friend_list_response GetFriendList(1: douyin_relation_friend_list_request req) (api.get="/douyin/relation/friend/list/")
}

service MessageService {
    douyin_message_action_response SendMessage(1: douyin_message_action_request req) (api.post="/douyin/message/action/")
    douyin_message_chat_response GetMessageChat(1: douyin_message_chat_request req) (api.get="/douyin/message/chat/")
}