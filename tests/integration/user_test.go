package integration

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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
				"password": password,
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
				"password": password,
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
				"password": password,
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
			Expect(respData.StatusCode).To(Equal(int64(10122)))
		})

		// TODO: password too short and not strong enough
	})

	Describe("login test", func() {
		const (
			path     = uPath + "/login"
			username = "fortest-login"
			password = password
		)
		var (
			userid int64
			// token  string
		)

		BeforeEach(func() {
			var err error
			userid, _, err = util.GetUserIDAndToken(username, password)
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
				"password": password,
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
			id, token, err = util.GetUserIDAndToken(username, password)
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
				"user_id": userid + "000",
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
