package main

import (
	util "douyin/test/testutil"
	"errors"
	"fmt"
	"net/http"
)

func main() {
	id, token, err := util.GetUseridAndToken("jmeter-comment", "123456")
	if err != nil {
		panic(err)
	}
	fmt.Printf("id: %v", id)
	_ = token
	err = AddComments(token, "1", "This is a comment for comment list test.", 100)
	if err != nil {
		panic(err)
	}
}

func AddComments(token, video_id, content string, num int) error {
	q := map[string]string{
		"token":        token,
		"video_id":     video_id,
		"action_type":  "1",
		"comment_text": content,
	}
	for i := 0; i < num; i++ {
		resp, err := http.Post(util.CreateURL("/douyin/comment/action", q), "", nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return errors.New("return code != 200")
		}
		fmt.Println("add a comment")
	}
	return nil
}
