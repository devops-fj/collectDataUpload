package example

import (
	"fmt"
	"time"
)

type ExamplePlugin struct {
	running bool
}

func NewExamplePlugin() *ExamplePlugin {
	return &ExamplePlugin{}
}
func (p *ExamplePlugin) Start() error {
	p.running = true
	return nil
}

func (p *ExamplePlugin) Stop() error {
	p.running = false
	return nil
}

func (p *ExamplePlugin) Name() string {
	return "example_plugin"
}

func (p *ExamplePlugin) Collect() (map[string]interface{}, error) {
	if !p.running {
		return nil, fmt.Errorf("plugin not running")
	}
	data := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"value":     42, // 模拟的数据
	}
	return data, nil
}
