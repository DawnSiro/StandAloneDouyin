// Code generated by hertz generator. DO NOT EDIT.

package Api

import (
	api "douyin/biz/handler/api"
	"github.com/cloudwego/hertz/pkg/app/server"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_douyin := root.Group("/douyin", _douyinMw()...)
		{
			_comment := _douyin.Group("/comment", _commentMw()...)
			{
				_action := _comment.Group("/action", _actionMw()...)
				_action.POST("/", append(_comment_ctionMw(), api.CommentAction)...)
			}
			{
				_list := _comment.Group("/list", _listMw()...)
				_list.GET("/", append(_getcommentlistMw(), api.GetCommentList)...)
			}
		}
		{
			_favorite := _douyin.Group("/favorite", _favoriteMw()...)
			{
				_action0 := _favorite.Group("/action", _action0Mw()...)
				_action0.POST("/", append(_favoritevideoMw(), api.FavoriteVideo)...)
			}
			{
				_list0 := _favorite.Group("/list", _list0Mw()...)
				_list0.GET("/", append(_getfavoritelistMw(), api.GetFavoriteList)...)
			}
		}
		{
			_feed := _douyin.Group("/feed", _feedMw()...)
			_feed.GET("/", append(_getfeedMw(), api.GetFeed)...)
		}
		{
			_message := _douyin.Group("/message", _messageMw()...)
			{
				_action1 := _message.Group("/action", _action1Mw()...)
				_action1.POST("/", append(_sendmessageMw(), api.SendMessage)...)
			}
			{
				_chat := _message.Group("/chat", _chatMw()...)
				_chat.GET("/", append(_getmessagechatMw(), api.GetMessageChat)...)
			}
		}
		{
			_relation := _douyin.Group("/relation", _relationMw()...)
			{
				_action2 := _relation.Group("/action", _action2Mw()...)
				_action2.POST("/", append(_followMw(), api.Follow)...)
			}
			{
				_follow0 := _relation.Group("/follow", _follow0Mw()...)
				{
					_list1 := _follow0.Group("/list", _list1Mw()...)
					_list1.GET("/", append(_getfollowlistMw(), api.GetFollowList)...)
				}
			}
			{
				_follower := _relation.Group("/follower", _followerMw()...)
				{
					_list2 := _follower.Group("/list", _list2Mw()...)
					_list2.GET("/", append(_getfollowerlistMw(), api.GetFollowerList)...)
				}
			}
			{
				_friend := _relation.Group("/friend", _friendMw()...)
				{
					_list3 := _friend.Group("/list", _list3Mw()...)
					_list3.GET("/", append(_getfriendlistMw(), api.GetFriendList)...)
				}
			}
		}
		{
			_user := _douyin.Group("/user", _userMw()...)
			_user.GET("/", append(_getuserinfoMw(), api.GetUserInfo)...)
			{
				_login := _user.Group("/login", _loginMw()...)
				_login.POST("/", append(_login0Mw(), api.Login)...)
			}
			{
				_register := _user.Group("/register", _registerMw()...)
				_register.POST("/", append(_register0Mw(), api.Register)...)
			}
		}
	}
}
