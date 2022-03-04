package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"orbit"
	"time"
)

func main() {
	go Client0()
	//go Client1()
	//go Client2()

	r := orbit.InitRouter()
	r.Handle(1, func(ctx *orbit.Context) {
		fmt.Printf("[ SERVER ] receive msg form client: protocol = %d, data = %s\n", ctx.Protocol(), ctx.RawData())
		ctx.Write([]byte("pong"))
	})

	srv := orbit.New(
		orbit.WithNetwork("tcp"),
		orbit.WithIP("127.0.0.1"),
		orbit.WithPort(62817),
		orbit.WithMaxConns(10),
		orbit.WithMaxMessagePacketSize(1024),
		orbit.WithMaxWorkerPoolSize(1),
		orbit.WithMaxWorkerTasksQueueLength(64),
		orbit.WithRouter(r),
		orbit.WithContext(context.Background()),
	)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}

func Client0() {
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:62817")
	if err != nil {
		panic(err)
	}

	read(conn)
	//conn.Close()
}

func Client1() {
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:62817")
	if err != nil {
		panic(err)
	}

	read(conn)
	//conn.Close()
}

func Client2() {
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:62817")
	if err != nil {
		panic(err)
	}

	read(conn)
	conn.Close()
}

func read(conn net.Conn) {
	i := 0
	for {
		dp := orbit.NewDataPacket()
		send, _ := dp.Pack(orbit.NewMessagePacket(1, []byte("ping")))
		_, err := conn.Write(send)
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
		}
		fmt.Printf("[ CLIENT ] receive msg form server: protocol = %d, data = %s\n", receive.GetProtocol(), string(data))

		if i >= 2 {
			break
		}

		time.Sleep(3 * time.Second)
		i++
	}
}
