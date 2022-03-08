package orbit

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestNewConnection(t *testing.T) {
	go Client4TestNewConnection()

	r := Setup()
	r.Handle(1, func(ctx *Context) {
		fmt.Printf("[ SERVER ] receive msg form client: id=%d, len=%d, data=%s\n", ctx.Protocol(), len(ctx.RawData()), ctx.RawData())
		ctx.Write([]byte("Server Testing NewConnection Function..."))
	})

	lis := New(
		WithNetwork("tcp4"),
		WithIP("127.0.0.1"),
		WithPort(11111),
		WithRouter(r),
	)

	go Close4TestNewConnection(lis)

	if err := lis.On(); err != nil && err != context.Canceled {
		t.Error(err)
	}
}

func Close4TestNewConnection(lis Server) {
	time.Sleep(10*time.Second)

	fmt.Println("close listener...")
	if e := lis.Off(); e != nil {
		fmt.Println(e)
	}
}

func Client4TestNewConnection() {
	time.Sleep(3*time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:11111")
	if err != nil {
		fmt.Println("client testing NewConnection function start err:", err)
		return
	}

	for {
		dp := NewDataPacket()
		send, _ := dp.Pack(NewMessagePacket(1, []byte("Client Testing NewConnection Function...")))
		_, err = conn.Write(send)
		if err !=nil {
			fmt.Println("client testing NewConnection function write err:", err)
			return
		}

		head := make([]byte, dp.GetHeadLength())
		_, err = io.ReadFull(conn, head)
		if err != nil {
			fmt.Println("client read head error")
			break
		}

		receive, err := dp.Unpack(head, 4096)
		if err != nil {
			fmt.Println("client unpack err:", err)
			return
		}

		var data []byte
		if receive.GetLength() > 0 {
			data = make([]byte, receive.GetLength())
			_, err = io.ReadFull(conn, data)
			if err != nil {
				fmt.Println("client unpack data err:", err)
				return
			}

			fmt.Println("[ CLIENT ] receive msg form server: id=", receive.GetProtocol(), ", len=", receive.GetLength(), ", data=", string(data))
		}

		time.Sleep(2*time.Second)
		conn.Close()

		return
	}
}