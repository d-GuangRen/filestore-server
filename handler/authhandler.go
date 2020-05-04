package handler

import (
	"fmt"
	"net/http"
)

// http 请求拦截器
func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		username := r.Form.Get("username")
		token := r.Form.Get("token")

		if len(username) > 16 || !isTokenValid(token) {
			fmt.Println("forbidden")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func isTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}
