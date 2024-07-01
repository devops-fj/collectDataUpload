package ano

import (
	"fmt"
	"log"
	"time"
)

type AnotherPlugin struct {
	running bool
}

func NewAnotherPlugin() *AnotherPlugin {
	return &AnotherPlugin{}
}

func (p *AnotherPlugin) Collect() (map[string]interface{}, error) {
	if !p.running {
		return nil, fmt.Errorf("plugin not running")
	}
	data := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"value":     10000, // 模拟的数据
	}
	log.Println("another_plugin collect data:", data)
	return data, nil
}

func (p *AnotherPlugin) Start() error {
	p.running = true
	return nil
}

func (p *AnotherPlugin) Stop() error {
	p.running = false
	return nil
}

// Name 返回插件名称
func (p *AnotherPlugin) Name() string {
	return "another_plugin"
}
