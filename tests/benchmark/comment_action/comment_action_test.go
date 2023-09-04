package comment_action

import (
	"fmt"
	"net/http"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gm "github.com/onsi/gomega/gmeasure"

	util "douyin/test/testutil"
)

var _ = Describe("comment action benchmark", func() {
	const (
		path     = "/douyin/comment/action/"
		password = "BenchmarkTest!2023"

		VideoID   = 3
		FanPrefix = "benchmark-commenter"
		NumFans   = 1000
		Content   = "comment for benchmark"
	)

	var (
		e *gm.Experiment
	)

	BeforeEach(func() {
		e = gm.NewExperiment("Test Experiment")
		AddReportEntry(e.Name, e)
	})

	It("[comment action]benchmark", func() {
		e.Sample(func(idx int) {
			n := fmt.Sprintf("%s-%d", FanPrefix, idx)
			_, token, err := util.GetUseridAndToken(n, password)
			Expect(err).To(BeNil())
			q := map[string]string{
				"token":        token,
				"video_id":     strconv.Itoa(VideoID),
				"action_type":  "1",
				"comment_text": Content,
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
