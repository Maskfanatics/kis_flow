package flow

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/conn"
	"kis-flow/function"
	"kis-flow/id"
	"kis-flow/kis"
	"kis-flow/log"
	"kis-flow/metrics"
	"sync"
	"time"
)

type KisFlow struct {
	Id   string
	Name string

	Conf *config.KisFlowConfig // Flow配置策略

	// Function列表
	Funcs          map[string]kis.Function // 当前flow拥有的全部管理的全部Function对象, key: FunctionName
	FlowHead       kis.Function            // 当前Flow所拥有的Function列表表头
	FlowTail       kis.Function            // 当前Flow所拥有的Function列表表尾
	flock          sync.RWMutex            // 管理链表插入读写的锁
	ThisFunction   kis.Function            // Flow当前正在执行的KisFunction对象
	ThisFunctionId string                  // 当前执行到的Function ID (策略配置ID)
	PrevFunctionId string                  // 当前执行到的Function 上一层FunctionID(策略配置ID)

	// Function列表参数
	funcParams map[string]config.FParam // flow在当前Function的自定义固定配置参数,Key:function的实例NsID, value:FParam
	fplock     sync.RWMutex             // 管理funcParams的读写锁

	// ++++++++ 数据 ++++++++++
	buffer common.KisRowArr  // 用来临时存放输入字节数据的内部Buf, 一条数据为interface{}, 多条数据为[]interface{} 也就是KisBatch
	data   common.KisDataMap // 流式计算各个层级的数据源
	inPut  common.KisRowArr  // 当前Function的计算输入数据

	// ++++++++ KisFlow Action++++++++++

	action kis.Action

	abort bool

	cache *cache.Cache

	metaData map[string]interface{}
	mLock    sync.RWMutex
}

func (flow *KisFlow) Run(ctx context.Context) error {

	var fn kis.Function

	fn = flow.FlowHead

	flow.abort = false

	if flow.Conf.Status == int(common.FlowDisable) {
		//flow被配置关闭
		return nil
	}

	var funcStart time.Time
	var flowStart time.Time

	flow.PrevFunctionId = common.FunctionIdFirstVirtual

	if err := flow.commitSrcData(ctx); err != nil {
		return err
	}

	if config.GlobalConfig.EnableProm {
		metrics.Metrics.FlowScheduleCntsTotal.WithLabelValues(flow.Name).Inc()
	}

	for fn != nil && flow.abort == false {

		// flow记录当前执行到的Function 标记
		fid := fn.GetId()
		flow.ThisFunction = fn
		flow.ThisFunctionId = fid

		fName := fn.GetConfig().FName
		fMode := fn.GetConfig().FMode

		if config.GlobalConfig.EnableProm {
			metrics.Metrics.FuncScheduleCntsTotal.WithLabelValues(fName, fMode).Inc()

			funcStart = time.Now()
		}

		// 得到当前Function要处理与的源数据
		if inputData, err := flow.getCurData(); err != nil {
			log.GetLogger().ErrorFX(ctx, "flow.Run(): getCurData err = %s\n", err.Error())
			return err
		} else {
			flow.inPut = inputData
		}

		if err := fn.Call(ctx, flow); err != nil {
			return err
		} else {
			fn, err = flow.dealAction(ctx, fn)
			if err != nil {
				return err
			}
		}

		if config.GlobalConfig.EnableProm {
			duration := time.Since(funcStart)

			metrics.Metrics.FunctionDuration.With(
				prometheus.Labels{
					common.LABEL_FUNCTION_NAME: fName,
					common.LABEL_FUNCTION_MODE: fMode,
				}).Observe(duration.Seconds() * 1000)
		}

	}

	if config.GlobalConfig.EnableProm {
		duration := time.Since(flowStart)

		metrics.Metrics.FlowDuration.WithLabelValues(flow.Name).Observe(duration.Seconds() * 1000)
	}

	return nil
}

func (flow *KisFlow) Link(fConfig *config.KisFuncConfig, fParams config.FParam) error {
	f := function.NewKisFunction(flow, fConfig)

	if fConfig.Option.CName != "" {
		connConfig, err := fConfig.GetConnConfig()
		if err != nil {
			panic(err)
		}
		connector := conn.NetKisConnector(connConfig)

		if err = connector.Init(); err != nil {
			panic(err)
		}

		_ = f.AddConnector(connector)
	}

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

	//数据
	flow.data = make(common.KisDataMap)

	flow.cache = cache.New(cache.NoExpiration, common.DefaultFlowCacheCleanUp*time.Minute)

	flow.metaData = make(map[string]interface{})

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
	flow.Funcs[function.GetConfig().FName] = function

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

func (flow *KisFlow) Input() common.KisRowArr {
	return flow.inPut
}

func (flow *KisFlow) GetName() string {
	return flow.Name
}

func (flow *KisFlow) GetThisFunction() kis.Function {
	return flow.ThisFunction
}

func (flow *KisFlow) GetThisFuncConf() *config.KisFuncConfig {
	return flow.ThisFunction.GetConfig()
}

func (flow *KisFlow) GetConnector() (kis.Connector, error) {
	if conn := flow.ThisFunction.GetConnector(); conn != nil {
		return conn, nil
	} else {
		return nil, errors.New("GetConnector(): Connector is nil")
	}
}

func (flow *KisFlow) GetConnConf() (*config.KisConnConfig, error) {
	if conn := flow.ThisFunction.GetConnector(); conn != nil {
		return conn.GetConfig(), nil
	} else {
		return nil, errors.New("GetConnConf(): Connector's conf is nil")
	}
}
func (flow *KisFlow) GetConfig() *config.KisFlowConfig {
	return flow.Conf
}

// GetFuncConfigByName 得到当前Flow的配置
func (flow *KisFlow) GetFuncConfigByName(funcName string) *config.KisFuncConfig {
	if f, ok := flow.Funcs[funcName]; ok {
		return f.GetConfig()
	} else {
		log.GetLogger().ErrorF("GetFuncConfigByName(): Function %s not found", funcName)
		return nil
	}
}

func (flow *KisFlow) Next(acts ...kis.ActionFunc) error {
	flow.action = kis.LoadActions(acts)

	return nil
}

func (flow *KisFlow) Fork(ctx context.Context) kis.Flow {
	config := flow.Conf

	// 通过之前的配置生成一个新的Flow
	newFlow := NewKisFlow(config)

	for _, fp := range flow.Conf.Flows {
		if _, ok := flow.funcParams[flow.Funcs[fp.FuncName].GetId()]; !ok {
			//当前function没有配置Params
			newFlow.Link(flow.Funcs[fp.FuncName].GetConfig(), nil)
		} else {
			newFlow.Link(flow.Funcs[fp.FuncName].GetConfig(), fp.Params)
		}
	}
	log.GetLogger().DebugFX(ctx, "=====>Flow Fork, oldFlow.funcParams = %+v\n", flow.funcParams)
	log.GetLogger().DebugFX(ctx, "=====>Flow Fork, newFlow.funcParams = %+v\n", newFlow.GetFuncParamsAllFuncs())

	return newFlow
}
