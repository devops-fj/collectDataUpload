package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics 定义 Prometheus metrics
type Metrics struct {
	pluginRuns *prometheus.CounterVec
}

// NewMetrics 创建一个新的 Metrics 实例
func NewMetrics() *Metrics {
	return &Metrics{
		pluginRuns: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "plugin_runs_total",
				Help: "Total number of plugin runs",
			},
			[]string{"plugin"},
		),
	}
}

// Register 注册 Prometheus metrics
func (m *Metrics) Register() {
	prometheus.MustRegister(m.pluginRuns)
}

// IncrementPluginRuns 增加插件运行计数
func (m *Metrics) IncrementPluginRuns(pluginName string) {
	m.pluginRuns.With(prometheus.Labels{"plugin": pluginName}).Inc()
}

// Handler 返回 Prometheus metrics handler
func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}
