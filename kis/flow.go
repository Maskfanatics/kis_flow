package kis

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
	"time"
)

type Flow interface {
	Run(ctx context.Context) error

	Link(fConfig *config.KisFuncConfig, fParams config.FParam) error

	CommitRow(row interface{}) error

	Input() common.KisRowArr

	// GetName 得到Flow的名称
	GetName() string
	// GetThisFunction 得到当前正在执行的Function
	GetThisFunction() Function
	// GetThisFuncConf 得到当前正在执行的Function的配置
	GetThisFuncConf() *config.KisFuncConfig

	// +++++++++++++++++++++++++++++++++
	GetConnector() (Connector, error)
	// GetConnConf 得到当前正在执行的Function的Connector的配置
	GetConnConf() (*config.KisConnConfig, error)
	// +++++++++++++++++++++++++++++++++
	// GetConfig 得到当前Flow的配置
	GetConfig() *config.KisFlowConfig
	// GetFuncConfigByName 得到当前Flow的配置
	GetFuncConfigByName(funcName string) *config.KisFuncConfig

	//  --- KisFlow Action ---
	// Next 当前Flow执行到的Function进入下一层Function所携带的Action动作
	Next(acts ...ActionFunc) error

	// ++++++++++++++++++++++++++++++++++++++++
	// GetCacheData 得到当前Flow的缓存数据
	GetCacheData(key string) interface{}
	// SetCacheData 设置当前FLow的缓存数据
	SetCacheData(key string, value interface{}, Exp time.Duration)

	// ++++++++++++++++++++++++++++
	// GetMetaData 得到当前Flow的临时数据
	GetMetaData(key string) interface{}
	// SetMetaData 设置当前Flow的临时数据
	SetMetaData(key string, value interface{})

	GetFuncParam(key string) string

	GetFuncParamAll() config.FParam

	Fork(ctx context.Context) Flow

	GetFuncParamsAllFuncs() map[string]config.FParam
}
