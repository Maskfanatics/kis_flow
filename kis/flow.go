package kis

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
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
}
