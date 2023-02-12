package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
)

func Videos(video []*db.Video) ([]*api.Video, error) {
	res := make([]*api.Video, 0)
	return res, nil
}
