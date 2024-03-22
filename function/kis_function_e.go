package function

import (
	"context"
	"fmt"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionE struct {
	BaseFunction
}

func (f *KisFunctionE) Call(ctx context.Context, flow kis.Flow) error {
	fmt.Printf("KisFunctionE, flow = %+v\n", flow)

	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.GetLogger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
