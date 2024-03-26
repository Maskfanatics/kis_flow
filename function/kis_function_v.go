package function

import (
	"context"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionV struct {
	BaseFunction
}

func (f *KisFunctionV) Call(ctx context.Context, flow kis.Flow) error {
	log.GetLogger().InfoF("KisFunctionV, flow = %+v\n", flow)

	//TODO 调用具体的 Function 执行方法
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}
	return nil
}

func NewKisFunctionV() kis.Function {
	f := new(KisFunctionV)

	f.metaData = make(map[string]interface{})

	return f

}
