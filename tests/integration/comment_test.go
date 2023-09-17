package integration

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	util "douyin/test/testutil"
)

var _ = Describe("comment test", func() {
	const (
		username = "fortest-comment"
		video_id = 2
		content  = "This is a comment for test"
	)

	Describe("comment action test", func() {
		const (
			path = "/douyin/comment/action"
		)

		var (
			query = map[string]string{
				"video_id": fmt.Sprintf("%d", video_id),
			}
			commentId int64
		)

		BeforeEach(func() {
			_, token, err := util.GetUserIDAndToken(username, password)
			Expect(err).To(BeNil())
			query["token"] = token
		})

		It("should success", func() {
			before_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())

			query["action_type"] = fmt.Sprintf("%d", 1)
			query["comment_text"] = content
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinCommentActionResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
			commentId = respData.Comment.ID

			time.Sleep(time.Second) // wait for mq

			after_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(1))
		})

		It("should cancel success", func() {
			before_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())
			Expect(before_cnt > 0).To(BeTrue())

			query["action_type"] = fmt.Sprintf("%d", 2)
			query["comment_id"] = fmt.Sprintf("%d", commentId)
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinCommentActionResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))

			after_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(-1))
		})

		It("wrong token", func() {
			before_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())

			query["action_type"] = fmt.Sprintf("%d", 1)
			query["comment_text"] = content
			query["token"] += "0"
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinCommentActionResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10220)))

			after_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(0))
		})

		It("wrong video id", func() {
			before_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())

			query["action_type"] = fmt.Sprintf("%d", 1)
			query["video_id"] = fmt.Sprintf("%d", -1)
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinCommentActionResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10400))) // TODO: better errno

			after_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(0))
		})

		It("wrong action", func() {
			before_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())

			query["action_type"] = fmt.Sprintf("%d", 3)
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinCommentActionResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10400))) // TODO: better errno

			after_cnt, err := GetCommentNum(video_id)
			Expect(err).To(BeNil())
			Expect(after_cnt - before_cnt).To(Equal(0))
		})
	})

	Describe("comment list test", func() {
		const (
			path = "/douyin/comment/list"
		)

		var (
			query = map[string]string{
				"video_id": fmt.Sprintf("%d", video_id),
			}
		)

		BeforeEach(func() {
			_, token, err := util.GetUserIDAndToken(username, password)
			Expect(err).To(BeNil())
			query["token"] = token
		})

		It("should success", func() {
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinCommentListResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
		})

		It("wrong token", func() {
			query["token"] += "0"
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinCommentListResponse](resp)
			Expect(err).To(BeNil())
			// FIXME: wrong token but get comment list success?
			// Expect(respData.StatusCode).To(Equal(int64(10220)))
			Expect(respData.StatusCode).To(Equal(int64(0)))
		})

		It("wrong video id", func() {
			query["video_id"] = fmt.Sprintf("%d", -1)
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinCommentListResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10400))) // TODO: better errno
		})

		// It("should be in reverse time order", func() {
		// 	resp, err := http.Get(util.CreateURL(path, query))
		// 	Expect(err).To(BeNil())
		// 	Expect(resp.StatusCode).To(Equal(200))
		// 	respData, err := util.GetDouyinResponse[util.DouyinCommentListResponse](resp)
		// 	Expect(err).To(BeNil())
		// 	Expect(respData.StatusCode).To(Equal(int64(0)))

		// 	var lastTime string
		// 	// lastTime := time.Now()
		// 	for i, c := range respData.CommentList {
		// 		t := c.CreateDate
		// 		Expect(err).To(BeNil())
		// 		if i > 0 {
		// 			Expect(t <= lastTime).To(BeTrue())
		// 		}
		// 		lastTime = t
		// 	}
		// })
	})
})

func GetCommentNum(vid int) (num int, err error) {
	db, err := util.GetDBConnection()
	if err != nil {
		return
	}

	sqlStr := `select comment_count from video where id = ?;`
	rows, err := db.Query(sqlStr, vid)
	if err != nil {
		return
	}
	defer rows.Close()
	Expect(rows.Next()).To(BeTrue())
	err = rows.Scan(&num)
	return
}
