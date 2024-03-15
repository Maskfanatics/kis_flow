package config

import "kis-flow/common"

type KisFlowFunctionParam struct {
	FuncName string `yaml:"fname"`
	Params   FParam `yaml:"params"`
}

type KisFlowConfig struct {
	KisType  string                 `yaml:"kistype"`
	Status   int                    `yaml:"status"`
	FlowName string                 `yaml:"flow_name"`
	Flows    []KisFlowFunctionParam `yaml:"flows"`
}

func NewFlowConfig(flowName string, enable common.KisOnOff) *KisFlowConfig {
	config := new(KisFlowConfig)
	config.FlowName = flowName
	config.Flows = make([]KisFlowFunctionParam, 0)
	config.Status = int(enable)
	return config
}

func (fConfig *KisFlowConfig) AppendFunctionConfig(params KisFlowFunctionParam) {
	fConfig.Flows = append(fConfig.Flows, params)
}
