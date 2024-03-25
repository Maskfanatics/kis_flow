package flow

import (
	"context"
	"errors"
	"fmt"
	"kis-flow/kis"
)

func (flow *KisFlow) dealAction(ctx context.Context, fn kis.Function) (kis.Function, error) {

	if flow.action.DataReuse {
		if err := flow.commitReuseData(ctx); err != nil {
			return nil, err
		}
	} else {
		if err := flow.commitCurData(ctx); err != nil {
			return nil, err
		}
	}

	//ForceEntryNext function
	if flow.action.ForceEntryNext {
		if err := flow.commitVoidData(ctx); err != nil {
			return nil, err
		}
		flow.abort = false
	}

	if flow.action.JumpFunc != "" {
		if _, ok := flow.Funcs[flow.action.JumpFunc]; !ok {
			return nil, errors.New(fmt.Sprintf("Flow Jump -> %s is not in Flow", flow.action.JumpFunc))
		}

		jumpFunction := flow.Funcs[flow.action.JumpFunc]

		flow.PrevFunctionId = jumpFunction.GetPrevId()

		fn = jumpFunction

		flow.abort = false

	} else {
		flow.PrevFunctionId = flow.ThisFunctionId
		fn = fn.Next()
	}

	if flow.action.Abort {
		flow.abort = true
	}

	flow.action = kis.Action{}

	return fn, nil
}
