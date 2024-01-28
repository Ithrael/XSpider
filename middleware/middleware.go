// middleware.go
package middleware

// Middleware 定义中间件的契约。
type Middleware interface {
	Init() error
	Process(data interface{}) error
	Close() error
}
