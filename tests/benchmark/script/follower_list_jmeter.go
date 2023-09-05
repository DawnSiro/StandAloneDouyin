package main

import (
	util "douyin/test/testutil"
	"fmt"
	"net/http"
)

func main() {
	id, _, err := util.GetUseridAndToken("jmeter-follower-list", "123456")
	if err != nil {
		panic(err)
	}
	err = DoFollower(100, id, "jmeter-fan", "1")
	if err != nil {
		panic(err)
	}
}

func DoFollower(num int, id int64, prefix, action string) (err error) {
	query := map[string]string{
		"to_user_id":  fmt.Sprintf("%d", id),
		"action_type": action,
	}
	for i := 0; i < num; i++ {
		f := prefix + fmt.Sprintf("%d", i)
		_, token, err := util.GetUseridAndToken(f, f)
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
