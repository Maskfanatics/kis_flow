package flow

import (
	"context"
	"errors"
	"fmt"
	"kis-flow/common"
	"kis-flow/log"
)

func (flow *KisFlow) CommitRow(row interface{}) error {
	flow.buffer = append(flow.buffer, row)
	return nil
}

func (flow *KisFlow) commitSrcData(ctx context.Context) error {
	dataCnt := len(flow.buffer)
	batch := make(common.KisRowArr, 0, dataCnt)
	for _, rows := range flow.buffer {
		batch = append(batch, rows)
	}
	flow.clearData(flow.data)

	flow.data[common.FunctionIdFirstVirtual] = batch

	flow.buffer = flow.buffer[0:0]

	log.GetLogger().DebugFX(ctx, "====> After CommitSrcData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

func (flow *KisFlow) commitCurData(ctx context.Context) error {
	if len(flow.buffer) == 0 {
		flow.abort = true
		return nil
	}

	batch := make(common.KisRowArr, 0, len(flow.buffer))

	for _, rows := range flow.buffer {
		batch = append(batch, rows)
	}

	flow.data[flow.ThisFunctionId] = batch

	flow.buffer = flow.buffer[0:0]

	log.GetLogger().DebugFX(ctx, " ====> After commitCurData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil

}

func (flow *KisFlow) clearData(data common.KisDataMap) {
	for k := range data {
		delete(data, k)
	}
}

// getCurData 获取flow当前Function层级的输入数据
func (flow *KisFlow) getCurData() (common.KisRowArr, error) {
	if flow.PrevFunctionId == "" {
		return nil, errors.New(fmt.Sprintf("flow.PrevFunctionId is not set"))
	}

	if _, ok := flow.data[flow.PrevFunctionId]; !ok {
		return nil, errors.New(fmt.Sprintf("[%s] is not in flow.data", flow.PrevFunctionId))
	}

	return flow.data[flow.PrevFunctionId], nil
}

func (flow *KisFlow) commitReuseData(ctx context.Context) error {
	// 判断上层是否有结果数据, 如果没有则退出本次Flow Run循环
	if len(flow.data[flow.PrevFunctionId]) == 0 {
		flow.abort = true
		return nil
	}

	//本层结果数据等于上层结果数据(复用上层数据于本层)
	flow.data[flow.ThisFunctionId] = flow.data[flow.PrevFunctionId]

	//清空缓冲区（如果是 ReuseData，那么本层提交的所有数据，都不会携带到下一层）
	flow.buffer = flow.buffer[0:0]

	log.GetLogger().DebugFX(ctx, " ====> After commitReuseData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}
