package business

import (
	"errors"
	"time"
)

const (
	ReachMaxLoopTime = iota // 达到最大的轮询次数
	SelfReturn              // 自己返回
)

// ExecuteReturnReason 执行结果，返回原因
type ExecuteReturnReason int

// Executor 执行的轮询动作，注意currentTime下标从1开始
type Executor func(currentTime int) bool

// ExecuteResult 返回值
type ExecuteResult struct {
	HasLoopCount int                 // 已经轮询的次数，从1开始
	ReturnReason ExecuteReturnReason // 返回原因
}

// Retryer 重试器
type Retryer struct {
	LoopCount            int           // 目标轮询次数
	LoopOnceTimeDuration time.Duration // 每次轮询的时间间隔
}

// NewRetryer 构建retryer
func NewRetryer(duration time.Duration, loopCount int) (*Retryer, error) {
	if duration.Seconds() < 1 || loopCount < 1 {
		return nil, errors.New("轮询时间间隔不能小于1秒，轮询次数不能小于1")
	}
	return &Retryer{
		LoopCount:            loopCount,
		LoopOnceTimeDuration: duration,
	}, nil
}

// Start 开始进行处理任务工作
func (r Retryer) Start(f Executor) *ExecuteResult {
	if !f(1) {
		return &ExecuteResult{
			HasLoopCount: 1,
			ReturnReason: SelfReturn,
		}
	}

	ticker := time.NewTicker(r.LoopOnceTimeDuration)
	defer ticker.Stop()
	length := r.LoopCount + 1
	for i := 2; i < length; i++ {
		<-ticker.C
		if !f(i) {
			return &ExecuteResult{
				HasLoopCount: i,
				ReturnReason: SelfReturn,
			}
		}
	}
	return &ExecuteResult{
		HasLoopCount: r.LoopCount,
		ReturnReason: ReachMaxLoopTime,
	}
}
