package main

import (
	"fmt"
	"net/http"

	util "douyin/test/testutil"
)

const (
	password = "BenchmarkTest!2023"

	userA = "benchmark-friend-userA"
	userB = "benchmark-friend-userB"
)

func AFollowB(token string, toUserID int64) {
	q := map[string]string{
		"token":       token,
		"to_user_id":  fmt.Sprintf("%d", toUserID),
		"action_type": "1",
	}
	resp, err := http.Post(util.CreateURL("/douyin/relation/action/", q), "", nil)
	assert(err)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("return code != 200")
	}
}

func main() {
	idA, tokenA, err := util.GetUserIDAndToken(userA, password)
	assert(err)
	idB, tokenB, err := util.GetUserIDAndToken(userB, password)
	assert(err)

	AFollowB(tokenA, idB)
	AFollowB(tokenB, idA)
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
