package log

import "context"
import "fmt"

type kisDefaultLog struct{}

func (log *kisDefaultLog) InfoF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
	fmt.Println()
}

func (log *kisDefaultLog) ErrorF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
	fmt.Println()
}

func (log *kisDefaultLog) DebugF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
	fmt.Println()
}

func (log *kisDefaultLog) InfoFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
	fmt.Println()
}

func (log *kisDefaultLog) ErrorFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
	fmt.Println()
}

func (log *kisDefaultLog) DebugFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
	fmt.Println()
}

func init() {
	if GetLogger() == nil {
		setLogger(&kisDefaultLog{})
	}
}
