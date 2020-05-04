package db

import (
	mydb "filestore-server/db/mysql"
	"time"
)

type UserFile struct {
	Username string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdate string
}

func OnUserFileUploadFinished(username, fileHash, fileName string, fileSize int64) bool {
	stmt, err := mydb.DbConn().Prepare("insert ignore into `tbl_user_file` (`username`, `file_sha1`, `file_name`, `file_size`, `upload_at`, `status`) " +
		"values (?,?,?,?,?,1)")
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, fileHash, fileName, fileSize, time.Now())
	if err != nil {
		return false
	}

	return true
}
