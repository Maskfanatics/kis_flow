package function

import (
	"context"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionS struct {
	BaseFunction
}

func (f *KisFunctionS) Call(ctx context.Context, flow kis.Flow) error {
	log.GetLogger().InfoF("KisFunctionS, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}
	return nil
}
