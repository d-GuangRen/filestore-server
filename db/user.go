package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

type TableUser struct {
	Username string
	Password string
}
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

func GetByUsername(username string) (*TableUser, error) {
	stmt, err := mydb.DbConn().Prepare("select username, password from tbl_user where username = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result := TableUser{}
	err = stmt.QueryRow(username).Scan(&result.Username, &result.Password)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func UpdateUserToken(username, userToken string) bool {
	stmt, err := mydb.DbConn().Prepare("replace into tbl_user_token (`username`, `user_token`) values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, userToken)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

