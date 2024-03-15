package flow

import (
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/id"
)

type KisFlow struct {
	Id   string
	Name string

	//TODO
}

func NewKisFlow(config *config.KisFlowConfig) *KisFlow {
	flow := new(KisFlow)

	flow.Id = id.KisId(common.KisIdTypeFlow)
	flow.Name = config.FlowName

	return flow
}
