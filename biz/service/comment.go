package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/util/sensitive"
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/go-redis/redis"
	"reflect"
	"strconv"
	"strings"
)

func PostComment(userID, videoID uint64, commentText string) (*api.DouyinCommentActionResponse, error) {
	filterErr := "评论带有敏感词"

	// 删除redis评论列表缓存
	// 使用 strings.Builder 来优化字符串的拼接
	var builder strings.Builder
	builder.WriteString(strconv.FormatUint(videoID, 10))
	builder.WriteString("_video_comments")
	delCommentListKey := builder.String()

	keysMatch, err := db.RDB.Do("keys", "*"+delCommentListKey).Result()
	if err != nil {
		hlog.Info("查询批量删除的redisKey失败", err)
	}
	if reflect.TypeOf(keysMatch).Kind() == reflect.Slice {
		val := reflect.ValueOf(keysMatch)
		// 删除key
		for i := 0; i < val.Len(); i++ {
			db.RDB.Del(val.Index(i).Interface().(string))
			hlog.Info("删除了rediskey:", val.Index(i).Interface().(string))
		}
	}

	//publish comment
	//检测了评论是否为空
	if commentText == "" {
		return nil, errors.New("评论不能为空")
	}
	//检测是否带有敏感词
	if sensitive.IsWordsFilter(commentText) {
		return nil, errors.New(filterErr)
	}
	dbc, err := db.CreateComment(videoID, commentText, userID)
	if err != nil {
		return nil, err
	}

	dbu, err := db.SelectUserByID(userID)
	if err != nil {
		return nil, err
	}
	authorID, err := db.SelectAuthorIDByVideoID(videoID)
	if err != nil {
		return nil, err
	}

	_, err = db.IncreaseCommentCount(videoID)

	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		Comment:    pack.Comment(dbc, dbu, db.IsFollow(userID, authorID)),
	}, nil
}

func DeleteComment(userID, videoID, commentID uint64) (*api.DouyinCommentActionResponse, error) {
	//delete comment
	//查询此评论是否是本人发送的
	isComment := db.IsCommentCreatedByMyself(userID, commentID)
	//非本人评论
	if !isComment {
		return nil, errors.New("delete failed")
	}
	dbc, err := db.DeleteCommentByID(commentID)
	if err != nil {
		return nil, err
	}
	dbu, err := db.SelectUserByID(userID)
	if err != nil {
		return nil, err
	}
	authorID, err := db.SelectAuthorIDByVideoID(videoID)
	if err != nil {
		return nil, err
	}

	_, err = db.DecreaseCommentCount(videoID)

	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		Comment:    pack.Comment(dbc, dbu, db.IsFollow(userID, authorID)),
	}, nil
}

func CommentList(userID, videoID uint64) (*api.DouyinCommentListResponse, error) {
	commentListKey := strconv.FormatUint(userID, 10) + "_userId_" + strconv.FormatUint(videoID, 10) + "_video" + "_comments"
	commentList, err := db.RDB.Get(commentListKey).Result()
	if err == redis.Nil {
		//find like_count in mysql
		dbcList, err := db.SelectCommentListByVideoID(videoID)
		if err != nil {
			return nil, err
		}

		cList := make([]*api.Comment, 0, len(dbcList))

		for i := 0; i < len(dbcList); i++ {
			u, _ := db.SelectUserByID(dbcList[i].UserID)
			cList = append(cList, pack.Comment(dbcList[i], u, db.IsFollow(userID, dbcList[i].UserID)))
		}

		//序列化
		marshalList, _ := json.Marshal(cList)
		_, err = db.RDB.Set(commentListKey, marshalList, 0).Result()
		if err != nil {
			hlog.Info("redis_error: ", err)
			return nil, err
		}
		commentList, err = db.RDB.Get(commentListKey).Result()
		if err != nil {
			hlog.Info("redis_error: ", err)
			return nil, err
		}
	}
	//反序列化
	var list []*api.Comment
	err = json.Unmarshal([]byte(commentList), &list)
	if err != nil {
		return nil, err
	}

	return &api.DouyinCommentListResponse{
		StatusCode:  0,
		StatusMsg:   nil,
		CommentList: list,
	}, nil
}
