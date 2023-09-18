package service

import (
	"context"
	"strconv"
	"time"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/model"
	"douyin/dal/pack"
	"douyin/dal/rdb"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"douyin/pkg/pulsar"
	"douyin/pkg/util"
	"douyin/pkg/util/sensitive"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func PostComment(userID, videoID uint64, commentText string) (*api.DouyinCommentActionResponse, error) {
	//检测是否带有敏感词
	if sensitive.IsWordsFilter(commentText) {
		return nil, errno.ContainsProhibitedSensitiveWordsError
	}

	// 基于雪花算法生成comment_id
	id, err := util.GetSonyFlakeID()
	if err != nil {
		hlog.Error("service.comment.PostComment err: failed to create comment id, ", err.Error())
	}

	// 发布消息队列
	msg := pulsar.PostCommentMessage{
		ID:          id,
		VideoID:     videoID,
		UserID:      userID,
		Content:     commentText,
		CreatedTime: time.Now(),
	}
	err = pulsar.GetPostCommentMQInstance().PostComment(msg)
	if err != nil {
		hlog.Error("service.comment.PostComment err: failed to publish mq ", err.Error())
		return nil, err
	}

	// 查询缓存数据
	dbu, err := rdb.GetUserInfo(userID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		// 缓存没有再查数据库
		dbu, err = db.SelectUserByID(userID)
		if err != nil {
			hlog.Error("service.comment.PostComment err:", err.Error())
			return nil, err
		}
		// 然后设置缓存
		err = rdb.SetUserInfo(dbu)
		if err != nil {
			// 要是设置出错，也不返回，继续执行逻辑
			hlog.Error("service.comment.PostComment err:", err.Error())
		}
	}

	// 这里的 isFollow 直接返回 false ，因为评论人自己当然不能关注自己
	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment((*model.Comment)(&msg), dbu, false),
	}, nil
}

func DeleteComment(userID, videoID, commentID uint64) (*api.DouyinCommentActionResponse, error) {
	// 查询此评论是否是本人发送的
	isComment, err := rdb.IsCommentCreatedByMyself(userID, videoID)
	// 这里因为使用了 string 存储，所以逻辑没有 ZSet 那么复杂
	if err != nil {
		isComment = db.IsCommentCreatedByMyself(userID, commentID)
	}

	// 非本人评论直接返回
	if !isComment {
		hlog.Error("service.comment.DeleteComment err:", errno.DeletePermissionError)
		return nil, errno.DeletePermissionError
	}

	// db 中会一起删除缓存数据
	dbc, err := db.DeleteCommentByID(videoID, commentID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}

	// 查询用户评论数据
	dbu, err := rdb.GetUserInfo(userID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		dbu, err = db.SelectUserByID(userID)
		if err != nil {
			hlog.Error("service.comment.DeleteComment err:", err.Error())
			return nil, err
		}
	}

	// 这里的 isFollow 直接返回 false ，因为评论人自己当然不能关注自己
	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment(dbc, dbu, false),
	}, nil
}

func GetCommentList(ctx context.Context, userID, videoID uint64) (*api.DouyinCommentListResponse, error) {
	logTag := "service.comment.GetCommentList err:"
	// userID可能为0，因为可能存在不登录也能查看视频评论的需求，但是videoID一定得为真实存在的ID
	// 使用布隆过滤器判断视频ID是否存在
	if !global.VideoIDBloomFilter.TestString(strconv.FormatUint(videoID, 10)) {
		hlog.Error(logTag, "布隆过滤器拦截")
		return nil, errno.UserRequestParameterError
	}

	// 加分布式锁
	lock := rdb.NewUserKeyLock(userID, constant.CommentRedisZSetPrefix)
	// 如果 redis 不可用，应该使用程序代码之外的方式进行限流
	_ = lock.Lock(ctx, global.CommentRC)

	// 获取评论基本数据
	cIDList, err := rdb.GetCommentIDByVideoID(videoID)
	if err != nil {
		hlog.Error(logTag, err)
		cIDList, err = db.SelectCommentIDByVideoID(videoID)
		if err != nil {
			hlog.Error(logTag, err)
			return nil, err
		}
	}

	cInfoList := make([]*rdb.CommentInfo, len(cIDList))
	lostCIDList := make([]uint64, 0, len(cIDList))
	for i := 0; i < len(cIDList); i++ {
		info, err := rdb.GetCommentInfo(cIDList[i])
		if err != nil {
			hlog.Error(logTag, err)
			// 没查询到就先记录起来
			lostCIDList = append(lostCIDList, cIDList[i])
		}
		// 如果 info 为 nil 就当添加一个占位用的，后续查完一个个填坑
		cInfoList[i] = info
	}

	// 一次性查询完剩下的
	lostCInfoList, err := db.SelectCommentInfoByCommentIDList(lostCIDList)
	if err != nil {
		hlog.Error(logTag, err)
		return nil, err
	}
	i := 0
	for _, data := range lostCInfoList {
		ci := &rdb.CommentInfo{
			ID:          data.ID,
			VideoID:     data.VideoID,
			UserID:      data.UserID,
			Content:     data.Content,
			CreatedTime: float64(data.CreatedTime.UnixMilli()),
		}
		for ; i < len(cInfoList); i++ {
			// 找到坑了则填入
			if cInfoList[i] == nil {
				cInfoList[i] = ci
				break
			}
		}
		// 顺手设置缓存
		err = rdb.SetCommentInfo(ci)
		if err != nil {
			hlog.Error(logTag, err)
		}
	}

	// 根据用户ID查询用户信息
	lostUIDList := make([]uint64, 0, len(cInfoList))
	uInfoList := make([]*model.User, len(cInfoList))
	for i := 0; i < len(cInfoList); i++ {
		ui, err := rdb.GetUserInfo(cInfoList[i].UserID)
		if err != nil {
			hlog.Error(logTag, err)
			lostUIDList = append(lostUIDList, cInfoList[i].UserID)
		}
		// 如果 info 为 nil 就当添加一个占位用的，后续查完一个个填坑
		uInfoList[i] = ui
	}

	// 还是一次性查询完剩下的
	// 这里需要注意一个多个评论可能对应同一个用户
	lostUInfoList, err := db.SelectUserByIDList(lostUIDList)
	if err != nil {
		hlog.Error(logTag, err)
		return nil, err
	}

	i = 0
	for _, data := range lostUInfoList {
		for ; i < len(uInfoList); i++ {
			// 找到坑了则填入
			if uInfoList[i] == nil && cInfoList[i].UserID == data.ID {
				uInfoList[i] = data
			}
		}
		// 顺手设置缓存
		err = rdb.SetUserInfo(data)
		if err != nil {
			hlog.Error(logTag, err)
		}
	}

	// 获取用户关注列表ID，判断是否关注并返回
	followUserIDSet, err := rdb.GetFollowUserIDSet(userID)
	if err != nil {
		// 缓存未命中就查询数据库
		hlog.Error(logTag, err)
		idList, err := db.SelectFollowUserIDSet(userID)
		if err != nil {
			hlog.Error(logTag, err)
		}
		for i := 0; i < len(idList); i++ {
			followUserIDSet[idList[i]] = struct{}{}
		}
		err = rdb.SetFollowUserIDSet(userID, idList)
		if err != nil {
			hlog.Error(logTag, err)
		}
	}

	// 解锁
	err = lock.Unlock(global.CommentRC)

	cList := make([]*api.Comment, len(cInfoList))
	for i := 0; i < len(cInfoList); i++ {
		isFollow := false
		if _, ok := followUserIDSet[cInfoList[i].UserID]; ok {
			isFollow = true
		}
		cList[i] = pack.ApiComment(cInfoList[i], uInfoList[i], isFollow)
	}

	return &api.DouyinCommentListResponse{
		StatusCode:  0,
		CommentList: cList,
	}, nil
}
