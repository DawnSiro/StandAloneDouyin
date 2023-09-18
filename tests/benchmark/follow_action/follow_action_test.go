package follow_action

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gm "github.com/onsi/gomega/gmeasure"

	util "douyin/test/testutil"
)

var _ = Describe("follow action benchmark", func() {
	const (
		path     = "/douyin/relation/action/"
		password = "jmeter-benchmark"

		StarName   = "jmeter-follow-action-star"
		StarUserid = 1
		FanPrefix  = "jmeter-follow-action-fan"
		NumFans    = 1000
	)

	var (
		e *gm.Experiment
	)

	BeforeEach(func() {
		e = gm.NewExperiment("Test Experiment")
		AddReportEntry(e.Name, e)
	})

	It("[follow action]benchmark", func() {
		e.Sample(func(idx int) {
			n := fmt.Sprintf("%s-%d", FanPrefix, idx)
			_, token, err := util.GetUserIDAndToken(n, password)
			Expect(err).To(BeNil())
			q := map[string]string{
				"token":       token,
				"to_user_id":  "1",
				"action_type": "1",
			}

			e.MeasureDuration("", func() {
				resp, err := http.Post(util.CreateURL(path, q), "", nil)
				Expect(err).To(BeNil())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(200))
			})
		}, gm.SamplingConfig{N: Times, NumParallel: Threads})
	})
})
