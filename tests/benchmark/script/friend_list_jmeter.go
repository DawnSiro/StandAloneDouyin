package main

import (
	util "douyin/test/testutil"
	"fmt"
	"net/http"
)

func main() {
	id, token, err := util.GetUserIDAndToken("jmeter-friend-list", "123456")
	if err != nil {
		panic(err)
	}
	err = DoFollowing(100, "jmeter-friend", token, "1")
	if err != nil {
		panic(err)
	}
	err = DoFollower(100, id, "jmeter-friend", "1")
	if err != nil {
		panic(err)
	}
}

func DoFollowing(num int, prefix, token, action string) (err error) {
	query := map[string]string{
		"token":       token,
		"action_type": action,
	}
	for i := 0; i < num; i++ {
		u := prefix + fmt.Sprintf("%d", i)
		var uid int64
		uid, _, err = util.GetUserIDAndToken(u, u)
		if err != nil {
			return
		}
		query["to_user_id"] = fmt.Sprintf("%d", uid)
		_, err = http.Post(util.CreateURL("/douyin/relation/action", query), "", nil)
		if err != nil {
			return
		}
	}
	return
}

func DoFollower(num int, id int64, prefix, action string) (err error) {
	query := map[string]string{
		"to_user_id":  fmt.Sprintf("%d", id),
		"action_type": action,
	}
	for i := 0; i < num; i++ {
		f := prefix + fmt.Sprintf("%d", i)
		_, token, err := util.GetUserIDAndToken(f, f)
		if err != nil {
			return err
		}
		query["token"] = token
		_, err = http.Post(util.CreateURL("/douyin/relation/action", query), "", nil)
		if err != nil {
			return err
		}
	}
	return nil
}
