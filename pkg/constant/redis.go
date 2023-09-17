package constant

import "time"

// Redis Key 相关
// redis key prefix
const (
	UserInfoRedisPrefix          = "user_info:"
	UserInfoCountHashRedisPrefix = "user_info_count:"
	FollowCountRedisFiled        = "follow_count:"
	FanCountRedisFiled           = "fan_count:"
	TotalFavoritedRedisFiled     = "favorited_count:"
	WorkCountRedisFiled          = "work_count:"
	FavoriteCountRedisFiled      = "favorite_count:"

	FollowIDRedisSetPrefix        = "follow_id:"
	FollowRedisZSetPrefix         = "follow:"
	FanIDRedisSetPrefix           = "fan_id:"
	FanRedisZSetPrefix            = "fan:"
	VideoInfoRedisPrefix          = "video_info:"
	VideoAuthorIDRedisPrefix      = "video_author:"
	FavoriteVideoIDRedisSetPrefix = "favorite_video_id:"
	FavoriteVideoRedisZSetPrefix  = "favorite_video:"

	PublishRedisZSetPrefix       = "publish:"
	VideoInfoCntRedisPrefix      = "video_info_count:"
	CommentCountRedisPrefix      = "comment_count:"
	LikeCountRedisPrefix         = "like_count:"
	CommentRedisZSetPrefix       = "comment:"
	CommentInfoRedisPrefix       = "comment_info:"
	MessageRedisPrefix           = "message:"
	MessageLatestTimeRedisPrefix = "message_latest_time:"
	LoginFailCounterRedisPrefix  = "login_fail_counter:"
	LoginFailLockRedisPrefix     = "user_login_fail_lock:"
	UserIDNilKey                 = "user_id_nil"
	VideoIDNilKey                = "video_id_nil"
)

// Redis 缓存过期时间相关
// 如果要调用 Lua 脚本的话，要转成 int64 ，因为 time.Duration 在 Lua 中无法序列化（marshal）
const (
	CommentInfoExpiration       = time.Hour * 2
	UserInfoExpiration          = time.Hour * 2
	VideoAuthorIDExpiration     = time.Hour * 2
	RelationRedisZSetExpiration = time.Hour * 2
)

// Redis Lua 脚本相关
const (
	// SetUserInfoLuaScript 设置用户信息 Lua 脚本
	SetUserInfoLuaScript = ``

	// CommentLuaScript 评论 Lua 脚本
	CommentLuaScript = `
		local infoKey = KEYS[1]
		local jsonString = ARGV[1]
		local commentInfoExpiration = tonumber(ARGV[2])
		local zSetKey = KEYS[2]
		local score = tonumber(ARGV[3])
		local member = ARGV[4]
		
		redis.call('SET', infoKey, jsonString, 'EX', commentInfoExpiration)
		
		if redis.call('EXISTS', zSetKey) > 0 then
			redis.call('ZAdd', zSetKey, score, member)
		end
		return 0
	`
	// DeleteCommentLuaScript 删除评论 Lua 脚本
	DeleteCommentLuaScript = `
		local infoKey = KEYS[1]
		local zSetKey = KEYS[2]
		local member = ARGV[1]
		
		redis.call('DEL', infoKey)
		
		if redis.call('EXISTS', zSetKey) > 0 then
			redis.call('ZRem', zSetKey, member)
		end
		
		return 0
	`
	// FollowLuaScript 关注的 Lua 脚本
	// 使用 relation 表的主键做为排序依据
	// 1. 关注者的关注 set 集合增加 ID
	// 2. 关注者的 ZSet 关注用户信息集合增加
	// 3. 关注者的关注数+1
	// 4. 被关注者的粉丝数+1
	// 5. 被关注者的 ZSet 粉丝信息集合增加
	FollowLuaScript = `
		local followIDKey = KEYS[1]
		local followKey = KEYS[2]
		local followInfoCountKey = KEYS[3]
		local fanKey = KEYS[4]
		local fanInfoCountKey = KEYS[5]
		local followCountFiled = KEYS[6]
		local fanCountFiled = KEYS[7]
		local createdTime = tonumber(ARGV[1])
		local toUserID = tonumber(ARGV[2])
		local userJSON = ARGV[3]
		local toUserJSON = ARGV[4]
		
		if redis.call("EXISTS", followIDKey) > 0 then 
			redis.call("SAdd", followIDKey, toUserID)
		end
		
		if redis.call("EXISTS", followKey) > 0 then
			redis.call("ZAdd", followKey, createdTime, toUserJSON)
		end
		
		if redis.call("EXISTS", followInfoCountKey) > 0 then 
			redis.call("HIncrBy", followInfoCountKey, followCountFiled, 1)
		end
		
		if redis.call("EXISTS", fanKey) > 0 then 
			redis.call("ZAdd", fanKey, createdTime, userJSON)
		end
		
		if redis.call("EXISTS", fanInfoCountKey) > 0 then
			redis.call("HIncrBy", fanInfoCountKey, fanCountFiled, 1)
		end
		
		return 0
	`
	// CancelFollowLuaScript 取消关注的 Lua 脚本
	CancelFollowLuaScript = `
		local followIDKey = KEYS[1]
		local followKey = KEYS[2]
		local followInfoCountKey = KEYS[3]
		local fanKey = KEYS[4]
		local fanInfoCountKey = KEYS[5]
		local followCountFiled = KEYS[6]
		local fanCountFiled = KEYS[7]
		local createdTime = tonumber(ARGV[1])
		local toUserID = tonumber(ARGV[2])
		local userJSON = ARGV[3]
		local toUserJSON = ARGV[4]
		
		if redis.call("EXISTS", followIDKey) > 0 then 
			redis.call("SRem", followIDKey, toUserID)
		end
		
		if redis.call("EXISTS", followKey) > 0 then
			redis.call("ZRem", followKey, toUserJSON)
		end
		
		if redis.call("EXISTS", followInfoCountKey) > 0 then 
			redis.call("HIncrBy", followInfoCountKey, followCountFiled, -1)
		end
		
		if redis.call("EXISTS", fanKey) > 0 then 
			redis.call("ZRem", fanKey, userJSON)
		end
		
		if redis.call("EXISTS", fanInfoCountKey) > 0 then
			redis.call("HIncrBy", fanInfoCountKey, fanCountFiled, -1)
		end
		
		return 0
	`
	// UnLockLuaScript 释放锁的 Lua 脚本，判断 Key 中的 Value
	UnLockLuaScript = `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
)
