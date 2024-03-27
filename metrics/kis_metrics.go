package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/log"
	"net/http"
)

type kisMetrics struct {
	DataTotal             prometheus.Counter
	FlowDataTotal         *prometheus.GaugeVec
	FlowScheduleCntsTotal *prometheus.GaugeVec
	FuncScheduleCntsTotal *prometheus.GaugeVec
	FunctionDuration      *prometheus.HistogramVec
	FlowDuration          *prometheus.HistogramVec
}

var Metrics *kisMetrics

func RunMetricsService(serverAddr string) error {
	http.Handle(common.METRICS_ROUTE, promhttp.Handler())

	err := http.ListenAndServe(serverAddr, nil)

	if err != nil {
		log.GetLogger().ErrorF("RunMetricsService err = %s\n", err)
	}

	return err
}

func InitMetrics() {
	Metrics = new(kisMetrics)

	Metrics.DataTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: common.COUNTER_KISFLOW_DATA_TOTAL_NAME,
		Help: common.COUNTER_KISFLOW_DATA_TOTAL_HELP,
	})

	Metrics.FlowDataTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: common.GANGE_FLOW_DATA_TOTAL_NAME,
		Help: common.COUNTER_KISFLOW_DATA_TOTAL_HELP,
	},
		[]string{common.LABEL_FLOW_NAME},
	)

	Metrics.FlowScheduleCntsTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: common.GANGE_FLOW_SCHE_CNTS_NAME,
		Help: common.GANGE_FLOW_SCHE_CNTS_HELP,
	}, []string{common.LABEL_FLOW_NAME})

	Metrics.FuncScheduleCntsTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: common.GANGE_FUNC_SCHE_CNTS_NAME,
		Help: common.GANGE_FUNC_SCHE_CNTS_HELP,
	}, []string{common.LABEL_FUNCTION_NAME, common.LABEL_FUNCTION_MODE})

	Metrics.FunctionDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    common.HISTOGRAM_FUNCTION_DURATION_NAME,
		Help:    common.HISTOGRAM_FUNCTION_DURATION_HELP,
		Buckets: []float64{0.005, 0.01, 0.03, 0.08, 0.1, 0.5, 1.0, 5.0, 10, 100, 1000, 5000, 30000}, //单位ms,最大半分钟
	}, []string{common.LABEL_FUNCTION_NAME, common.LABEL_FUNCTION_MODE})

	Metrics.FlowDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    common.HISTOGRAM_FLOW_DURATION_NAME,
		Help:    common.HISTOGRAM_FLOW_DURATION_HELP,
		Buckets: []float64{0.005, 0.01, 0.03, 0.08, 0.1, 0.5, 1.0, 5.0, 10, 100, 1000, 5000, 30000, 60000}, //单位ms,最大1分钟
	}, []string{common.LABEL_FLOW_NAME})

	prometheus.MustRegister(Metrics.DataTotal)
	prometheus.MustRegister(Metrics.FlowDataTotal)
	prometheus.MustRegister(Metrics.FlowScheduleCntsTotal)
	prometheus.MustRegister(Metrics.FuncScheduleCntsTotal)
	prometheus.MustRegister(Metrics.FunctionDuration)
	prometheus.MustRegister(Metrics.FlowDuration)
}

func RunMetrics() {
	InitMetrics()
	if config.GlobalConfig.EnableProm && config.GlobalConfig.PrometheusListen {
		go RunMetricsService(config.GlobalConfig.PrometheusServe)
	}
}
