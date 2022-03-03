package orbit

import (
	"context"
	"os"
	"time"
)

// Option 选项闭包函数
type Option func(*options)

// options 自定义选项
type options struct {
	network string
	ip      string
	port    int

	conns  int
	pool   int
	tasks  int
	packet uint32

	ctx     context.Context
	signals []os.Signal
	timeout time.Duration

	router Router
}

// WithNetwork 网络
func WithNetwork(network string) Option {
	return func(o *options) {
		o.network = network
	}
}

// WithIP 地址
func WithIP(ip string) Option {
	return func(o *options) {
		o.ip = ip
	}
}

// WithPort 端口
func WithPort(port int) Option {
	return func(o *options) {
		o.port = port
	}
}

// WithMaxConns 最大连接数
func WithMaxConns(conns int) Option {
	return func(o *options) {
		o.conns = conns
	}
}

// WithMaxWorkerPoolSize 工作池最大数
func WithMaxWorkerPoolSize(size int) Option {
	return func(o *options) {
		o.pool = size
	}
}

// WithMaxWorkerTasksQueueLength 工作池任务队列最大长度
func WithMaxWorkerTasksQueueLength(length int) Option {
	return func(o *options) {
		o.tasks = length
	}
}

// WithMaxMessagePacketSize 最大消息数据包
func WithMaxMessagePacketSize(size uint32) Option {
	return func(o *options) {
		o.packet = size
	}
}

// WithContext 上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithSignal 停止服务信号
func WithSignal(signals ...os.Signal) Option {
	return func(o *options) {
		o.signals = signals
	}
}

// WithStopTimeout 超时时间
func WithStopTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithRouter 路由
func WithRouter(r Router) Option {
	return func(o *options) {
		o.router = r
	}
}
