package follow_action

import (
	"flag"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	Times   int
	Threads int
)

func init() {
	flag.IntVar(&Times, "repeat", 1000, "-times=<Repeat times>")
	flag.IntVar(&Threads, "thread", 5, "-threads=<Thread numbers>")
}

func TestDouyin(t *testing.T) {
	flag.Parse()
	fmt.Printf("\033[33mrepeat times = %v, thread numbers = %v\n\033[0m", Times, Threads)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Douyin benchmark test for follow action")
}
