package handler

import (
	"encoding/json"
	"filestore-server/db"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 返回上传的html页面
		fileData, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal server error")
			return
		}
		io.WriteString(w, string(fileData))

	} else if r.Method == http.MethodPost {
		// 接收上传文件流并存储到本地目录
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Faild to get data, err: %s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileSha1: "",
			FileName: header.Filename,
			FileSize: 0,
			Location: "/tmp/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create("/tmp/" + header.Filename)
		if err != nil {
			fmt.Printf("Faild to create file, err: %s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Faild to save data to file, err: %s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		r.ParseForm()
		username := r.Form.Get("username")
		status := db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if !status {
			w.Write([]byte("Upload Failed."))
			return
		}

		http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	}
}

// 文件上传成功
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// 文件元信息查询
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form["fileHash"][0]
	fileMeta := meta.GetFileMeta(fileHash)

	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(marshal)
}

func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	count, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	userFiles, err := db.QueryUserFileMetas(username, count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshal, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
}

// 文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileSha1 := r.Form.Get("fileHash")
	fileMeta := meta.GetFileMeta(fileSha1)

	file, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\"" + fileMeta.FileName + "\"")
	w.Write(fileData)
}

// 文件元信息更新
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	option := r.Form.Get("option")
	if option != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fileSha1 := r.Form.Get("fileHash")
	filename := r.Form.Get("filename")

	fileMeta := meta.GetFileMeta(fileSha1)

	fileMeta.FileName = filename
	meta.UpdateFileMeta(fileMeta)

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
}

// 文件删除
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form.Get("fileHash")

	fileMeta := meta.GetFileMeta(fileHash)
	os.Remove(fileMeta.Location)
	meta.RemoveFileMeta(fileHash)

	w.WriteHeader(http.StatusOK)
}

// 尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	//// 解析请求参数
	//username := r.Form.Get("username")
	//fileHash := r.Form.Get("fileHash")
	//fileName := r.Form.Get("fileName")
	//fileSize := r.Form.Get("fileSize")
	//
	//// 从文件表中查询相同hash的文件记录
	//fileMeta, err := meta.GetFileMeta(fileHash)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//if fileMeta == nil {
	//
	//}
}













