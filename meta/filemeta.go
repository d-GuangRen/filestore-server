package meta

import (
	"filestore-server/db"
	"sort"
)

// FileMeta: 文件元信息结构
type FileMeta struct {
	FileSha1 string // 作为文件唯一标识
	FileName string // 文件名称
	FileSize int64 // 文件大小
	Location string // 文件路径
	UploadAt string // 文件上传时间
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 新增/更新文件元信息
func UpdateFileMeta(fileMeta FileMeta) {
	// fileMetas[fileMeta.FileSha1] = fileMeta
	db.OnFileUploadFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.Location)
}

// 通过sha1获取文件元信息
func GetFileMeta(fileSha1 string) (FileMeta, error) {
	tableFile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	result := FileMeta{
		FileSha1: tableFile.FileHash,
		FileName: tableFile.FileName.String,
		FileSize: tableFile.FileSize.Int64,
		Location: tableFile.FileAddr.String,
		UploadAt: tableFile.CreateAt.Time.String(),
	}
	return result, nil
}

func GetLastFileMetas(count int) []FileMeta {

	// 将 map 转换为切片
	var metas []FileMeta
	for _, v := range fileMetas {
		metas = append(metas, v)
	}

	// 对切片中数据进行排序
	sort.Sort(ByUploadTime(metas))
	// 截取指定长度返回
	return metas[0:count]
}

// 删除文件元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}