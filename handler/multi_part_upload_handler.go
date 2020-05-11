package handler

import (
	myCache "filestore-server/cache/reids"
	"filestore-server/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type MultiPartUploadInfo struct {
	UploadId string
	FileHash string
	FileSize int
	ChunkSize int
	ChunkCount int
}

// 初始化分块上传
func InitialMultiPartUploadHandler(w http.ResponseWriter, r *http.Request) {
	 r.ParseForm()

	username := r.Form.Get("username")
	fileHash := r.Form.Get("fileHash")
	fileSize, err := strconv.Atoi(r.Form.Get("fileSize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Params Invalid", nil).JSONBytes())
		return
	}

	// 获取redis连接
	conn := myCache.RedisPool().Get()
	defer conn.Close()

	uploadInfo := MultiPartUploadInfo{
		FileHash:   fileHash,
		FileSize:   fileSize,
		UploadId:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))),
	}

	conn.Do("HSET", "MP_" + uploadInfo.UploadId, "chunkCount", uploadInfo.ChunkCount)
	conn.Do("HSET", "MP_" + uploadInfo.UploadId, "fileHash", uploadInfo.FileHash)
	conn.Do("HSET", "MP_" + uploadInfo.UploadId, "fileSize", uploadInfo.FileSize)

	w.Write(util.NewRespMsg(0, "OK", uploadInfo).JSONBytes())
}

// 分块上传
func MultiPartUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	uploadId := r.Form.Get("uploadId")
	chunkIndex := r.Form.Get("chunkIndex")

	conn := myCache.RedisPool().Get()
	defer conn.Close()

	fileData, err := os.Create("/data/" + uploadId + "/" + chunkIndex)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fileData.Close()

	buffer := make([]byte, 1024 * 1024)
	for {
		n, err := r.Body.Read(buffer)
		fileData.Write(buffer[:n])
		if err != nil {
			break
		}
	}

	conn.Do("HSET", "MP_" + uploadId, "chunkIndex_" + chunkIndex, 1)

	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}