package function

import (
	"context"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionE struct {
	BaseFunction
}

func (f *KisFunctionE) Call(ctx context.Context, flow kis.Flow) error {
	log.GetLogger().InfoF("KisFunctionE, flow = %+v\n", flow)

	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
func NewKisFunctionE() kis.Function {
	f := new(KisFunctionE)

	f.metaData = make(map[string]interface{})

	return f

}
