package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	passwordSalt = "xQ619G%S"
	tokenSalt = "_tokensalt"
)

// 用户注册
func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		fileData, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(fileData)
		return
	}

	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) > 16 || len(password) > 16 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	// 对密码进行加密
	encodePassword := util.Sha1([]byte(password + passwordSalt))
	status := db.SignUp(username, encodePassword)
	if status {
		w.Write([]byte("SUCCESS"))
		return
	}
	w.Write([]byte("FAILED"))
}

func SignInHandler (w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	tableUser, err := db.GetByUsername(username)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("FAILED"))
		return
	}
	usernameDb := tableUser.Username
	passwordDb := tableUser.Password

	if passwordDb != util.Sha1([]byte(password + passwordSalt)) {
		w.Write([]byte("FAILED"))
		return
	}
	token := generateToken(usernameDb)
	updateStatus := db.UpdateUserToken(usernameDb, token)
	if !updateStatus {
		w.Write([]byte("FAILED"))
		return
	}

	// w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	respMsg := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: usernameDb,
			Token:    token,
		},
	}
	w.Write(respMsg.JSONBytes())
}

func generateToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + tokenSalt))
	return tokenPrefix + ts[:8]
}
