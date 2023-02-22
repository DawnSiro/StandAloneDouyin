package service

import (
	"reflect"
	"strconv"
	"strings"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/util/sensitive"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/go-redis/redis"
)

func PostComment(userID, videoID uint64, commentText string) (*api.DouyinCommentActionResponse, error) {
	// 删除redis评论列表缓存
	// 使用 strings.Builder 来优化字符串的拼接
	var builder strings.Builder
	builder.WriteString(strconv.FormatUint(videoID, 10))
	builder.WriteString("_video_comments")
	delCommentListKey := builder.String()

	// TODO 业务优化
	keysMatch, err := db.VideoCRDB.Do("keys", "*"+delCommentListKey).Result()
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
	}
	if reflect.TypeOf(keysMatch).Kind() == reflect.Slice {
		val := reflect.ValueOf(keysMatch)
		// 删除key
		for i := 0; i < val.Len(); i++ {
			db.VideoCRDB.Del(val.Index(i).Interface().(string))
			hlog.Info("删除了RedisKey:", val.Index(i).Interface().(string))
		}
	}

	//检测是否带有敏感词
	if sensitive.IsWordsFilter(commentText) {
		return nil, errno.ContainsProhibitedSensitiveWordsError
	}

	dbc, err := db.CreateComment(videoID, commentText, userID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		return nil, err
	}

	dbu, err := db.SelectUserByID(userID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		return nil, err
	}
	authorID, err := db.SelectAuthorIDByVideoID(videoID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		return nil, err
	}

	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment(dbc, dbu, db.IsFollow(userID, authorID)),
	}, nil
}

func DeleteComment(userID, videoID, commentID uint64) (*api.DouyinCommentActionResponse, error) {
	// 查询此评论是否是本人发送的
	isComment := db.IsCommentCreatedByMyself(userID, commentID)
	// 非本人评论
	if !isComment {
		hlog.Error("service.comment.DeleteComment err:", errno.DeletePermissionError)
		return nil, errno.DeletePermissionError
	}

	dbc, err := db.DeleteCommentByID(videoID, commentID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}
	dbu, err := db.SelectUserByID(userID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}
	authorID, err := db.SelectAuthorIDByVideoID(videoID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}

	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment(dbc, dbu, db.IsFollow(userID, authorID)),
	}, nil
}

func CommentList(userID, videoID uint64) (*api.DouyinCommentListResponse, error) {
	var builder strings.Builder
	builder.WriteString(strconv.FormatUint(userID, 10))
	builder.WriteString("_userId_")
	builder.WriteString(strconv.FormatUint(videoID, 10))
	builder.WriteString("_video_comments")
	commentListKey := builder.String()

	commentList, err := db.VideoCRDB.Get(commentListKey).Result()
	if err == redis.Nil {
		dbcList, err := db.SelectCommentListByVideoID(videoID)
		if err != nil {
			hlog.Error("service.comment.CommentList err:", err.Error())
			return nil, err
		}

		cList := make([]*api.Comment, 0, len(dbcList))

		for i := 0; i < len(dbcList); i++ {
			u, _ := db.SelectUserByID(dbcList[i].UserID)
			cList = append(cList, pack.Comment(dbcList[i], u, db.IsFollow(userID, dbcList[i].UserID)))
		}

		//序列化
		marshalList, _ := json.Marshal(cList)
		_, err = db.VideoCRDB.Set(commentListKey, marshalList, 0).Result()
		if err != nil {
			hlog.Error("service.comment.CommentList err:", err.Error())
			return nil, err
		}
		commentList, err = db.VideoCRDB.Get(commentListKey).Result()
		if err != nil {
			hlog.Error("service.comment.CommentList err:", err.Error())
			return nil, err
		}
	}
	//反序列化
	var list []*api.Comment
	err = json.Unmarshal([]byte(commentList), &list)
	if err != nil {
		hlog.Error("service.comment.CommentList err:", err.Error())
		return nil, err
	}

	return &api.DouyinCommentListResponse{
		StatusCode:  0,
		CommentList: list,
	}, nil
}
