package testutil

import (
	"database/sql"
	"net/http"
)

const ROOT = "http://127.0.0.1:30000"

func CreateURL(path string, query map[string]string) (res string) {
	res = ROOT + path
	if len(query) > 0 {
		res += "?"
		for k, v := range query {
			res += k + "=" + v + "&"
		}
		res = res[:len(res)-1]
	}
	return
}

func GetDBConnection() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "longfar:Ning@tcp(127.0.0.1:3306)/douyin_test")
	if err != nil {
		return
	}
	err = db.Ping()
	return
}

func CreateUser(username, password string) (userid int64, token string, err error) {
	query := map[string]string{
		"username": username,
		"password": password,
	}
	resp, err := http.Post(CreateURL("/douyin/user/register", query), "", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respData, err := GetDouyinResponse[DouyinUserRegisterResponse](resp)
	if err != nil {
		return
	}
	return respData.UserID, respData.Token, nil
}

func DeleteUser(username string) (err error) {
	db, err := GetDBConnection()
	if err != nil {
		return
	}
	defer db.Close()

	strDel := "delete from user where username = ?"
	_, err = db.Exec(strDel, username)
	return err
}

func GetUseridAndToken(username, password string) (userid int64, token string, err error) {
	query := map[string]string{
		"username": username,
		"password": password,
	}
	resp, err := http.Post(CreateURL("/douyin/user/register", query), "", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respData, err := GetDouyinResponse[DouyinUserRegisterResponse](resp)
	if err != nil {
		return
	} else if respData.StatusCode == 10111 {
		resp, err = http.Post(CreateURL("/douyin/user/login", query), "", nil)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		var respData DouyinUserLoginResponse
		respData, err = GetDouyinResponse[DouyinUserLoginResponse](resp)
		if err != nil {
			return
		}
		return respData.UserID, respData.Token, nil
	}
	return respData.UserID, respData.Token, nil
}
