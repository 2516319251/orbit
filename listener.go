package orbit

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

// Server 监听者接口
type Server interface {
	Run() error
	On() error
	Off() error
}

// listener 监听器结构体
type listener struct {
	opts options
	lis  *net.TCPListener

	mgr    Manager
	router Router
	work   Worker
}

// New 实例化监听器
func New(opts ...Option) Server {
	// 初始化默认配置
	o := options{
		network: "tcp",
		ip:      "0.0.0.0",
		port:    62817,
		pool:    8,
		conns:   512,
		tasks:   1024,
		packet:  4096,
	}

	// 加载自定义配置
	for _, opt := range opts {
		opt(&o)
	}

	if o.router == nil {
		panic("router is nil")
	}

	return &listener{
		opts: o,
		mgr:  &manager{conns: make(map[string]Connection)},
		work: &worker{
			poolSize:  o.pool,
			taskLen:   o.tasks,
			taskQueue: make([]chan *Context, o.pool),
			router:    o.router,
		},
	}
}

// Run 开始监听服务
func (l *listener) Run() error {
	go func() {
		if e := l.On(); e != nil {
			panic(e)
		}
	}()

	// 接收系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, l.opts.signals...)
	<-quit

	return l.Off()
}

func (l *listener) On() error {
	log.Println(fmt.Sprintf("[ LISTENER ] startup"))

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
		go newConnection(conn, l.mgr, l.work, l.opts.packet).Handle()
	}
}

func (l *listener) Off() error {
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
