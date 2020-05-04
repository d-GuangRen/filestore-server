package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"io/ioutil"
	"net/http"
)

const passwordSalt = "xQ619G%S"

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


