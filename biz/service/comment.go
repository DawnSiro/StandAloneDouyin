package service

import (
	"sort"
	"sync"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/util/sensitive"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func PostComment(userID, videoID uint64, commentText string) (*api.DouyinCommentActionResponse, error) {
	// 删除redis评论列表缓存
	// 使用 strings.Builder 来优化字符串的拼接
	//var builder strings.Builder
	//builder.WriteString(strconv.FormatUint(videoID, 10))
	//builder.WriteString("_video_comments")
	//delCommentListKey := builder.String()
	//hlog.Info("service.comment.PostComment delCommentListKey:", delCommentListKey)

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

func GetCommentList(userID, videoID uint64) (*api.DouyinCommentListResponse, error) {
	//var builder strings.Builder
	//builder.WriteString(strconv.FormatUint(userID, 10))
	//builder.WriteString("_userId_")
	//builder.WriteString(strconv.FormatUint(videoID, 10))
	//builder.WriteString("_video_comments")
	//commentListKey := builder.String()

	//commentList, err := global.VideoCRC.Get(commentListKey).Result()
	//if err == redis.Nil {
	dbcList, err := db.SelectCommentListByVideoID(videoID)
	if err != nil {
		hlog.Error("service.comment.GetCommentList err:", err.Error())
		return nil, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(dbcList))
	// TODO 优化循环查询
	cList := make([]*api.Comment, 0, len(dbcList))
	for i := 0; i < len(dbcList); i++ {
		go func(c *db.Comment) {
			u, _ := db.SelectUserByID(c.UserID)
			cList = append(cList, pack.Comment(c, u, db.IsFollow(userID, c.UserID)))
			wg.Done()
		}(dbcList[i])
	}

	wg.Wait()
	sort.Sort(CommentSlice(cList))

	return &api.DouyinCommentListResponse{
		StatusCode:  0,
		CommentList: cList,
	}, nil
}

// CommentSlice 排序用的变量类型，用于实现三个排序需要的方法
type CommentSlice []*api.Comment

func (a CommentSlice) Len() int { //重写Len()方法
	return len(a)
}
func (a CommentSlice) Swap(i, j int) { //重写Swap()方法
	a[i], a[j] = a[j], a[i]
}
func (a CommentSlice) Less(i, j int) bool { //重写Less()方法
	return a[i].ID > a[j].ID
}
