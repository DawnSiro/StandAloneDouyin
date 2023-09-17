package main

import (
	util "douyin/test/testutil"
)

func main() {
	id, _, err := util.GetUserIDAndToken("jmeter-follower-list", "123456")
	if err != nil {
		panic(err)
	}
	err = DoFollower(100, id, "jmeter-fan", "1")
	if err != nil {
		panic(err)
	}
}
