package orbit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Server 监听者接口
type Server interface {
	Run() error
	Shutdown() error
}

// listener 监听器结构体
type listener struct {
	opts options
	lis  *net.TCPListener

	mgr    Manager
	router Router
	work   Worker

	ctx    context.Context
	cancel context.CancelFunc
}

// New 实例化监听器
func New(opts ...Option) Server {
	// 初始化默认配置
	o := options{
		network: "tcp",
		ip:      "127.0.0.1",
		port:    62817,
		pool:    8,
		conns:   512,
		tasks:   1024,
		packet:  4096,
		timeout: 1 * time.Second,
		ctx:     context.Background(),
	}

	// 加载自定义配置
	for _, opt := range opts {
		opt(&o)
	}

	if o.router == nil {
		panic("router is nil")
	}

	ctx, cancel := context.WithCancel(o.ctx)
	return &listener{
		opts: o,
		mgr:  &manager{conns: make(map[string]Connection)},
		work: &worker{
			poolSize:  o.pool,
			taskLen:   o.tasks,
			taskQueue: make([]chan *Context, o.pool),
			router:    o.router,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run 开始监听服务
func (l *listener) Run() error {
	log.Println(fmt.Sprintf("[ LISTENER ] startup"))

	// 根据上下文信息做停止监听处理
	wg := sync.WaitGroup{}
	eg, ctx := errgroup.WithContext(l.ctx)
	eg.Go(func() error {
		<-ctx.Done()
		return l.off()
	})

	// 启动监听服务
	wg.Add(1)
	eg.Go(func() error {
		wg.Done()
		return l.on()
	})
	wg.Wait()

	// 接收系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, l.opts.signals...)

	// 处理关闭 goroutine
	eg.Go(func() error {
		for {
			select {
			// 来自上下文的关闭通知
			case <-ctx.Done():
				return ctx.Err()
			// 来自信号关闭的通知
			case <-quit:
				if err := l.Shutdown(); err != nil {
					return err
				}
				return nil
			}
		}
	})

	// 返回错误信息
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Shutdown 关闭监听服务
func (l *listener) Shutdown() error {
	if l.cancel != nil {
		l.cancel()
	}
	return nil
}

func (l *listener) on() error {
	// 解析地址
	address := fmt.Sprintf("%s:%d", l.opts.ip, l.opts.port)
	addr, err := net.ResolveTCPAddr(l.opts.network, address)
	if err != nil {
		return err
	}

	// 监听对应端口
	lis, err := net.ListenTCP(l.opts.network, addr)
	if err != nil {
		return err
	}
	l.lis = lis
	log.Println(fmt.Sprintf("[ LISTENER ] listen on %s:%d", l.opts.ip, l.opts.port))

	// 启用工作池机制
	l.work.UseWorkerPool()

	for {
		// 阻塞等待客户端建立连接
		conn, e := l.lis.AcceptTCP()
		if e != nil {
			// 如果 listener 已关闭
			if errors.Is(e, net.ErrClosed) {
				return nil
			}
			log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s accept err: %e", conn.RemoteAddr().String(), e))
			continue
		}
		log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s established", conn.RemoteAddr().String()))

		// 如果当前连接数量超过最大连接数，则关闭新的连接
		if l.mgr.Len() >= l.opts.conns {
			conn.Close()
			log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s because limit connect is ignored", conn.RemoteAddr().String()))
			continue
		}

		// 开启协程处理当前连接任务
		go newConnection(conn, l.mgr, l.router, l.work, l.opts.packet).Handle()
	}
}

func (l *listener) off() error {
	log.Println(fmt.Sprintf("[ LISTENER ] listener is closeing"))

	// 关闭所有连接
	l.mgr.Clear()

	// 停止监听
	if e := l.lis.Close(); e != nil && !errors.Is(e, net.ErrClosed) {
		return e
	}

	log.Println(fmt.Sprintf("[ LISTENER ] closed"))
	return nil
}
