package main

import (
	"fmt"

	util "douyin/test/testutil"
)

const (
	FanPrefix = "benchmark-commenter"
	NumFans   = 1000

	password = "BenchmarkTest!2023"
)

func main() {
	for i := 0; i < NumFans; i++ {
		n := fmt.Sprintf("%s-%d", FanPrefix, i)
		_, _, err := util.GetUseridAndToken(n, password)
		assert(err)
	}
	fmt.Println("Create fans ok")
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
