package function

import (
	"context"
	"errors"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/id"
	"kis-flow/kis"
	"sync"
)

type BaseFunction struct {
	Id     string
	Config *config.KisFuncConfig

	Flow kis.Flow

	connector kis.Connector

	metaData map[string]interface{}
	mLock    sync.RWMutex

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
		f = NewKisFunctionV()
	case common.S:
		f = NewKisFunctionS()
	case common.L:
		f = NewKisFunctionL()
	case common.C:
		f = NewKisFunctionC()
	case common.E:
		f = NewKisFunctionE()
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

func (base *BaseFunction) AddConnector(conn kis.Connector) error {
	if conn == nil {
		return errors.New("conn is nil")
	}
	base.connector = conn

	return nil
}

func (base *BaseFunction) GetConnector() kis.Connector {
	return base.connector
}

func (base *BaseFunction) GetMetaData(key string) interface{} {
	base.mLock.RLock()
	defer base.mLock.RUnlock()

	data, ok := base.metaData[key]
	if !ok {
		return nil
	}

	return data
}

func (base *BaseFunction) SetMetaData(key string, value interface{}) {
	base.mLock.Lock()
	defer base.mLock.Unlock()

	base.metaData[key] = value

}
