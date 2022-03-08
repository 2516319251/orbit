package orbit

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	go Client4New()

	r := Setup()
	r.Handle(1, func(ctx *Context) {
		fmt.Printf("[ SERVER ] receive msg form client: protocol = %d, len = %d, data = %s\n",  ctx.Protocol(), len(ctx.RawData()), string(ctx.RawData()))
		ctx.Write(ctx.RawData())
	})

	srv := New(WithRouter(r))
	if e := srv.Run(); e != nil {
		t.Error(e)
	}
}

func Client4New() {
	time.Sleep(3*time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:62817")
	if err != nil {
		panic(err)
	}

	i := 0
	for {
		dp := NewDataPacket()
		send, _ := dp.Pack(NewMessagePacket(1, []byte("test for new listener")))
		_, err = conn.Write(send)
		if err != nil {
			panic(err)
		}

		head := make([]byte, dp.GetHeadLength())
		_, err = io.ReadFull(conn, head)
		if err != nil {
			panic(err)
		}

		receive, err := dp.Unpack(head, 4096)
		if err != nil {
			panic(err)
		}

		var data []byte
		if receive.GetLength() > 0 {
			data = make([]byte, receive.GetLength())
			_, e := io.ReadFull(conn, data)
			if e != nil {
				panic(e)
			}

			fmt.Printf("[ CLIENT ] receive msg form server: protocol = %d, len = %d, data = %s\n",  receive.GetProtocol(), receive.GetLength(), string(data))
		}

		if i >= 2 {
			break
		}

		time.Sleep(1*time.Second)
		i++
	}

	conn.Close()
}
