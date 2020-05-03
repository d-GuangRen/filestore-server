package db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
} 

func OnFileUploadFinished(fileSha1 string, filename string, fileSize int64, fileAddr string) bool {
	conn := mydb.DbConn()
	defer conn.Close()
	
	stmt, err := conn.Prepare("insert ignore into tbl_file (`file_sha1`, `file_name`, `file_size`, `file_addr`) values (?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("Failed to prepare statement, err: %s", err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(fileSha1, filename, fileSize, fileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := result.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash: %s has been uploaded before ", fileSha1)
		}
		return true
	}
	return false
}

func GetFileMeta(fileHash string) (*TableFile, error) {
	conn := mydb.DbConn()
	defer conn.Close()

	stmt, err := conn.Prepare("select file_sha1, file_addr, file_name, file_size from tbl_file where file_sha1 = ? and status = 1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	result := TableFile{}
	err = stmt.QueryRow(fileHash).Scan(&result.FileHash, &result.FileAddr, &result.FileName, &result.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &result, nil
} 
