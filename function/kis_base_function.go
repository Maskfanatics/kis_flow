package function

import (
	"context"
	"errors"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/id"
	"kis-flow/kis"
)

type BaseFunction struct {
	Id     string
	Config *config.KisFuncConfig

	Flow kis.Flow

	N kis.Function
	P kis.Function
}

func (base *BaseFunction) Call(ctx context.Context, flow kis.Flow) error {
	return nil
}

func (base *BaseFunction) Next() kis.Function {
	return base.N
}

func (base *BaseFunction) Prev() kis.Function {
	return base.P
}

func (base *BaseFunction) SetN(f kis.Function) {
	base.N = f
}

func (base *BaseFunction) SetP(f kis.Function) {
	base.P = f
}

func (base *BaseFunction) SetConfig(s *config.KisFuncConfig) error {
	if s == nil {
		return errors.New("KisFuncConfig is nil")
	}

	base.Config = s
	return nil
}

func (base *BaseFunction) GetId() string {
	return base.Id
}

func (base *BaseFunction) GetPrevId() string {
	if base.P == nil {
		return common.FunctionIdFirstVirtual
	}
	return base.P.GetId()
}

func (base *BaseFunction) GetNextId() string {
	if base.N == nil {
		return common.FunctionIdLastVirtual
	}
	return base.N.GetId()
}

func (base *BaseFunction) GetConfig() *config.KisFuncConfig {
	return base.Config
}

func (base *BaseFunction) SetFlow(f kis.Flow) error {
	if f == nil {
		return errors.New("KisFlow is nil")
	}

	base.Flow = f
	return nil
}

func (base *BaseFunction) GetFlow() kis.Flow {
	return base.Flow
}

func (base *BaseFunction) CreateId() {
	base.Id = id.KisId(common.KisIdTypeFunction)
}

// NewKisFunction 创建一个NsFunction
// flow: 当前所属的flow实例
// s : 当前function的配置策略

func NewKisFunction(flow kis.Flow, config *config.KisFuncConfig) kis.Function {
	var f kis.Function

	switch common.KisMode(config.FMode) {
	case common.V:
		f = new(KisFunctionV)
		break
	case common.S:
		f = new(KisFunctionS)
	case common.L:
		f = new(KisFunctionL)
	case common.C:
		f = new(KisFunctionC)
	case common.E:
		f = new(KisFunctionE)
	default:
		//LOG ERROR
		return nil
	}
	f.CreateId()

	if err := f.SetConfig(config); err != nil {
		panic(err)
	}

	if err := f.SetFlow(flow); err != nil {
		panic(err)
	}
	return f
}
