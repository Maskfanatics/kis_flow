package common

import "time"

type KisMode string
type KisOnOff int
type KisConnType string

const (
	// V 为校验特征的KisFunction,
	// 主要进行数据的过滤，验证，字段梳理，幂等等前置数据处理
	V KisMode = "Verify"

	// S 为存储特征的KisFunction,
	// S会通过NsConnector进行将数据进行存储，数据的临时声明周期为NsWindow
	S KisMode = "Save"

	// L 为加载特征的KisFunction，
	// L会通过KisConnector进行数据加载，通过该Function可以从逻辑上与对应的S Function进行并流
	L KisMode = "Load"

	// C 为计算特征的KisFunction,
	// C会通过KisFlow中的数据计算，生成新的字段，将数据流传递给下游S进行存储，或者自己也已直接通过KisConnector进行存储
	C KisMode = "Calculate"

	// E 为扩展特征的KisFunction，
	// 作为流式计算的自定义特征Function，如，Notify 调度器触发任务的消息发送，删除一些数据，重置状态等。
	E KisMode = "Expand"

	REDIS KisConnType = "Redis"

	MYSQL KisConnType = "MySQL"

	KAFKA KisConnType = "Kafka"

	KisEnable KisOnOff = 1

	KisDisable KisOnOff = 0

	FlowEnable KisOnOff = 1

	FlowDisable KisOnOff = 0
)
const (
	KisIdTypeFlow       = "flow"
	KisIdTypeConnnector = "conn"
	KisIdTypeFunction   = "func"
	KisIdTypeGlobal     = "global"
	KisIdJoinChar       = "-"
)

const (
	FunctionIdFirstVirtual = "FunctionIdFirstVirtual"
	FunctionIdLastVirtual  = "FunctionalIdLastVirtual"
)

//cache

const (
	DefaultFlowCacheCleanUp = 5 //min

	DefaultExpiration time.Duration = 0
)

//Metrics

const (
	METRICS_ROUTE string = "/metrics"

	LABEL_FLOW_NAME     string = "flow_name"
	LABEL_FLOW_ID       string = "flow_id"
	LABEL_FUNCTION_NAME string = "func_name"
	LABEL_FUNCTION_MODE string = "func_mode"

	COUNTER_KISFLOW_DATA_TOTAL_NAME string = "kisflow_data_total"
	COUNTER_KISFLOW_DATA_TOTAL_HELP string = "KisFlow全部Flow的数据总量"

	GANGE_FLOW_DATA_TOTAL_NAME string = "flow_data_total"
	GANGE_FLOW_DATA_TOTAL_HELP string = "KisFlow各个FlowID数据流的数据数量总量"

	GANGE_FLOW_SCHE_CNTS_NAME string = "flow_schedule_cnts"
	GANGE_FLOW_SCHE_CNTS_HELP string = "KisFlow各个FlowID被调度的次数"

	GANGE_FUNC_SCHE_CNTS_NAME string = "func_schedule_cnts"
	GANGE_FUNC_SCHE_CNTS_HELP string = "KisFlow各个FunctionID被调度的次数"

	HISTOGRAM_FUNCTION_DURATION_NAME string = "func_run_duration"
	HISTOGRAM_FUNCTION_DURATION_HELP string = "Function执行耗时"

	HISTOGRAM_FLOW_DURATION_NAME string = "flow_run_duration"
	HISTOGRAM_FLOW_DURATION_HELP string = "Flow执行耗时"
)
