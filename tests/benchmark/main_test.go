package benchmark

import(
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDouyin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Douyin benchmark test")
}
