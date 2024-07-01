package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPReporter HTTP上报器
type HTTPReporter struct {
	URL string
}

// NewHTTPReporter 创建一个新的HTTPReporter实例
func NewHTTPReporter(url string) *HTTPReporter {
	return &HTTPReporter{URL: url}
}

// Report 实现Reporter接口的上报方法
func (r *HTTPReporter) Report(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(r.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
