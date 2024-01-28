// csv_middleware_test.go
package middleware

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVMiddleware_Process(t *testing.T) {
	// 创建一个临时文件用于测试
	tempFile, err := os.CreateTemp("", "test.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// 创建 CSVMiddleware 实例
	csvMiddleware := &CSVMiddleware{}
	csvMiddleware.file = tempFile
	csvMiddleware.writer = csv.NewWriter(tempFile)

	// 创建一个用于传递数据的 channel
	detailsCh := make(chan *PageDetail)

	// 启动一个 goroutine 来运行 Process
	go func() {
		// 将测试数据发送到 channel
		detailsCh <- &PageDetail{
			Url:          "https://example.com",
			Title:        "Example Title",
			Host:         "example.com",
			ResponseCode: 200,
			Fingerprint:  "fingerprint123",
			Timestamp:    "2023-01-12 17:12:16",
		}

		// 关闭 channel，表示数据发送完毕
		close(detailsCh)
	}()

	// 运行 Process 函数
	err = csvMiddleware.Process(detailsCh)
	assert.NoError(t, err)

	// 关闭 CSVMiddleware，触发关闭操作
	err = csvMiddleware.Close()
	assert.NoError(t, err)

	// 重新打开临时文件以验证写入的内容
	contents, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	// 验证文件内容是否符合预期
	expectedContents := "https://example.com,Example Title,example.com,200,fingerprint123,2023-01-12 17:12:16\n"
	fmt.Println(string(contents))
	fmt.Println(expectedContents)
	assert.Equal(t, expectedContents, string(contents))
}
