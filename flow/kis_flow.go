package flow

import (
	"context"
	"errors"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/function"
	"kis-flow/id"
	"kis-flow/kis"
	"sync"
)

type KisFlow struct {
	Id   string
	Name string

	Conf *config.KisFlowConfig // Flow配置策略

	// Function列表
	Funcs          map[string]kis.Function // 当前flow拥有的全部管理的全部Function对象, key: FunctionID
	FlowHead       kis.Function            // 当前Flow所拥有的Function列表表头
	FlowTail       kis.Function            // 当前Flow所拥有的Function列表表尾
	flock          sync.RWMutex            // 管理链表插入读写的锁
	ThisFunction   kis.Function            // Flow当前正在执行的KisFunction对象
	ThisFunctionId string                  // 当前执行到的Function ID (策略配置ID)
	PrevFunctionId string                  // 当前执行到的Function 上一层FunctionID(策略配置ID)

	// Function列表参数
	funcParams map[string]config.FParam // flow在当前Function的自定义固定配置参数,Key:function的实例NsID, value:FParam
	fplock     sync.RWMutex             // 管理funcParams的读写锁
}

func (flow *KisFlow) Run(ctx context.Context) error {
	var fn kis.Function

	fn = flow.FlowHead

	if flow.Conf.Status == int(common.FlowDisable) {
		//flow被配置关闭
		return nil
	}

	for fn != nil {
		if err := fn.Call(ctx, flow); err != nil {
			return err
		} else {
			fn = fn.Next()
		}
	}

	return nil
}

func (flow *KisFlow) Link(fConfig *config.KisFuncConfig, fParams config.FParam) error {
	f := function.NewKisFunction(flow, fConfig)

	if err := flow.appendFunc(f, fParams); err != nil {
		return err
	}
	return nil
}

func NewKisFlow(conf *config.KisFlowConfig) *KisFlow {
	flow := new(KisFlow)

	flow.Id = id.KisId(common.KisIdTypeFlow)
	flow.Name = conf.FlowName
	flow.Conf = conf

	// Function列表
	flow.Funcs = make(map[string]kis.Function)
	flow.funcParams = make(map[string]config.FParam)
	return flow
}

func (flow *KisFlow) appendFunc(function kis.Function, fParam config.FParam) error {
	if function == nil {
		return errors.New("AppendFunc append nil to List")
	}

	flow.flock.Lock()
	defer flow.flock.Unlock()

	if flow.FlowHead == nil {
		flow.FlowHead = function
		flow.FlowTail = function

		function.SetP(nil)
		function.SetN(nil)
	} else {
		function.SetP(flow.FlowTail)
		function.SetN(nil)
		flow.FlowTail.SetN(function)
		flow.FlowTail = function
	}
	flow.Funcs[function.GetId()] = function

	params := make(config.FParam)
	for k, v := range function.GetConfig().Option.Params {
		params[k] = v
	}

	for k, v := range fParam {
		params[k] = v
	}

	flow.funcParams[function.GetId()] = params

	return nil
}
