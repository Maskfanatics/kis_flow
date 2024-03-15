package log

import "context"

type KisLogger interface {
	InfoFX(ctx context.Context, str string, v ...interface{})
	ErrorFX(ctx context.Context, str string, v ...interface{})
	DebugFX(ctx context.Context, str string, v ...interface{})

	InfoF(str string, v ...interface{})
	ErrorF(str string, v ...interface{})
	DebugF(str string, v ...interface{})
}

// 定义一个全局的默认kisLog 对象
var kisLog KisLogger

// Logger 设置 KisLog对象，可以是用户自定义的Logger 对象
func setLogger(newloger KisLogger) {
	kisLog = newloger
}

// Logger 获取到 kisLog 对象
func GetLogger() KisLogger {
	return kisLog
}
