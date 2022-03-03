package orbit

// Message 消息接口
type Message interface {
	GetLength() uint32
	GetProtocol() uint32
	GetData() []byte

	SetLength(length uint32)
	SetProtocol(id uint32)
	SetData(data []byte)
}

// message 消息结构体
type message struct {
	length   uint32
	protocol uint32
	data     []byte
}

// NewMessagePacket 创建消息数据包
func NewMessagePacket(protocol uint32, data []byte) Message {
	return &message{
		length:   uint32(len(data)),
		protocol: protocol,
		data:     data,
	}
}

// GetLength 获取消息长度
func (msg *message) GetLength() uint32 {
	return msg.length
}

// GetProtocol 获取消息协议
func (msg *message) GetProtocol() uint32 {
	return msg.protocol
}

// GetData 获取消息内容
func (msg *message) GetData() []byte {
	return msg.data
}

// SetLength 设置消息长度
func (msg *message) SetLength(length uint32) {
	msg.length = length
}

// SetProtocol 设置消息协议
func (msg *message) SetProtocol(protocol uint32) {
	msg.protocol = protocol
}

// SetData 设置消息内容
func (msg *message) SetData(data []byte) {
	msg.data = data
}
