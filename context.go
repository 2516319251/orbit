package orbit

// Context 结构体
type Context struct {
	protocol uint32
	data     []byte
	conn     Connection
}

// RemoteAddr 获取客户端地址
func (ctx *Context) RemoteAddr() string {
	return ctx.conn.RemoteAddr()
}

// Protocol 获取当前服务所属模块
func (ctx *Context) Protocol() uint32 {
	return ctx.protocol
}

// RawData 获取未处理过的请求数据
func (ctx *Context) RawData() []byte {
	return ctx.data
}

// Write 返回数据
func (ctx *Context) Write(b []byte) error {
	return ctx.conn.Send(ctx.protocol, b)
}
