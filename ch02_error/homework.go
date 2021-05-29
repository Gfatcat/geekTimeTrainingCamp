package ch02_error

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"time"
)

var db *sql.DB

type Employee struct {
	id   int
	name string
	age  int
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
		s = s + "\n" + "query sql: " + e.Query
	}
	return s
}

type Logger struct {
	fileName string      // 日志文件名
	file     *os.File    // 日志文件
	Info     *log.Logger // 日志
	Warning  *log.Logger // 警告
	Error    *log.Logger // 错误
}

func (logger *Logger) New() *Logger {
	logger.fileName = time.Now().Format("20060102") + ".log"
	var err error
	logger.file, err = os.OpenFile(logger.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open error log file:", err)
	}
	mw := io.MultiWriter(os.Stdout, logger.file)
	logger.Info = log.New(mw, "[Info]", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Warning = log.New(mw, "[Warning]", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Error = log.New(mw, "[Error]", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

func (logger *Logger) Close() {
	err := logger.file.Close()
	if err != nil {
		log.Println("Failed to close error log file:", err)
	}
	return
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
	err = db.Ping()
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
			return nil, errors.Wrap(&QueryError{err, querySQL}, "fail to query employee")
		}
	}
	return &name, nil
}

func RunHomework(id int) {
	// 初始化
	Logger := Logger{"", nil, nil, nil, nil}
	Logger.New()

	err := InitDbConn()
	if err != nil {
		//ErrLogger.Printf("original error %T %v \n", errors.Cause(err), errors.Cause(err))
		//log.Printf("original error: %T\nstack trace:\n%+v", errors.Cause(err),err)
		Logger.Error.Printf("original error: %T\nstack trace:\n%+v", errors.Cause(err), err)
		os.Exit(1)
	}
	name, err := GetEmployee(id)
	if err != nil {
		//ErrLogger.Printf("original error %T %v \n", errors.Cause(err), errors.Cause(err))
		//log.Printf("stack trace:\n%+v\n", err)
		Logger.Error.Printf("stack trace:\n%+v\n", err)
		os.Exit(1)
	}
	if name == nil {
		//log.Println("no employee found")
		Logger.Error.Println("no employee found")
	} else {
		fmt.Print(&name)
	}
	Logger.Close()
	return
}
