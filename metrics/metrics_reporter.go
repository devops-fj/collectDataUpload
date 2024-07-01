package metrics

import (
	"net/http"

	"github.com/devops-fj/collectDataUpload/reporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ReporterMetrics 定义 Reporter 相关 Prometheus metrics
type ReporterMetrics struct {
	reportsTotal *prometheus.CounterVec
}

// NewReporterMetrics 创建一个新的 ReporterMetrics 实例
func NewReporterMetrics() *ReporterMetrics {
	return &ReporterMetrics{
		reportsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "reporter_reports_total",
				Help: "Total number of reports by reporter type",
			},
			[]string{"reporter_type"},
		),
	}
}

// Register 注册 Prometheus metrics
func (m *ReporterMetrics) Register() {
	prometheus.MustRegister(m.reportsTotal)
}

// IncrementReports 增加报告计数
func (m *ReporterMetrics) IncrementReports(reporterType string) {
	m.reportsTotal.With(prometheus.Labels{"reporter_type": reporterType}).Inc()
}

// Handler 返回 Prometheus metrics handler
func (m *ReporterMetrics) Handler() http.Handler {
	return promhttp.Handler()
}

// ReporterObserver 用于观察报告并更新 Prometheus metrics
type ReporterObserver struct {
	metrics *ReporterMetrics
}

// NewReporterObserver 创建一个新的 ReporterObserver 实例
func NewReporterObserver(metrics *ReporterMetrics) *ReporterObserver {
	return &ReporterObserver{metrics: metrics}
}

// ObserveReports 观察报告并更新 Prometheus metrics
func (o *ReporterObserver) ObserveReports(reports []reporter.Report) {
	for _, report := range reports {
		o.metrics.IncrementReports(report.Type)
	}
}
