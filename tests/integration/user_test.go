package integration

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "github.com/go-sql-driver/mysql"

	util "douyin/test/testutil"
)

var _ = Describe("user test", func() {
	const uPath = "/douyin/user"

	Describe("register test", func() {
		const (
			path     = uPath + "/register"
			testuser = "fortest-register"
		)

		It("should register success", func() {
			query := map[string]string{
				"username": testuser,
				"password": "123456",
			}
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserRegisterResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
		})

		It("username exited", func() {
			query := map[string]string{
				"username": testuser,
				"password": "123456",
			}
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserRegisterResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10111)))

			err = util.DeleteUser(testuser)
			Expect(err).To(BeNil())
		})

		It("username too long", func() {
			longuser := "1111111111111111111111111111111111"
			Expect(len(longuser) > 32).To(BeTrue())
			query := map[string]string{
				"username": longuser,
				"password": "000000",
			}
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserRegisterResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10400))) // TODO: improve errno
		})

		It("password too long", func() {
			longpasswd := "1111111111111111111111111111111111"
			Expect(len(longpasswd) > 32).To(BeTrue())
			query := map[string]string{
				"username": "fortest",
				"password": longpasswd,
			}
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserRegisterResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10400)))
		})

		// TODO: password too short and not strong enough
	})

	Describe("login test", func() {
		const (
			path     = uPath + "/login"
			username = "fortest-login"
			password = "fortest-login"
		)
		var (
			userid int64
			// token  string
		)

		BeforeEach(func() {
			var err error
			userid, _, err = util.GetUseridAndToken(username, password)
			Expect(err).To(BeNil())
		})

		It("should login success", func() {
			query := map[string]string{
				"username": username,
				"password": password,
			}
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserLoginResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
			Expect(respData.UserID).To(Equal(userid))
			// Expect(respData.Token).To(Equal(token))
		})

		It("username not exist", func() {
			user := "user-not-exist"
			err := util.DeleteUser(user)
			Expect(err).To(BeNil())

			query := map[string]string{
				"username": user,
				"password": "hhhhhh",
			}
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserLoginResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10201)))
		})

		It("wrong password", func() {
			query := map[string]string{
				"username": username,
				"password": password + "0",
			}
			resp, err := http.Post(util.CreateURL(path, query), "", nil)
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserLoginResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10210)))
		})
	})

	Describe("get info test", func() {
		const (
			path     = uPath
			username = "fortest-userinfo"
		)
		var (
			id     int64
			userid string
			token  string
		)

		BeforeEach(func() {
			var err error
			id, token, err = util.GetUseridAndToken(username, username)
			Expect(err).To(BeNil())
			userid = fmt.Sprintf("%d", id)
		})

		It("should get info", func() {
			query := map[string]string{
				"user_id": userid,
				"token":   token,
			}
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(0)))
			Expect(respData.User.ID).To(Equal(id))
			Expect(respData.User.Name).To(Equal(username))
		})

		// It("user_id not exist", func() {
		// 	err := util.DeleteUser(username)
		// 	Expect(err).To(BeNil())

		// 	query := map[string]string{
		// 		"user_id": userid,
		// 		"token":   token,
		// 	}
		// 	resp, err := http.Get(util.CreateURL(path, query))
		// 	Expect(err).To(BeNil())
		// 	defer resp.Body.Close()
		// 	Expect(resp.StatusCode).To(Equal(200))

		// 	respData, err := util.GetDouyinResponse[util.DouyinUserResponse](resp)
		// 	Expect(err).To(BeNil())
		// 	Expect(respData.StatusCode).To(Equal(int64(20000)))
		// })

		It("wrong userid", func() {
			query := map[string]string{
				"user_id": userid + "0",
				"token":   token,
			}
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(20000)))
		})

		It("wrong token", func() {
			query := map[string]string{
				"user_id": userid,
				"token":   token + "0",
			}
			resp, err := http.Get(util.CreateURL(path, query))
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))

			respData, err := util.GetDouyinResponse[util.DouyinUserResponse](resp)
			Expect(err).To(BeNil())
			Expect(respData.StatusCode).To(Equal(int64(10220)))
		})
	})

})
