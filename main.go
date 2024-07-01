// main.go
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/devops-fj/collectDataUpload/config"
	"github.com/devops-fj/collectDataUpload/logger"
	"github.com/devops-fj/collectDataUpload/metrics"
	"github.com/devops-fj/collectDataUpload/plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 初始化配置
	cfg := initConfig()

	// 初始化日志管理器
	logger.InitLogger(cfg.Log.Level, cfg.Log.Output)

	// 初始化 Metrics
	metric := metrics.NewMetrics()
	metric.Register()

	// 初始化 ReporterMetrics
	reporterMetrics := metrics.NewReporterMetrics()
	reporterMetrics.Register()

	// 启动 Metrics HTTP 服务器
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	// 创建插件管理器
	pm := plugin.NewPluginManager(cfg, metric, reporterMetrics)

	// 注册所有插件
	pm.RegisterPlugins()

	// 启动插件
	pm.StartPlugins()

	// 注册插件管理器的 HTTP 处理器
	http.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// 获取所有插件的数据
			data := pm.GetPluginsData()

			// 将数据转换为 JSON 格式并写入响应
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	})

	// 等待中断信号以优雅地关闭插件
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	// 停止插件
	pm.StopPlugins()
}

// 初始化配置文件
func initConfig() *config.Config {
	var cfg config.Config

	// 读取配置文件
	if _, err := toml.DecodeFile("config/config.toml", &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
