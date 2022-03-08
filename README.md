# orbit

在学习 [zinx](https://github.com/aceld/zinx) 的时候搞出来的。


## Quick start

```go
package main

import "github.com/2516319251/orbit"

func main() {
	r := orbit.Setup()
	r.Handle(1, func(ctx *orbit.Context) {
		fmt.Printf("[ SERVER ] receive msg form client: protocol = %d, data = %s\n", ctx.Protocol(), ctx.RawData())
		ctx.Write([]byte("pong"))
	})

	srv := orbit.New(
		orbit.WithNetwork("tcp"),
		orbit.WithIP("127.0.0.1"),
		orbit.WithPort(62817),
		orbit.WithMaxConns(512),
		orbit.WithMaxMessagePacketSize(4096),
		orbit.WithMaxWorkerPoolSize(8),
		orbit.WithMaxWorkerTasksQueueLength(1024),
		orbit.WithRouter(r),
	)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
```