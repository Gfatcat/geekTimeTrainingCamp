package ch02_error

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"io"
	"os"
)

var db *sql.DB

type Employee struct {
	id   int
	name string
	age  int
}

type ErrLogger struct {
	io.Writer
	err error
}

func (e *ErrLogger) Write(buf []byte) (int, error) {
	if e.err != nil {
		return 0, nil
	}
	var n int
	n, e.err = e.Writer.Write(buf)
	return n, e.err
}

type QueryError struct {
	Err   error
	Query string
}

func (e *QueryError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Err.Error()
	if e.Query != "" {
		s = "query sql  " + e.Query + ": " + s
	}
	return s
}

func InitDbConn() error {
	var getMysqlDsn = func(host, port, pwd, user, database string) string {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&allowOldPasswords=True", user, pwd, host, port, database)
	}
	chatFilterDsn := getMysqlDsn("127.0.0.1", "3306", "123456", "root", "dev")
	var err error
	db, err = sql.Open("mysql", chatFilterDsn)
	if err != nil {
		return errors.Wrap(err, "fail to connect database")
	}
	return nil
}

func GetEmployee(id int) (*string, error) {
	querySQL := fmt.Sprintf("select name from employees where id = %d", id)
	var name string
	err := db.QueryRow(querySQL).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, errors.Wrap(&QueryError{err, querySQL}, "fail to Query employee")
		}
	}
	return &name, nil
}

func Run(id int) {
	err := InitDbConn()
	if err != nil {
		fmt.Printf("original error %T %v \n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace:\n%+v\n", err)
		os.Exit(1)
	}
	name, err := GetEmployee(id)
	if err != nil {
		fmt.Printf("original error %T %v \n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace:\n%+v\n", err)
		os.Exit(1)
	}
	if name == nil {
		fmt.Println("no employee found")
	} else {
		fmt.Println(*name)
	}
	return
}
