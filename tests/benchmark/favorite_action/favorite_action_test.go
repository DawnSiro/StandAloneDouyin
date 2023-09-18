package favorite_action

import (
	"fmt"
	"net/http"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gm "github.com/onsi/gomega/gmeasure"

	util "douyin/test/testutil"
)

var _ = Describe("favorite action benchmark", func() {
	const (
		path     = "/douyin/favorite/action/"
		password = "BenchmarkTest!2023"

		VideoID   = 2
		FanPrefix = "benchmark-favorite-fan"
		NumFans   = 1000
	)

	var (
		e *gm.Experiment
	)

	BeforeEach(func() {
		e = gm.NewExperiment("Test Experiment")
		AddReportEntry(e.Name, e)
	})

	It("[favorite action]benchmark", func() {
		e.Sample(func(idx int) {
			n := fmt.Sprintf("%s-%d", FanPrefix, idx)
			_, token, err := util.GetUserIDAndToken(n, password)
			Expect(err).To(BeNil())
			q := map[string]string{
				"token":       token,
				"video_id":    strconv.Itoa(VideoID),
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
