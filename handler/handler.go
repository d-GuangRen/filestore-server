package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// 文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传的html页面
		fileData, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal server error")
			return
		}
		io.WriteString(w, string(fileData))
	} else if r.Method == "POST" {
		// 接收上传文件流并存储到本地目录
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Faild to get data, err: %s\n", err.Error())
			return
		}
		defer file.Close()

		newFile, err := os.Create("/tmp/" + header.Filename)
		if err != nil {
			fmt.Printf("Faild to create file, err: %s\n", err.Error())
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Faild to save data to file, err: %s\n", err.Error())
			return
		}

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

// 文件上传成功
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}
