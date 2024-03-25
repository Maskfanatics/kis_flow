package kis

type Action struct {
	DataReuse bool

	Abort bool

	ForceEntryNext bool

	JumpFunc string
}

type ActionFunc func(ops *Action)

func LoadActions(acts []ActionFunc) Action {
	action := Action{}

	if acts == nil {
		return action
	}

	for _, act := range acts {
		act(&action)
	}

	return action
}

func ActionAbort(action *Action) {
	action.Abort = true
}

func ActionDataReuse(action *Action) {
	action.DataReuse = true
}

func ActionForceEntry(action *Action) {
	action.ForceEntryNext = true
}

func ActionJumpFunc(funcName string) ActionFunc {
	return func(action *Action) {
		action.JumpFunc = funcName
	}
}
