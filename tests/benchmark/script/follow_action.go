package main

import (
	util "douyin/test/testutil"
	"fmt"
)

const (
	StarName  = "jmeter-follow-action-star"
	FanPrefix = "jmeter-follow-action-fan"
	NumFans   = 1000

	password = "jmeter-benchmark"
)

func main() {
	_, _, err := util.GetUseridAndToken(StarName, password)
	assert(err)
	fmt.Println("Create star ok")
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
