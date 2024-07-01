package plugin

import (
	"fmt"
	"sync"
	"time"

	"github.com/devops-fj/collectDataUpload/config"
	"github.com/devops-fj/collectDataUpload/metrics"
	"github.com/devops-fj/collectDataUpload/plugins/ano"
	"github.com/devops-fj/collectDataUpload/plugins/example"
	"github.com/devops-fj/collectDataUpload/reporter"
)

const DefaultReportInterval = 10 // 默认上报周期为10秒

// PluginManager 插件管理器
type PluginManager struct {
	plugins         []Plugin
	wg              sync.WaitGroup
	stopCh          chan struct{}
	multiReporter   *reporter.MultiReporter
	config          *config.Config
	metrics         *metrics.Metrics
	reporterMetrics *metrics.ReporterMetrics
	mu              sync.Mutex
	inMemoryData    map[string]map[string]interface{} // 内存数据存储
	queryFromMemory bool                              // 是否从内存中查询数据
}

// NewPluginManager 创建一个新的 PluginManager 实例
func NewPluginManager(cfg *config.Config, metrics *metrics.Metrics, reporterMetrics *metrics.ReporterMetrics) *PluginManager {
	return &PluginManager{
		stopCh:          make(chan struct{}),
		multiReporter:   setupReporters(cfg, reporterMetrics),
		config:          cfg,
		metrics:         metrics,
		reporterMetrics: reporterMetrics,
		inMemoryData:    make(map[string]map[string]interface{}),
		queryFromMemory: cfg.QueryFromMemory, // 从配置文件中读取是否从内存中查询数据
	}
}

// RegisterPlugins 注册所有插件
func (pm *PluginManager) RegisterPlugins() {
	pm.plugins = append(pm.plugins, example.NewExamplePlugin())
	pm.plugins = append(pm.plugins, ano.NewAnotherPlugin())
}

// GetPluginsData 返回所有插件的数据列表
func (pm *PluginManager) GetPluginsData() []map[string]interface{} {
	if pm.queryFromMemory {
		pm.mu.Lock()
		defer pm.mu.Unlock()
		var data []map[string]interface{}
		for _, d := range pm.inMemoryData {
			data = append(data, d)
		}
		return data
	}

	var wg sync.WaitGroup
	dataChan := make(chan map[string]interface{}, len(pm.plugins))

	for _, p := range pm.plugins {
		wg.Add(1)
		go func(plugin Plugin) {
			defer wg.Done()
			data, err := plugin.Collect()
			if err != nil {
				println("Failed to collect data from plugin:", err)
				return
			}
			dataChan <- data
		}(p)
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	var data []map[string]interface{}
	for d := range dataChan {
		data = append(data, d)
	}

	return data
}

// StartPlugins 启动所有插件
func (pm *PluginManager) StartPlugins() {
	for _, p := range pm.plugins {
		pm.wg.Add(1)
		go pm.runPlugin(p)
	}
}

// runPlugin 在单独的协程中运行插件
func (pm *PluginManager) runPlugin(p Plugin) {
	defer pm.wg.Done()

	err := p.Start()
	if err != nil {
		println("Failed to start plugin:", err)
		return
	}

	// 立即采集数据
	pm.collectAndReport(p)

	reportInterval := pm.getReportInterval(p.Name())
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pm.collectAndReport(p)
		case <-pm.stopCh:
			err := p.Stop()
			if err != nil {
				println("Failed to stop plugin:", err)
			}
			return
		}
	}
}

// collectAndReport 采集并上报数据
func (pm *PluginManager) collectAndReport(p Plugin) {
	data, err := p.Collect()
	if err != nil {
		println("Failed to collect data from plugin:", err)
	} else {
		println("Collected data:", data)
		err := pm.multiReporter.Report(data)
		if err != nil {
			fmt.Println("Failed to report data:", err)
		}
		pm.metrics.IncrementPluginRuns(p.Name())

		// 存储数据到内存
		pm.mu.Lock()
		pm.inMemoryData[p.Name()] = data
		pm.mu.Unlock()
	}
}

// getReportInterval 获取插件的上报周期
func (pm *PluginManager) getReportInterval(pluginName string) int {
	for _, pluginCfg := range pm.config.Plugin {
		if pluginCfg.Name == pluginName {
			if pluginCfg.ReportInterval > 0 {
				return pluginCfg.ReportInterval
			}
		}
	}
	return DefaultReportInterval
}

// StopPlugins 停止所有插件
func (pm *PluginManager) StopPlugins() {
	close(pm.stopCh)
	pm.wg.Wait()
}

// setupReporters 根据配置文件设置上报器
func setupReporters(cfg *config.Config, reporterMetrics *metrics.ReporterMetrics) *reporter.MultiReporter {
	var reporters []reporter.Reporter

	for _, writeCfg := range cfg.Write {
		switch writeCfg.Type {
		case "http":
			if writeCfg.URL != "" {
				httpReporter := reporter.NewHTTPReporter(writeCfg.URL)
				reporters = append(reporters, httpReporter)
			}
		case "mq":
			if writeCfg.Broker != "" && writeCfg.Topic != "" {
				mqReporter := reporter.NewMQReporter(writeCfg.Broker, writeCfg.Topic)
				reporters = append(reporters, mqReporter)
			}
		}
	}

	multiReporter := reporter.NewMultiReporter(reporters...)
	observer := metrics.NewReporterObserver(reporterMetrics)
	multiReporter.SetObserver(observer)

	return multiReporter
}
