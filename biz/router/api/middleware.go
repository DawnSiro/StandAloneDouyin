// Code generated by hertz generator.

package Api

import (
	"douyin/biz/mw"
	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _douyinMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _commentMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _actionMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _comment_ctionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _listMw() []app.HandlerFunc {
	// your code...
	// 未登录可以查看评论
	return nil
}

func _getcommentlistMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _favoriteMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _list0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getfavoritelistMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _feedMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getfeedMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _messageMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _action0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _sendmessageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _chatMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getmessagechatMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _publishMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _action1Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _publish_ctionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list1Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getpublishvideosMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _relationMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _action2Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _followMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _follow0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list2Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getfollowlistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _followerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list3Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getfollowerlistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _friendMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list4Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getfriendlistMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getuserinfoMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _login0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _register0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _favoritevideoMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		// use jwt mw
		mw.JwtMiddleware.MiddlewareFunc(),
	}
}

func _action3Mw() []app.HandlerFunc {
	// your code...
	return nil
}
