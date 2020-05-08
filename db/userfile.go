package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	Username   string
	FileHash   string
	FileName   string
	FileSize   int64
	UploadAt   string
	LastUpdated string
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

// 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DbConn().Prepare("select file_sha1, file_name, file_size, upload_at, last_update from tbl_user_file where username = ? limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}
	var userFiles []UserFile
	for rows.Next() {
		userFile := UserFile{}
		err := rows.Scan(&userFile.FileHash, &userFile.FileName, &userFile.FileSize, &userFile.UploadAt, &userFile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, userFile)
	}
	return userFiles, nil
}