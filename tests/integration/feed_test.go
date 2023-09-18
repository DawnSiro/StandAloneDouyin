package integration

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	util "douyin/test/testutil"
)

var _ = Describe("/douyin/feed api request", func() {
	const path = "/douyin/feed"

	Context("no token", func() {
		It("should feed success", func() {
			resp, err := http.Get(util.CreateURL(path, nil))
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))
		})

		It("should less than 30", func() {
			resp, err := http.Get(util.CreateURL(path, nil))
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinFeedResponse](resp)
			Expect(err).To(BeNil())
			Expect(len(respData.VideoList) <= 30).To(BeTrue())
		})

		// It("should reverse order", func() {
		// 	resp, err := http.Get(CreateURL(path, nil))
		// 	Expect(err).To(BeNil())
		// 	defer resp.Body.Close()
		// 	Expect(resp.StatusCode).To(Equal(200))

		// 	respData, err := getFeedResponse(resp)
		// 	Expect(err).To(BeNil())
		// 	// last := 0
		// })

		// It("should filter with latest time", func() {
		// 	resp, err := http.Get(CreateURL(path, nil))
		// 	Expect(err).To(BeNil())
		// 	defer resp.Body.Close()
		// 	Expect(resp.StatusCode).To(Equal(200))

		// 	respData, err := getFeedResponse(resp)
		// 	latestTime := respData.NextTime

		// 	query := map[string]string{
		// 		"latest_time": fmt.Sprintf("%d", *latestTime),
		// 	}
		// 	resp, err = http.Get(CreateURL(path, query))
		// 	Expect(err).To(BeNil())
		// 	defer resp.Body.Close()
		// 	Expect(resp.StatusCode).To(Equal(200))
		// 	respData, err = getFeedResponse(resp)
		// })
	})

	Context("has token", func() {
		const (
			username = "fortest-feed"
		)

		var (
			query = make(map[string]string)
		)

		BeforeEach(func() {
			_, token, err := util.GetUserIDAndToken(username, password)
			Expect(err).To(BeNil())
			query["token"] = token
		})

		AfterEach(func() {
			delete(query, "token")
		})

		It("should feed success", func() {
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))
		})

		It("should less than 30", func() {
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinFeedResponse](resp)
			Expect(err).To(BeNil())
			Expect(len(respData.VideoList) <= 30).To(BeTrue())
		})
	})
})
