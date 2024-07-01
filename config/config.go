package config

// WriteConfig Write配置结构
type WriteConfig struct {
	Type   string `toml:"type"`
	URL    string `toml:"url,omitempty"`
	Broker string `toml:"broker,omitempty"`
	Topic  string `toml:"topic,omitempty"`
}

// PluginConfig 插件配置结构
type PluginConfig struct {
	Name           string `toml:"name"`
	ReportInterval int    `toml:"report_interval"`
}

// LogConfig 包含日志配置
type LogConfig struct {
	Level  string `toml:"level"`
	Output string `toml:"output"`
}

// Config 配置文件结构
type Config struct {
	QueryFromMemory bool           `toml:"query_from_memory"`
	Write           []WriteConfig  `toml:"write"`
	Plugin          []PluginConfig `toml:"plugin"`
	Log             LogConfig      `toml:"log"`
}
