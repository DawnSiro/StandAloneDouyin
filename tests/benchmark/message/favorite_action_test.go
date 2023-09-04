package message

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gm "github.com/onsi/gomega/gmeasure"

	util "douyin/test/testutil"
)

type UserInfo struct {
	Uid   int64
	Token string
}

var _ = Describe("favorite action benchmark", func() {
	const (
		path          = "/douyin/message/action/"
		password      = "BenchmarkTest!2023"
		messagePrefix = "message for benchmark"
	)

	var (
		userA = UserInfo{2002, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MjAwMiwiZXhwIjoxNjkzODYzNTIwLCJuYmYiOjE2OTM4MjAzMjAsImlhdCI6MTY5MzgyMDMyMH0.Dr_vX7FFIzMIeJqWi_X8xI6VKAO7fjKHB-oMcl_lEGk"}
		userB = UserInfo{2003, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MjAwMywiZXhwIjoxNjkzODYzNTgwLCJuYmYiOjE2OTM4MjAzODAsImlhdCI6MTY5MzgyMDM4MH0.ockEQKRI2fPSyTYuvrjckOIdTnu10wm0h6Ma74283GQ"}
		e     *gm.Experiment
	)

	BeforeEach(func() {
		e = gm.NewExperiment("Test Experiment")
		AddReportEntry(e.Name, e)
	})

	It("[message]benchmark", func() {
		e.Sample(func(idx int) {
			qa2b := map[string]string{
				"token":       userA.Token,
				"to_user_id":  fmt.Sprintf("%d", userB.Uid),
				"action_type": "1",
				"content":     fmt.Sprintf("%s (A2B-%d)", messagePrefix, idx),
			}
			qb2a := map[string]string{
				"token":       userB.Token,
				"to_user_id":  fmt.Sprintf("%d", userA.Uid),
				"action_type": "1",
				"content":     fmt.Sprintf("%s (B2A-%d)", messagePrefix, idx),
			}

			e.MeasureDuration("AtoB", func() {
				resp, err := http.Post(util.CreateURL(path, qa2b), "", nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(200))
			})

			e.MeasureDuration("BtoA", func() {
				resp, err := http.Post(util.CreateURL(path, qb2a), "", nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(200))
			})
		}, gm.SamplingConfig{N: Times, NumParallel: Threads})
	})
})
