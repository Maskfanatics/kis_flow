package kis

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type Faas interface{}

type FaaSDesc struct {
	Serialize
	FnName    string         // Function名称
	f         interface{}    // FaaS 函数
	fName     string         // 函数名称
	ArgsType  []reflect.Type // 函数参数类型（集合）
	ArgNum    int            // 函数参数个数
	FuncType  reflect.Type   // 函数类型
	FuncValue reflect.Value  // 函数值(函数地址)
}

func NewFaasDesc(fnName string, f Faas) (*FaaSDesc, error) {

	var serializeImpl Serialize

	// 传入的回调函数FaaS,函数值(函数地址)
	funcValue := reflect.ValueOf(f)

	// 传入的回调函数FaaS 类型
	funcType := funcValue.Type()

	// 判断传递的FaaS指针是否是函数类型
	if !isFuncType(funcType) {
		return nil, fmt.Errorf("provided FaaS type is %s, not a function", funcType.Name())
	}

	// 判断传递的FaaS函数是否有返回值类型是只包括(error)
	if funcType.NumOut() != 1 || funcType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("function must have exactly one return value of type error")
	}

	// FaaS函数的参数类型集合
	argsType := make([]reflect.Type, funcType.NumIn())

	// 获取FaaS的函数名称
	fullName := runtime.FuncForPC(funcValue.Pointer()).Name()

	// 确保  FaaS func(context.Context, Flow, ...interface{}) error 形参列表，存在context.Context 和 kis.Flow

	// 是否包含kis.Flow类型的形参
	containsKisFlow := false
	// 是否包含context.Context类型的形参
	containsCtx := false

	// 遍历FaaS的形参类型
	for i := 0; i < funcType.NumIn(); i++ {

		// 取出第i个形式参数类型
		paramType := funcType.In(i)

		if isFlowType(paramType) {
			// 判断是否包含kis.Flow类型的形参
			containsKisFlow = true

		} else if isContextType(paramType) {
			// 判断是否包含context.Context类型的形参
			containsCtx = true

		} else if isSliceType(paramType) {
			// 判断是否包含Slice类型的形参

			// 获取当前参数Slice的元素类型
			itemType := paramType.Elem()

			// 如果当前参数是一个指针类型，则获取指针指向的结构体类型
			if itemType.Kind() == reflect.Ptr {
				itemType = itemType.Elem() // 获取指针指向的结构体类型
			}

			if isSerialize(itemType) {
				serializeImpl = reflect.New(itemType).Interface().(Serialize)
			} else {
				serializeImpl = defaultSerialize

			}
		} else {
			// Other types are not supported...
		}

		// 将当前形参类型追加到argsType集合中
		argsType[i] = paramType
	}

	if !containsKisFlow {
		// 不包含kis.Flow类型的形参，返回错误
		return nil, errors.New("function parameters must have kis.Flow param, please use FaaS type like: [type FaaS func(context.Context, Flow, ...interface{}) error]")
	}

	if !containsCtx {
		// 不包含context.Context类型的形参，返回错误
		return nil, errors.New("function parameters must have context, please use FaaS type like: [type FaaS func(context.Context, Flow, ...interface{}) error]")
	}

	// 返回FaaSDesc描述实例
	return &FaaSDesc{
		Serialize: serializeImpl,
		FnName:    fnName,
		f:         f,
		fName:     fullName,
		ArgsType:  argsType,
		ArgNum:    len(argsType),
		FuncType:  funcType,
		FuncValue: funcValue,
	}, nil
}

func isFuncType(paramType reflect.Type) bool {
	return paramType.Kind() == reflect.Func
}

func isFlowType(paramType reflect.Type) bool {
	var flowInterfaceType = reflect.TypeOf((*Flow)(nil)).Elem()
	return paramType.Implements(flowInterfaceType)
}

func isContextType(paramType reflect.Type) bool {
	typeName := paramType.Name()

	return strings.Contains(typeName, "Context")
}

func isSliceType(paramType reflect.Type) bool {
	return paramType.Kind() == reflect.Slice
}
