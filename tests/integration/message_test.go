package integration

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	util "douyin/test/testutil"
)

var _ = Describe("message test", func() {
	const (
		userA = "fortest-message-a"
		userB = "fortest-message-b"
	)

	var (
		idA    int64
		idB    int64
		tokenA string
		tokenB string
	)

	BeforeEach(func() {
		var err error
		idA, tokenA, err = util.GetUseridAndToken(userA, userA)
		Expect(err).To(BeNil())
		idB, tokenB, err = util.GetUseridAndToken(userB, userB)
		Expect(err).To(BeNil())
		// DoRelationAction(map[string]string{
		// 	"token":       tokenA,
		// 	"to_user_id":  fmt.Sprint("%d", idB),
		// 	"action_type": "1",
		// })
		// DoRelationAction(map[string]string{
		// 	"token":       tokenB,
		// 	"to_user_id":  fmt.Sprintf("%d", idA),
		// 	"action_type": "1",
		// })
	})

	// AfterEach(func() {
	// DoRelationAction(map[string]string{
	// 	"token":       tokenA,
	// 	"to_user_id":  fmt.Sprint("%d", idB),
	// 	"action_type": "2",
	// })
	// DoRelationAction(map[string]string{
	// 	"token":       tokenB,
	// 	"to_user_id":  fmt.Sprintf("%d", idA),
	// 	"action_type": "2",
	// })
	// })

	Describe("message action test", func() {
		const (
			path = "/douyin/message/action"
		)

		var (
			query = map[string]string{
				"action_type": "1",
				"content":     "This is a message for test.",
			}
		)

		Context("no friendship", func() {
			It("should fail", func() {
				Expect(idB).NotTo(BeZero())
				query["token"] = tokenA
				query["to_user_id"] = fmt.Sprintf("%d", idB)
				resp, err := http.Post(util.CreateURL(path, query), "", nil)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))

				respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(10400))) // TODO: better error
			})
		})

		Context("has friendship", func() {
			BeforeEach(func() {
				DoRelationAction(map[string]string{
					"token":       tokenA,
					"to_user_id":  fmt.Sprintf("%d", idB),
					"action_type": "1",
				})
				DoRelationAction(map[string]string{
					"token":       tokenB,
					"to_user_id":  fmt.Sprintf("%d", idA),
					"action_type": "1",
				})
				time.Sleep(1 * time.Second) // wait for MQ
			})

			AfterEach(func() {
				DoRelationAction(map[string]string{
					"token":       tokenA,
					"to_user_id":  fmt.Sprintf("%d", idB),
					"action_type": "2",
				})
				DoRelationAction(map[string]string{
					"token":       tokenB,
					"to_user_id":  fmt.Sprintf("%d", idA),
					"action_type": "2",
				})
			})

			It("A send to B", func() {
				query["token"] = tokenA
				query["to_user_id"] = fmt.Sprintf("%d", idB)
				resp, err := http.Post(util.CreateURL(path, query), "", nil)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))

				respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))
			})

			It("B send to A", func() {
				query["token"] = tokenB
				query["to_user_id"] = fmt.Sprintf("%d", idA)
				resp, err := http.Post(util.CreateURL(path, query), "", nil)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))

				respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))
			})
		})
	})

	Describe("message chat test", func() {
		const (
			path = "/douyin/message/chat"
		)

		var (
			queryA = make(map[string]string)
			queryB = make(map[string]string)
		)

		BeforeEach(func() {
			queryA["token"] = tokenA
			queryA["to_user_id"] = fmt.Sprintf("%d", idB)
			queryB["token"] = tokenB
			queryB["to_user_id"] = fmt.Sprintf("%d", idA)
		})

		It("should get message list of A", func() {
			resp, err := http.Get(util.CreateURL(path, queryA))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinMessageChatResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
		})

		It("should get message list of B", func() {
			resp, err := http.Get(util.CreateURL(path, queryB))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinMessageChatResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
		})

		It("A should = B", func() {
			resp, err := http.Get(util.CreateURL(path, queryA))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respDataA, err := util.GetDouyinResponse[util.DouyinMessageChatResponse](resp)
			Expect(err).To(BeNil())
			Expect(respDataA.StatusCode).To(Equal(int64(0)))

			resp, err = http.Get(util.CreateURL(path, queryB))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respDataB, err := util.GetDouyinResponse[util.DouyinMessageChatResponse](resp)
			Expect(err).To(BeNil())
			Expect(respDataB.StatusCode).To(Equal(int64(0)))

			Expect(len(respDataA.MessageList)).To(Equal(len(respDataB.MessageList)))
		})

		// Context("has friendship", func() {
		// 	BeforeEach(func() {
		// 		DoRelationAction(map[string]string{
		// 			"token":       tokenA,
		// 			"to_user_id":  fmt.Sprintf("%d", idB),
		// 			"action_type": "1",
		// 		})
		// 		DoRelationAction(map[string]string{
		// 			"token":       tokenB,
		// 			"to_user_id":  fmt.Sprintf("%d", idA),
		// 			"action_type": "1",
		// 		})
		// 	})

		// 	AfterEach(func() {
		// 		DoRelationAction(map[string]string{
		// 			"token":       tokenA,
		// 			"to_user_id":  fmt.Sprintf("%d", idB),
		// 			"action_type": "2",
		// 		})
		// 		DoRelationAction(map[string]string{
		// 			"token":       tokenB,
		// 			"to_user_id":  fmt.Sprintf("%d", idA),
		// 			"action_type": "2",
		// 		})
		// 	})

		// 	It("should success", func() {
		// 		resp, err := http.Get(util.CreateURL(path, queryA))
		// 		Expect(err).To(BeNil())
		// 		Expect(resp.StatusCode).To(Equal(200))
		// 		respData, err := util.GetDouyinResponse[util.DouyinMessageChatResponse](resp)
		// 		Expect(err).To(BeNil())
		// 		_ = respData
		// 	})
		// })
	})
})
