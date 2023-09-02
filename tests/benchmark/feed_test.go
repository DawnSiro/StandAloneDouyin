package benchmark

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gm "github.com/onsi/gomega/gmeasure"

	util "douyin/test/testutil"
)

var _ = Describe("feed test", func() {
	const (
		path = "/douyin/feed"
	)
	var (
		e *gm.Experiment
	)
	BeforeEach(func() {
		e = gm.NewExperiment("Test Experiment")
		AddReportEntry(e.Name, e)
	})

	Context("no token", func() {
		It("should feed success", func() {
			e.Sample(func(idx int) {
				e.MeasureDuration("", func() {
					resp, err := http.Get(util.CreateURL(path, nil))
					Expect(err).To(BeNil())
					defer resp.Body.Close()
					Expect(resp.StatusCode).To(Equal(200))
				})
			}, gm.SamplingConfig{N: Times, NumParallel: Threads})
		})
	})

	Context("has token", func() {
		const (
			username = "fortest-feed"
			password = "fortest-feed"
		)

		var (
			query = make(map[string]string)
		)

		BeforeEach(func() {
			_, token, err := util.GetUseridAndToken(username, password)
			Expect(err).To(BeNil())
			query["token"] = token
		})

		AfterEach(func() {
			delete(query, "token")
		})

		It("should feed success", func() {
			e.Sample(func(idx int) {
				e.MeasureDuration("", func() {
					resp, err := http.Get(util.CreateURL(path, query))
					Expect(err).To(BeNil())
					defer resp.Body.Close()
					Expect(resp.StatusCode).To(Equal(200))
				})
			}, gm.SamplingConfig{N: Times, NumParallel: Threads})
		})
	})
})
