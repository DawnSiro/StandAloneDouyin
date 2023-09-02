package integration

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	util "douyin/test/testutil"
)

var _ = Describe("favorite test", func() {
	const (
		username = "fortest-favorite"
		password = "fortest-favorite"
	)

	var (
		query = make(map[string]string)
	)

	Describe("favorite action test", func() {
		const (
			video_id = 1
			path     = "/douyin/favorite/action"
		)

		BeforeEach(func() {
			_, token, err := util.GetUseridAndToken(username, password)
			Expect(err).To(BeNil())
			query["token"] = token
			query["video_id"] = fmt.Sprintf("%d", video_id)
		})

		It("should success", func() {
			before_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())

			query["action_type"] = fmt.Sprintf("%d", 1)
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))

			after_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(1))
			
			query["action_type"] = fmt.Sprintf("%d", 2)
			resp, err = http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err = util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
		})

		It("should cancel success", func() {
			query["action_type"] = fmt.Sprintf("%d", 1)
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))

			before_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())
			Expect(before_cnt).NotTo(BeZero())

			query["action_type"] = fmt.Sprintf("%d", 2)
			resp, err = http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err = util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))

			after_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(-1))
		})

		// It("double favorite", func() {
		// 	before_cnt, err := GetFavoriteNum(video_id)
		// 	Expect(err).To(BeNil())

		// 	query["action_type"] = fmt.Sprintf("%d", 1)
		// 	resp, err := http.Post(CreateURL(path+"/action", query), "", nil)
		// 	Expect(err).To(BeNil())
		// 	Expect(resp.StatusCode).To(Equal(200))
		// 	respData, err := GetDouyinResponse[util.DouyinSimpleResponse](resp)
		// 	Expect(err).To(BeNil())
		// 	Expect(respData.StatusCode).To(Equal(int64(0)))
		// 	resp, err = http.Post(CreateURL(path+"/action", query), "", nil)
		// 	Expect(err).To(BeNil())
		// 	Expect(resp.StatusCode).To(Equal(200))
		// 	respData, err = GetDouyinResponse[util.DouyinSimpleResponse](resp)
		// 	Expect(err).To(BeNil())
		// 	Expect(respData.StatusCode).To(Equal(int64(0)))

		// 	after_cnt, err := GetFavoriteNum(video_id)
		// 	Expect(err).To(BeNil())
		// 	Expect(after_cnt - before_cnt).To(Equal(1))

		// 	// 恢复原来的点赞数据
		// 	err = RecoverFavoriteData(int(userid), video_id, 1)
		// 	Expect(err).To(BeNil())
		// })

		It("wrong token", func() {
			before_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())

			query["action_type"] = fmt.Sprintf("%d", 1)
			query["token"] += "0"
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10220)))

			after_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(0))
		})

		It("wrong video id", func() {
			before_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())

			query["action_type"] = fmt.Sprintf("%d", 1)
			query["video_id"] = fmt.Sprintf("%d", -1)
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10400))) // TODO: better errno

			after_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(0))
		})

		It("wrong action", func() {
			before_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())

			query["action_type"] = fmt.Sprintf("%d", 3)
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10400))) // TODO: better errno

			after_cnt, err := GetFavoriteNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(0))
		})
	})

	Describe("favorite list test", func() {
		const (
			path = "/douyin/favorite/list"
		)

		BeforeEach(func() {
			userid, token, err := util.GetUseridAndToken(username, password)
			Expect(err).To(BeNil())
			query["user_id"] = fmt.Sprintf("%d", userid)
			query["token"] = token
		})

		It("should success", func() {
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinFavoriteListResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
		})

		It("should equal 3", func() {
			err := doFavoriteAction(query["token"])
			Expect(err).To(BeNil())

			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinFavoriteListResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
			Expect(len(respData.VideoList)).To(Equal(3))

			err = cancelFavoriteAction(query["token"])
			Expect(err).To(BeNil())
		})

		It("wrong token", func() {
			query["token"] += "0"
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinFavoriteListResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10220)))
		})

		It("wrong user id", func() {
			query["user_id"] = "-1"
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinFavoriteListResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10400)))
		})

	})
})

func GetFavoriteNum(vid int) (num int, err error) {
	db, err := util.GetDBConnection()
	if err != nil {
		return
	}

	sqlStr := `select favorite_count from video where id = ?;`
	rows, err := db.Query(sqlStr, vid)
	if err != nil {
		return
	}
	defer rows.Close()
	Expect(rows.Next()).To(BeTrue())
	err = rows.Scan(&num)
	return
}

func doFavoriteAction(token string) (err error) {
	q := map[string]string{
		"token":       token,
		"video_id":    "1",
		"action_type": "1",
	}
	_, err = http.Post(util.CreateURL("/douyin/favorite/action", q), "", nil)
	if err != nil {
		return
	}
	q["video_id"] = "2"
	_, err = http.Post(util.CreateURL("/douyin/favorite/action", q), "", nil)
	if err != nil {
		return
	}
	q["video_id"] = "3"
	_, err = http.Post(util.CreateURL("/douyin/favorite/action", q), "", nil)
	return
}

func cancelFavoriteAction(token string) (err error) {
	q := map[string]string{
		"token":       token,
		"video_id":    "1",
		"action_type": "2",
	}
	_, err = http.Post(util.CreateURL("/douyin/favorite/action", q), "", nil)
	if err != nil {
		return
	}
	q["video_id"] = "2"
	_, err = http.Post(util.CreateURL("/douyin/favorite/action", q), "", nil)
	if err != nil {
		return
	}
	q["video_id"] = "3"
	_, err = http.Post(util.CreateURL("/douyin/favorite/action", q), "", nil)
	return
}

// func SetFavoriteNum(vid, num int) (err error) {
// 	db, err := util.GetDBConnection()
// 	if err != nil {
// 		return
// 	}

// 	sqlStr := `update video set favorite_count=? where id = ?;`
// 	_, err = db.Exec(sqlStr, num, vid)
// 	return
// }

// func RecoverFavoriteData(uid, vid, cnt int) (err error) {
// 	db, err := GetDBConnection()
// 	if err != nil {
// 		return
// 	}

// 	q := `update video set favorite_count=? where id=?`
// 	_, err = db.Exec(q, cnt, vid)
// 	if err != nil {
// 		return
// 	}
// 	q = `delete from user_favorite_video where user_id=?;`
// 	_, err = db.Exec(q, uid)
// 	if err != nil {
// 		return
// 	}
// 	q = `delete from user where id=?;`
// 	_, err = db.Exec(q, uid)
// 	return
// }
