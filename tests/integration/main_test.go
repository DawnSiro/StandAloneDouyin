package integration

import(
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	password = "IntegrationTest!2023"
)

func TestDouyin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Douyin integration test")
}
