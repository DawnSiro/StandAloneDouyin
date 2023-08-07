package integration

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	util "douyin/test/testutil"
)

var _ = Describe("relation test", func() {
	const (
		up  = "fortest-up"
		fan = "fortest-fan"
	)

	var (
		upid  int64
		fanid int64
	)

	Describe("action test", func() {
		const (
			path = "/douyin/relation/action"
		)

		var (
			query = make(map[string]string)
		)

		BeforeEach(func() {
			fanid_, token, err := util.GetUseridAndToken(fan, fan)
			Expect(err).To(BeNil())
			upid_, _, err := util.GetUseridAndToken(up, up)
			Expect(err).To(BeNil())
			fanid = fanid_
			upid = upid_
			query["token"] = token
			query["to_user_id"] = fmt.Sprintf("%d", upid)
		})

		It("should follow success", func() {
			fansBefore, err := GetFansNum(upid)
			Expect(err).To(BeNil())
			upsBefore, err := GetFollowingNum(fanid)
			Expect(err).To(BeNil())

			query["action_type"] = "1"
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))

			fansAfter, err := GetFansNum(upid)
			Expect(err).To(BeNil())
			Expect(fansAfter - fansBefore).To(Equal(1))
			upsAfter, err := GetFollowingNum(fanid)
			Expect(err).To(BeNil())
			Expect(upsAfter - upsBefore).To(Equal(1))

			// cancel the follow action
			query["action_type"] = "2"
			resp, err = http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err = util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
		})

		It("should cancel success", func() {
			// follow firstly
			query["action_type"] = "1"
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))

			fansBefore, err := GetFansNum(upid)
			Expect(err).To(BeNil())
			upsBefore, err := GetFollowingNum(fanid)
			Expect(err).To(BeNil())

			query["action_type"] = "2"
			resp, err = http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err = util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))

			fansAfter, err := GetFansNum(upid)
			Expect(err).To(BeNil())
			Expect(fansAfter - fansBefore).To(Equal(-1))
			upsAfter, err := GetFollowingNum(fanid)
			Expect(err).To(BeNil())
			Expect(upsAfter - upsBefore).To(Equal(-1))
		})

		It("wrong token", func() {
			query["action_type"] = "1"
			query["token"] += "0"
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10220)))
		})
	})

	Describe("follow list test", func() {
		const (
			path = "/douyin/relation/follow/list"
		)

		var (
			query = make(map[string]string)
		)

		BeforeEach(func() {
			id, token, err := util.GetUseridAndToken(fan, fan)
			Expect(err).To(BeNil())
			query["user_id"] = fmt.Sprintf("%d", id)
			query["token"] = token
		})

		Context("no following", func() {
			It("should get following list", func() {
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinRelationFollowListResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))

				cnt, err := GetFollowingNum(fanid)
				Expect(err).To(BeNil())
				for _, u := range respData.UserList {
					if u.IsFollow {
						cnt--
					}
				}
				Expect(cnt).To(BeZero())
			})

			It("wrong token", func() {
				query["token"] += "0"
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(10220)))
			})
		})

		Context("has following", func() {
			BeforeEach(func() {
				err := actionFollowing(query["token"], "1")
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				err := actionFollowing(query["token"], "2")
				Expect(err).To(BeNil())
			})

			It("should get following list", func() {
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinRelationFollowListResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))

				cnt, err := GetFollowingNum(fanid)
				Expect(err).To(BeNil())
				for _, u := range respData.UserList {
					if u.IsFollow {
						cnt--
					}
				}
				Expect(cnt).To(BeZero())
			})

			It("wrong token", func() {
				query["token"] += "0"
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(10220)))
			})
		})
	})

	Describe("follower list test", func() {
		const (
			path = "/douyin/relation/follower/list"
		)

		var (
			query = make(map[string]string)
		)

		BeforeEach(func() {
			id, token, err := util.GetUseridAndToken(up, up)
			Expect(err).To(BeNil())
			query["user_id"] = fmt.Sprintf("%d", id)
			query["token"] = token
		})

		Context("no follower", func() {
			It("should get follower list", func() {
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinRelationFollowListResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))

				cnt, err := GetFollowerNum(upid)
				Expect(err).To(BeNil())
				for _, u := range respData.UserList {
					if u.IsFollow {
						cnt--
					}
				}
				Expect(cnt).To(BeZero())
			})

			It("wrong token", func() {
				query["token"] += "0"
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(10220)))
			})
		})

		Context("has follower", func() {
			BeforeEach(func() {
				err := actionFollower(upid, "1")
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				err := actionFollower(upid, "2")
				Expect(err).To(BeNil())
			})

			It("should get follower list", func() {
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinRelationFollowListResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))

				cnt, err := GetFollowerNum(upid)
				Expect(err).To(BeNil())
				for _, u := range respData.UserList {
					if u.IsFollow {
						cnt--
					}
				}
				Expect(cnt).To(BeZero())
			})

			It("wrong token", func() {
				query["token"] += "0"
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(10220)))
			})
		})

	})

	Describe("friend list test", func() {
		const (
			path = "/douyin/relation/friend/list"
		)

		Context("no friendship", func() {
			var (
				query = make(map[string]string)
			)

			BeforeEach(func() {
				id, token, err := util.GetUseridAndToken(up, up)
				Expect(err).To(BeNil())
				query["user_id"] = fmt.Sprintf("%d", id)
				query["token"] = token
			})

			It("should fail", func() {
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinRelationFriendListResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))

				cnt, err := GetFollowerNum(upid)
				Expect(cnt).To(BeZero())
				Expect(err).To(BeNil())
				for _, u := range respData.UserList {
					if u.IsFollow {
						cnt--
					}
				}
				Expect(cnt).To(BeZero())
			})

			It("wrong token", func() {
				query["token"] += "0"
				resp, err := http.Get(util.CreateURL(path, query))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinRelationFriendListResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(10220)))
			})
		})

		Context("has friendship", func() {
			var (
				queryA = make(map[string]string)
				queryB = make(map[string]string)
			)

			BeforeEach(func() {
				id, token, err := util.GetUseridAndToken(up, up)
				Expect(err).To(BeNil())
				queryA["user_id"] = fmt.Sprintf("%d", id)
				queryA["token"] = token

				id, token, err = util.GetUseridAndToken(fan, fan)
				Expect(err).To(BeNil())
				queryB["user_id"] = fmt.Sprintf("%d", id)
				queryB["token"] = token

				DoRelationAction(map[string]string{
					"token":       queryA["token"],
					"to_user_id":  queryB["user_id"],
					"action_type": "1",
				})
				DoRelationAction(map[string]string{
					"token":       queryB["token"],
					"to_user_id":  queryA["user_id"],
					"action_type": "1",
				})
			})

			AfterEach(func() {
				DoRelationAction(map[string]string{
					"token":       queryA["token"],
					"to_user_id":  queryB["user_id"],
					"action_type": "2",
				})
				DoRelationAction(map[string]string{
					"token":       queryB["token"],
					"to_user_id":  queryA["user_id"],
					"action_type": "2",
				})
			})

			It("should get friend list of A", func() {
				resp, err := http.Get(util.CreateURL(path, queryA))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinRelationFriendListResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))
				Expect(respData.UserList[0].Name).To(Equal(fan))
			})

			It("should get friend list of B", func() {
				resp, err := http.Get(util.CreateURL(path, queryB))
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				respData, err := util.GetDouyinResponse[util.DouyinRelationFriendListResponse](resp)
				Expect(err).To(BeNil())
				Expect(respData.StatusCode).To(Equal(int64(0)))
				Expect(respData.UserList[0].Name).To(Equal(up))
			})

			// It("wrong token", func() {
			// 	query["token"] += "0"
			// 	resp, err := http.Get(util.CreateURL(path, query))
			// 	Expect(err).To(BeNil())
			// 	Expect(resp.StatusCode).To(Equal(200))
			// 	respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
			// 	Expect(err).To(BeNil())
			// 	Expect(respData.StatusCode).To(Equal(int64(10220)))
			// })
		})

	})
})

func DoRelationAction(q map[string]string) {
	resp, err := http.Post(util.CreateURL("/douyin/relation/action", q), "", nil)
	Expect(err).To(BeNil())
	Expect(resp.StatusCode).To(Equal(200))
	respData, err := util.GetDouyinResponse[util.DouyinSimpleResponse](resp)
	Expect(err).To(BeNil())
	Expect(respData.StatusCode).To(Equal(int64(0)))
}

func GetFansNum(id int64) (res int, err error) {
	db, err := util.GetDBConnection()
	if err != nil {
		return
	}

	q := `select follower_count from user where id = ?`
	rows, err := db.Query(q, id)
	if err != nil {
		return
	}
	rows.Next()
	err = rows.Scan(&res)
	return
}

func GetFollowingNum(id int64) (res int, err error) {
	db, err := util.GetDBConnection()
	if err != nil {
		return
	}

	q := `select following_count from user where id = ?`
	rows, err := db.Query(q, id)
	if err != nil {
		return
	}
	rows.Next()
	err = rows.Scan(&res)
	return
}

func GetFollowerNum(id int64) (res int, err error) {
	db, err := util.GetDBConnection()
	if err != nil {
		return
	}

	q := `select follower_count from user where id = ?`
	rows, err := db.Query(q, id)
	if err != nil {
		return
	}
	rows.Next()
	err = rows.Scan(&res)
	return
}

func actionFollowing(token, action string) (err error) {
	query := map[string]string{
		"token":       token,
		"action_type": action,
	}
	for i := 1; i <= 10; i++ {
		query["to_user_id"] = fmt.Sprintf("%d", i)
		_, err = http.Post(util.CreateURL("/douyin/relation/action", query), "", nil)
		if err != nil {
			return
		}
	}
	return
}

func actionFollower(upid int64, action string) error {
	const user = "fortest-fans"
	query := map[string]string{
		"to_user_id":  fmt.Sprintf("%d", upid),
		"action_type": action,
	}
	for i := 0; i < 10; i++ {
		f := user + fmt.Sprintf("%d", i)
		_, token, err := util.GetUseridAndToken(f, f)
		if err != nil {
			return err
		}
		query["token"] = token
		_, err = http.Post(util.CreateURL("/douyin/relation/action", query), "", nil)
		if err != nil {
			return err
		}
	}
	return nil
}
