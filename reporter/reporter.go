package reporter

import (
	"fmt"
	"sync"
)

// Reporter 上报器接口
type Reporter interface {
	Report(data map[string]interface{}) error
}

// Report 表示一个上报的数据
type Report struct {
	Type string                 // 上报的类型，例如 "http" 或 "mq"
	Data map[string]interface{} // 上报的数据
}

// Observer 用于观察报告
type Observer interface {
	ObserveReports(reports []Report)
}

// MultiReporter 多上报器，可以同时上报到多个目标
type MultiReporter struct {
	reporters []Reporter
	observer  Observer // 添加观察器
}

// NewMultiReporter 创建一个新的 MultiReporter 实例
func NewMultiReporter(reporters ...Reporter) *MultiReporter {
	return &MultiReporter{
		reporters: reporters,
	}
}

// Report 实现 Reporter 接口
//
//	func (mr *MultiReporter) Report(data map[string]interface{}) error {
//		reports := make([]Report, len(mr.reporters))
//		for i, r := range mr.reporters {
//			err := r.Report(data)
//			reports[i] = Report{
//				Type: getType(r),
//				Data: data,
//			}
//			if err != nil {
//				// fmt.Println(err)
//				return err
//			}
//		}
//		// 添加观察器的处理
//		if mr.observer != nil {
//			mr.observer.ObserveReports(reports)
//		}
//		return nil
//	}
//
// Report 实现 Reporter 接口
// func (mr *MultiReporter) Report(data map[string]interface{}) error {
// 	var errs []error
// 	reports := make([]Report, len(mr.reporters))
//
// 	for i, r := range mr.reporters {
// 		err := r.Report(data)
// 		reports[i] = Report{
// 			Type: getType(r),
// 			Data: data,
// 		}
// 		if err != nil {
// 			errs = append(errs, err)
// 			// 可以选择在这里记录日志或其他处理
// 		}
// 	}
//
// 	// 添加观察器的处理
// 	if mr.observer != nil {
// 		mr.observer.ObserveReports(reports)
// 	}
//
// 	// 如果errs不为空，则返回多个错误；否则返回nil
// 	if len(errs) > 0 {
// 		return fmt.Errorf("errors occurred during reporting: %v", errs)
// 	}
// 	return nil
// }

// Report 实现 Reporter 接口
func (mr *MultiReporter) Report(data map[string]interface{}) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(mr.reporters))
	reports := make([]Report, len(mr.reporters))

	for i, r := range mr.reporters {
		wg.Add(1)
		go func(index int, reporter Reporter) {
			defer wg.Done()

			err := reporter.Report(data)
			reports[index] = Report{
				Type: getType(reporter),
				Data: data,
			}
			if err != nil {
				errCh <- err
				// 可以选择在这里记录日志或其他处理
			}
		}(i, r)
	}
	wg.Wait()
	close(errCh)
	// 如果有错误，则从通道中读取错误并返回
	if len(errCh) > 0 {
		var errs []error
		for err := range errCh {
			errs = append(errs, err)
		}
		return fmt.Errorf("errors occurred during reporting: %v", errs)
	}
	// 添加观察器的处理
	if mr.observer != nil {
		mr.observer.ObserveReports(reports)
	}

	return nil
}

// getType 获取上报器类型
func getType(r Reporter) string {
	switch r.(type) {
	case *HTTPReporter:
		return "http"
	case *MQReporter:
		return "mq"
	default:
		return "unknown"
	}
}

// SetObserver 设置观察器
func (mr *MultiReporter) SetObserver(observer Observer) {
	mr.observer = observer
}
