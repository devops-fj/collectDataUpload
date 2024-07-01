package plugin

// Plugin 插件接口
type Plugin interface {
	Start() error
	Stop() error
	Collect() (map[string]interface{}, error)
	Name() string
}
