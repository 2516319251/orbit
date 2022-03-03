package orbit

import (
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestNewDataPacket(t *testing.T) {
	go Client4NewDataPacket()

	lis, err := net.Listen("tcp", "127.0.0.1:11111")
	if err != nil {
		t.Error(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			t.Error(err)
		}

		dp := NewDataPacket()
		for {
			head := make([]byte, dp.GetHeadLength())
			if _, e := io.ReadFull(conn, head); e != nil {
				if errors.Is(e, io.EOF) {
					return
				}
				t.Error(e)
				return
			}

			msg, err := dp.Unpack(head, 4096)
			if err != nil {
				t.Error(err)
			}

			data := make([]byte, msg.GetLength())
			if msg.GetLength() > 0 {
				if _, e := io.ReadFull(conn, data); e != nil {
					if errors.Is(e, io.EOF) {
						return
					}
					t.Error(e)
					return
				}
				msg.SetData(data)
				fmt.Printf("protocol: %d, length: %d, data: %s\n", msg.GetProtocol(), msg.GetLength(), msg.GetData())
			}
		}
	}

	lis.Close()
}

func Client4NewDataPacket() {
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:11111")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	dp := NewDataPacket()
	msg1 := NewMessagePacket(0, []byte{'h', 'e', 'l', 'l', 'o'})
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := NewMessagePacket(1, []byte{'w', 'o', 'r', 'l', 'd', '!', '!'})
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client temp msg2 err:", err)
		return
	}

	conn.Write(append(sendData1, sendData2...))
	conn.Close()
}
