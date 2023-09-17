package rdb

import (
	"douyin/dal/model"
	"douyin/pkg/constant"
	"douyin/pkg/global"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
)

type CommentInfo struct {
	ID          uint64
	VideoID     uint64
	UserID      uint64
	Content     string
	CreatedTime float64
}

type CommentIDZSet struct {
	CommentID  uint64
	CreateTime int64
}

func AddComment(videoID uint64, comment CommentInfo) error {
	jsonString, err := json.Marshal(comment)
	if err != nil {
		hlog.Error("dal.rdb.comment.AddComment err:", err.Error())
		return err
	}
	infoKey := constant.CommentInfoRedisPrefix + strconv.FormatUint(comment.ID, 10)
	zSetKey := constant.CommentRedisZSetPrefix + strconv.FormatUint(videoID, 10)
	keys := []string{infoKey, zSetKey}
	args := []interface{}{jsonString, int64(constant.CommentInfoExpiration), comment.CreatedTime, comment.ID}
	// 执行 Redis 服务端缓存的 Lua 脚本
	err = global.UserRC.EvalSha(global.CommentLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error("dal.rdb.comment.AddComment err:", err.Error())
		return err
	}
	return nil
}

func DeleteComment(videoID, commentID uint64) error {
	infoKey := constant.CommentInfoRedisPrefix + strconv.FormatUint(commentID, 10)
	zSetKey := constant.CommentRedisZSetPrefix + strconv.FormatUint(videoID, 10)
	keys := []string{infoKey, zSetKey}
	args := []interface{}{commentID}
	// 执行 Redis 服务端缓存的 Lua 脚本
	err := global.UserRC.EvalSha(global.DeleteCommentLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error("dal.rdb.comment.DeleteComment err:", err.Error())
		return err
	}
	return nil
}

// IsCommentCreatedByMyself 获取评论用户ID
func IsCommentCreatedByMyself(userID, commentID uint64) (bool, error) {
	key := constant.CommentInfoRedisPrefix + strconv.FormatUint(commentID, 10)
	// 这里查询的是一个 string ，查不到就是查不到
	result, err := global.CommentRC.Get(key).Result()
	if err != nil {
		hlog.Error("dal.rdb.comment.IsCommentCreatedByMyself err:", err.Error())
		return false, err
	}
	var commentInfo CommentInfo
	err = json.Unmarshal([]byte(result), &commentInfo)
	if err != nil {
		hlog.Error("dal.rdb.comment.IsCommentCreatedByMyself err:", err.Error())
		return false, err
	}
	// 解析完 JSON 判断是否相等，不相等就不是自己发的
	return userID == commentInfo.UserID, nil
}

func SetCommentIDByVideoID(videoID uint64, cIDZ []*CommentIDZSet) error {
	key := constant.CommentRedisZSetPrefix + strconv.FormatUint(videoID, 10)
	global.CommentRC.ZAdd(key)
	return nil
}

// GetCommentIDByVideoID 通过视频ID获取评论ID
func GetCommentIDByVideoID(videoID uint64) ([]string, error) {
	key := constant.CommentRedisZSetPrefix + strconv.FormatUint(videoID, 10)
	result, err := global.CommentRC.ZRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetCommentInfo(commentID uint64) (*model.Comment, error) {
	return nil, nil
}

// GetCommentList 获取
func GetCommentList(commentIDList []string) ([]*model.Comment, error) {
	jsonList, err := global.CommentRC.MGet(commentIDList...).Result()
	if err != nil {
		return nil, err
	}

	dbcList := make([]*model.Comment, len(jsonList))
	for i := 0; i < len(jsonList); i++ {
		err = json.Unmarshal(jsonList[i].([]byte), &dbcList[i])
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// GetCommentListByVideoID 获取视频评论数据
func GetCommentListByVideoID(videoID uint64) ([]*model.Comment, error) {
	commentIDList, err := GetCommentIDByVideoID(videoID)
	if err != nil {
		return nil, err
	}
	dbcList, err := GetCommentList(commentIDList)
	if err != nil {
		return nil, err
	}
	return dbcList, nil
}
