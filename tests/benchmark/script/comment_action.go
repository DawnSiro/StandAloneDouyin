package main

import (
	"fmt"

	util "douyin/test/testutil"
)

const (
	CommentPrefix = "benchmark-commenter"
	NumComment    = 1000
)

func main() {
	for i := 0; i < NumFans; i++ {
		n := fmt.Sprintf("%s-%d", CommentPrefix, i)
		_, _, err := util.GetUserIDAndToken(n, password)
		assert(err)
	}
	fmt.Println("Create fans ok")
}
