// csv_middleware.go
package middleware

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sync"
)

// CSVMiddleware 实现 Middleware 接口的 CSV 中间件。
type CSVMiddleware struct {
	file   *os.File
	writer *csv.Writer
	Once   sync.Once
}

// 全局变量，用于保存csv的唯一实例
var csvInstance *CSVMiddleware

// GetCsvInstance 返回csv的唯一实例
func GetCsvInstance() *CSVMiddleware {
	if csvInstance == nil {
		csvInstance = &CSVMiddleware{}
		// 使用sync.Once确保初始化只执行一次
		csvInstance.Once.Do(func() {
			// 在这里初始化csv writer
			f, err := os.Create("output.csv")
			if err != nil {
				panic(err.Error())
			}
			csvInstance.file = f
			csvInstance.writer = csv.NewWriter(f)
		})
	}
	return csvInstance
}

func (c *CSVMiddleware) Init() error {
	return nil
}

// Process 处理数据并将其存储在 CSV 中。
func (c *CSVMiddleware) Process(data interface{}) error {
	detailsCh, ok := data.(chan *PageDetail)
	if !ok {
		return errors.New("failed to convert data to PageDetail type")
	}

	for detail := range detailsCh {
		record := []string{
			detail.Url,
			detail.Title,
			detail.Host,
			fmt.Sprint(detail.ResponseCode),
			detail.Fingerprint,
			detail.Timestamp,
		}
		if err := c.writer.Write(record); err != nil {
			return fmt.Errorf("failed to write data to CSV: %v", err)
		}
		c.writer.Flush()
	}
	return nil
}

// Close 关闭 CSV 中间件使用的任何资源。
func (c *CSVMiddleware) Close() error {
	if err := c.writer.Error(); err != nil {
		return fmt.Errorf("failed to flush CSV writer: %v", err)
	}

	if err := c.file.Close(); err != nil {
		return fmt.Errorf("failed to close CSV file: %v", err)
	}
	return nil
}
