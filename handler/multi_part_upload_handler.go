package handler

import (
	myCache "filestore-server/cache/redis"
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
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

	filePath := "/data/" + uploadId + "/" + chunkIndex
	os.MkdirAll(path.Dir(filePath), 0744)
	fileData, err := os.Create(filePath)
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

// 通知上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	uploadId := r.Form.Get("uploadId")
	fileSize := r.Form.Get("fileSize")
	fileHash := r.Form.Get("fileHash")
	fileName := r.Form.Get("fileName")
	username := r.Form.Get("username")

	conn := myCache.RedisPool().Get()
	defer conn.Close()

	cacheData, err := redis.Values(conn.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
		return
	}

	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(cacheData); i += 2 {
		key := string(cacheData[i].([]byte))
		value := string(cacheData[i+1].([]byte))
		if key == "chunkCount" {
			totalCount, _ = strconv.Atoi(key)
		} else if strings.HasPrefix(key, "chunkIndex_") && value == "1" {
			chunkCount++
		}
	}

	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "params invalid", nil).JSONBytes())
		return
	}

	// TODO 文件块合并

	fSize, _ := strconv.Atoi(fileSize)
	db.OnFileUploadFinished(fileHash, fileName, int64(fSize), "")
	db.OnUserFileUploadFinished(username, fileHash, fileName, int64(fSize))

	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}