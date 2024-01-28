// mysql_middleware.go
package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLMiddleware 是用于将数据存储在MySQL中的中间件。
type MySQLMiddleware struct {
	DB   *sql.DB
	Once sync.Once
}

// 全局变量，用于保存MySQLDB的唯一实例
var mySQLDBInstance *MySQLMiddleware

// GetMySQLDBInstance 返回MySQLDB的唯一实例
func GetMySQLDBInstance() *MySQLMiddleware {
	if mySQLDBInstance == nil {
		mySQLDBInstance = &MySQLMiddleware{}
		// 使用sync.Once确保初始化只执行一次
		mySQLDBInstance.Once.Do(func() {
			log.Println("MYSQL_URL: ", os.Getenv("MYSQL_URL"))
			// 在这里初始化MySQL连接池
			db, err := sql.Open("mysql", os.Getenv("MYSQL_URL"))
			if err != nil {
				panic(err.Error())
			}
			mySQLDBInstance.DB = db
		})
	}
	return mySQLDBInstance
}

// Init 初始化 MySQLMiddleware。
func (m *MySQLMiddleware) Init() error {
	return nil
}

// Process 处理数据并将其存储在MySQL中。
func (m *MySQLMiddleware) Process(data interface{}) error {
	detailsCh, ok := data.(chan *PageDetail)
	if !ok {
		return errors.New("failed to convert data to PageDetail type")
	}

	for detail := range detailsCh {
		_, err := mySQLDBInstance.DB.Exec("INSERT INTO websites (url, title, host, code, finger, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			detail.Url, detail.Title, detail.Host, detail.ResponseCode, detail.Fingerprint, detail.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to insert data into MySQL: %v", err)
		}
	}
	return nil
}

// Close 关闭MySQLMiddleware使用的任何资源。
func (m *MySQLMiddleware) Close() error {
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			return fmt.Errorf("failed to close MySQL connection: %v", err)
		}
	}
	return nil
}
