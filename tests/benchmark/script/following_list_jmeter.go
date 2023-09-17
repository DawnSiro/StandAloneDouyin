package main

import (
	util "douyin/test/testutil"
)

func main() {
	_, token, err := util.GetUserIDAndToken("jmeter-following-list", "123456")
	if err != nil {
		panic(err)
	}
	DoFollowing(100, "jmeter-star", token, "1")
}
