package orbit

import (
	"fmt"
	"log"
)

// HandlerFunc 执行方法
type HandlerFunc func(ctx *Context)

// HandlersChain 执行方法切片
type HandlersChain []HandlerFunc

// Router 路由接口
type Router interface {
	Handle(protocol uint32, handler HandlerFunc)
	do(ctx *Context)
}

// router 路由结构体
type router struct {
	api map[uint32]HandlerFunc
}

// InitRouter 路由初始化
func InitRouter() Router {
	log.Println(fmt.Sprintf("[ ROUTER ] router init"))
	return &router{
		api: make(map[uint32]HandlerFunc),
	}
}

// Handle 添加处理句柄
func (r *router) Handle(protocol uint32, handler HandlerFunc) {
	if _, ok := r.api[protocol]; ok {
		panic(fmt.Sprintf("repeated protocol: %d", protocol))
	}
	r.api[protocol] = handler

	log.Println(fmt.Sprintf("[ ROUTER ] add protocol %d", protocol))
}

// do 执行
func (r *router) do(ctx *Context) {
	handler, ok := r.api[ctx.Protocol()]
	if !ok {
		return
	}
	handler(ctx)
}
