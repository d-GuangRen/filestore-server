package meta

import "sort"

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
	fileMetas[fileMeta.FileSha1] = fileMeta
}

// 通过sha1获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
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
