package function

import (
	"context"
	"fmt"
	"kis-flow/kis"
)

type KisFunctionV struct {
	BaseFunction
}

func (f *KisFunctionV) Call(ctx context.Context, flow kis.Flow) error {
	fmt.Printf("KisFunctionV, flow = %+v\n", flow)

	//TODO 调用具体的 Function 执行方法

	return nil
}
