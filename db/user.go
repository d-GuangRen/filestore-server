package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

// 用户注册
func SignUp(username, password string) bool {
	stmt, err := mydb.DbConn().Prepare("insert ignore into tbl_user(`username`, `password`) values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	execResult, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rowsAffected, err := execResult.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}

	return false
}

func SignIn(username, password string) {

}
