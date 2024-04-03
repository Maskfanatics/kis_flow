package kis

import (
	"context"
	"errors"
	"fmt"
	"kis-flow/common"
	"kis-flow/log"
	"reflect"
	"sync"
)

var _poolOnce sync.Once

// kisPool 用于管理全部的Function和Flow配置的池子
type kisPool struct {
	fnRouter funcRouter   // 全部的Function管理路由
	fnLock   sync.RWMutex // fnRouter 锁

	flowRouter flowRouter   // 全部的flow对象
	flowLock   sync.RWMutex // flowRouter 锁

	// +++++++++++++++++
	cInitRouter connInitRouter // 全部的Connector初始化路由
	ciLock      sync.RWMutex   // cInitRouter 锁

	cTree connTree     //全部Connector管理路由
	cLock sync.RWMutex // cTree 锁
	// +++++++++++++++++
}

// 单例
var _pool *kisPool

// Pool 单例构造
func Pool() *kisPool {
	_poolOnce.Do(func() {
		//创建kisPool对象
		_pool = new(kisPool)

		// fnRouter初始化
		_pool.fnRouter = make(funcRouter)

		// flowRouter初始化
		_pool.flowRouter = make(flowRouter)

		// +++++++++++++++++++++++++
		// connTree初始化
		_pool.cTree = make(connTree)
		_pool.cInitRouter = make(connInitRouter)
		// +++++++++++++++++++++++++
	})
	return _pool
}

func (pool *kisPool) AddFlow(name string, flow Flow) {
	pool.flowLock.Lock()
	defer pool.flowLock.Unlock()

	if _, ok := pool.flowRouter[name]; !ok {
		pool.flowRouter[name] = flow
	} else {
		errString := fmt.Sprintf("Pool AddFlow Repeat FlowName=%s\n", name)
		panic(errString)
	}

	log.GetLogger().InfoF("Add FlowRouter FlowName=%s\n", name)
}

func (pool *kisPool) GetFlow(name string) Flow {
	pool.flowLock.RLock()
	defer pool.flowLock.RUnlock()

	if flow, ok := pool.flowRouter[name]; ok {
		return flow
	} else {
		return nil
	}
}

// FaaS 注册 Function 计算业务逻辑, 通过Function Name 索引及注册
func (pool *kisPool) FaaS(fnName string, f Faas) {
	faasDesc, err := NewFaasDesc(fnName, f)
	if err != nil {
		panic(err)
	}
	pool.fnLock.Lock()
	defer pool.fnLock.Unlock()

	if _, ok := pool.fnRouter[fnName]; !ok {
		pool.fnRouter[fnName] = faasDesc
	} else {
		errString := fmt.Sprintf("KisPool Faas Request FuncName=%s", fnName)
		panic(errString)
	}
	log.GetLogger().InfoF("Add KisPool FuncName=%s", fnName)
}

// CallFunction 调度 Function
func (pool *kisPool) CallFunction(ctx context.Context, fnName string, flow Flow) error {
	if funcDesc, ok := pool.fnRouter[fnName]; ok {
		params := make([]reflect.Value, 0, funcDesc.ArgNum)

		for _, argType := range funcDesc.ArgsType {
			if isFlowType(argType) {
				params = append(params, reflect.ValueOf(flow))
			}

			if isContextType(argType) {
				params = append(params, reflect.ValueOf(ctx))
			}

			if isSliceType(argType) {
				value, err := funcDesc.Serialize.UnMarshal(flow.Input(), argType)
				if err != nil {
					log.GetLogger().ErrorFX(ctx, "funcDesc.Serialize.DecodeParam err=%v", err)
				} else {
					params = append(params, value)
					continue
				}
			}

			params = append(params, reflect.Zero(argType))
		}
		retValues := funcDesc.FuncValue.Call(params)

		ret := retValues[0].Interface()
		if ret == nil {
			return nil
		}

		return retValues[0].Interface().(error)
	}
	log.GetLogger().ErrorFX(ctx, "FuncName: %s Can not find in KisPool, Not Added.\n", fnName)

	return errors.New("FuncName: " + fnName + " Can not find in NsPool, Not Added.")
}

func (pool *kisPool) CaaSInit(cname string, c ConnInit) {
	pool.ciLock.Lock()
	defer pool.ciLock.Unlock()

	if _, ok := pool.cInitRouter[cname]; !ok {
		pool.cInitRouter[cname] = c
	} else {
		errString := fmt.Sprintf("KisPool Reg CaaSInit Repeat CName=%s\n", cname)
		panic(errString)
	}
	log.GetLogger().InfoF("Add KisPool CaaSInit CName=%s", cname)
}

// CallConnInit 调度 ConnInit
func (pool *kisPool) CallConnInit(conn Connector) error {
	pool.ciLock.RLock() // 读锁
	defer pool.ciLock.RUnlock()
	init, ok := pool.cInitRouter[conn.GetName()]
	fmt.Println(init)
	if !ok {
		panic(errors.New(fmt.Sprintf("init connector cname = %s not reg..", conn.GetName())))
	}

	return init(conn)
}

// CaaS 注册Connector Call业务
func (pool *kisPool) CaaS(cname string, fname string, mode common.KisMode, c CaaS) {
	pool.cLock.Lock() // 写锁
	defer pool.cLock.Unlock()

	if _, ok := pool.cTree[cname]; !ok {
		//cid 首次注册，不存在，创建二级树NsConnSL
		pool.cTree[cname] = make(connSL)

		//初始化各类型FunctionMode
		pool.cTree[cname][common.S] = make(connFuncRouter)
		pool.cTree[cname][common.L] = make(connFuncRouter)
	}

	if _, ok := pool.cTree[cname][mode][fname]; !ok {
		pool.cTree[cname][mode][fname] = c
	} else {
		errString := fmt.Sprintf("CaaS Repeat CName=%s, FName=%s, Mode =%s\n", cname, fname, mode)
		panic(errString)
	}

	log.GetLogger().InfoF("Add KisPool CaaS CName=%s, FName=%s, Mode =%s", cname, fname, mode)
}

// CallConnector 调度 Connector
func (pool *kisPool) CallConnector(ctx context.Context, flow Flow, conn Connector, args interface{}) (interface{}, error) {
	fn := flow.GetThisFunction()
	fnConf := fn.GetConfig()
	mode := common.KisMode(fnConf.FMode)

	if callback, ok := pool.cTree[conn.GetName()][mode][fnConf.FName]; ok {
		return callback(ctx, conn, fn, flow, args)
	}

	log.GetLogger().ErrorFX(ctx, "CName:%s FName:%s mode:%s Can not find in KisPool, Not Added.\n", conn.GetName(), fnConf.FName, mode)

	return nil, errors.New(fmt.Sprintf("CName:%s FName:%s mode:%s Can not find in KisPool, Not Added.", conn.GetName(), fnConf.FName, mode))
}

func (pool *kisPool) GetFlows() []Flow {
	pool.flowLock.RLock()
	defer pool.flowLock.RUnlock()

	var flows []Flow

	for _, flow := range pool.flowRouter {
		flows = append(flows, flow)
	}

	return flows
}
