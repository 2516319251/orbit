package orbit

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// defaultHeadLength 消息头部长度
const defaultHeadLength = 8

// Packet 数据包接口
type Packet interface {
	GetHeadLength() uint32
	Pack(msg Message) ([]byte, error)
	Unpack(data []byte, maxSize uint32) (Message, error)
}

// packet 数据包结构体
type packet struct{}

// NewDataPacket 数据包实例化
func NewDataPacket() Packet {
	return &packet{}
}

// GetHeadLength 获取包头长度
func (pk *packet) GetHeadLength() uint32 {
	// msg.length 4 字节 + msg.protocol 4 字节
	return defaultHeadLength
}

// Pack 封包
func (pk *packet) Pack(msg Message) ([]byte, error) {
	// 创建缓冲区
	buff := bytes.NewBuffer([]byte{})

	// 写数据长度
	if err := binary.Write(buff, binary.LittleEndian, msg.GetLength()); err != nil {
		return nil, err
	}

	// 写数据协议
	if err := binary.Write(buff, binary.LittleEndian, msg.GetProtocol()); err != nil {
		return nil, err
	}

	// 写数据内容
	if err := binary.Write(buff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// Unpack 拆包
func (pk *packet) Unpack(data []byte, maxSize uint32) (Message, error) {
	// 创建 io reader
	buff := bytes.NewReader(data)

	// 读取消息长度
	var length uint32
	if err := binary.Read(buff, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	// 读取消息协议
	var protocol uint32
	if err := binary.Read(buff, binary.LittleEndian, &protocol); err != nil {
		return nil, err
	}

	// 只解压 head 的消息，获取 protocol 和 length
	msg := NewMessagePacket(protocol, []byte{})
	msg.SetLength(length)

	// 判断数据长度是否超过允许值
	if maxSize > 0 && msg.GetLength() > maxSize {
		return nil, errors.New("received too large message")
	}

	// 通过 head 的长度，后续需要在从 conn 读取一次数据
	return msg, nil
}
