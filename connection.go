package orbit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// Connection 连接接口
type Connection interface {
	Handle()
	Close()
	Send(protocol uint32, data []byte) error
	RemoteAddr() string
}

// connection 连接结构体
type connection struct {
	conn    *net.TCPConn
	manager Manager
	router  Router
	worker  Worker

	size  uint32
	msgCh chan []byte

	ctx    context.Context
	cancel context.CancelFunc

	close bool
}

// newConnection 创建连接
func newConnection(conn *net.TCPConn, manager Manager, router Router, worker Worker, size uint32) Connection {
	c := &connection{
		conn:    conn,
		manager: manager,
		router:  router,
		worker:  worker,

		size:  size,
		msgCh: make(chan []byte, 1024),

		close: false,
	}

	c.manager.Add(c)

	return c
}

// Handle 处理连接
func (c *connection) Handle() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	// 开启读取客户端数据流的 Goroutine
	go c.readProcessor()
	// 开启返回数据给客户端的 Goroutine
	go c.writeProcessor()

	// 阻塞等待上下文的取消信号
	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

// readProcessor 读处理器
func (c *connection) readProcessor() {
	log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s reader goroutine is running", c.RemoteAddr()))
	defer log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s reader exit", c.RemoteAddr()))
	defer c.Close()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 读取客户端消息的 head
			dp := NewDataPacket()
			head := make([]byte, dp.GetHeadLength())
			if _, err := io.ReadFull(c.conn, head); err != nil {
				log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s read msg head err: %e", c.RemoteAddr(), err))
				return
			}

			// 拆包，获取消息 id 和长度
			msg, err := dp.Unpack(head, c.size)
			if err != nil {
				log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s unpack msg err: %e", c.RemoteAddr(), err))
				return
			}

			// 根据消息长度读取 data
			var data []byte
			if msg.GetLength() > 0 {
				data = make([]byte, msg.GetLength())
				if _, e := io.ReadFull(c.conn, data); e != nil {
					log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s read data head err: %e", c.RemoteAddr(), err))
					return
				}
			}
			msg.SetData(data)

			// 初始化上下文信息
			ctx := &Context{
				protocol: msg.GetProtocol(),
				data:     msg.GetData(),
				conn:     c,
			}

			if c.worker.GetWorkerPoolSize() > 0 {
				// 将消息交给工作池的任务队列中进行处理处理
				c.worker.JoinTaskQueue(ctx)
			} else {
				// 使用路由处理消息
				c.router.do(ctx)
			}
		}
	}
}

// writeProcessor 写处理器
func (c *connection) writeProcessor() {
	log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s writer goroutine is running", c.RemoteAddr()))
	defer log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s writer exit", c.RemoteAddr()))

	for {
		select {
		case <-c.ctx.Done():
			return
		case data, ok := <-c.msgCh:
			if !ok {
				log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s msg buff chan is closed", c.RemoteAddr()))
				return
			}
			if _, err := c.conn.Write(data); err != nil {
				log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s write buff err: %e", c.RemoteAddr(), err))
				return
			}
		}
	}
}

func (c *connection) Send(protocol uint32, data []byte) error {
	if c.close {
		return errors.New("connection closed when send buff msg")
	}

	// 将数据封包
	dp := NewDataPacket()
	msg, err := dp.Pack(NewMessagePacket(protocol, data))
	if err != nil {
		log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s pack msg err: %e", c.RemoteAddr(), err))
		return errors.New(fmt.Sprintf("pack error msg"))
	}

	// 如果管道关闭做超时处理
	timeout := time.NewTimer(5 * time.Millisecond)
	defer timeout.Stop()
	select {
	case <-timeout.C:
		return errors.New("send buff msg timeout")
	case c.msgCh <- msg:
		return nil
	}
}

// Close 关闭连接
func (c *connection) Close() {
	c.cancel()
}

// finalizer 连接关闭后的处理
func (c *connection) finalizer() {
	if c.close {
		return
	}

	log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s connection is closeing", c.RemoteAddr()))

	c.conn.Close()

	c.manager.Del(c)

	close(c.msgCh)

	log.Println(fmt.Sprintf("[ CONNECT ] remote addr %s connection closed", c.RemoteAddr()))

	c.close = true
}

// RemoteAddr 获取远程客户端地址
func (c *connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}
