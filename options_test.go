package orbit

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithNetwork(t *testing.T) {
	o := &options{}
	v := "tcp"
	WithNetwork(v)(o)
	assert.Equal(t, v, o.network)
}

func TestWithIP(t *testing.T) {
	o := &options{}
	v := "127.0.0.1"
	WithIP(v)(o)
	assert.Equal(t, v, o.ip)
}

func TestWithPort(t *testing.T) {
	o := &options{}
	v := 62817
	WithPort(v)(o)
	assert.Equal(t, v, o.port)
}

func TestWithMaxConns(t *testing.T) {
	o := &options{}
	v := 1
	WithMaxConns(v)(o)
	assert.Equal(t, v, o.conns)
}

func TestWithMaxWorkPoolSize(t *testing.T) {
	o := &options{}
	v := 10
	WithMaxWorkerPoolSize(v)(o)
	assert.Equal(t, v, o.pool)
}

func TestWithMaxWorkerTasksQueueLength(t *testing.T) {
	o := &options{}
	v := 1024
	WithMaxWorkerTasksQueueLength(v)(o)
	assert.Equal(t, v, o.tasks)
}

func TestWithMaxMessagePacketSize(t *testing.T) {
	o := &options{}
	v := uint32(4096)
	WithMaxMessagePacketSize(v)(o)
	assert.Equal(t, v, o.packet)
}

func TestWithContext(t *testing.T) {
	type ctxKey = struct{}
	o := &options{}
	v := context.WithValue(context.TODO(), ctxKey{}, "ctx")
	WithContext(v)(o)
	assert.Equal(t, v, o.ctx)
}

type mockSig struct{}

func (m *mockSig) String() string { return "sig" }
func (m *mockSig) Signal()        {}

func TestWithSignal(t *testing.T) {
	o := &options{}
	v := []os.Signal{
		&mockSig{}, &mockSig{},
	}
	WithSignal(v...)(o)
	assert.Equal(t, v, o.signals)
}

func TestWithTimeout(t *testing.T) {
	o := &options{}
	v := time.Duration(10)
	WithStopTimeout(v)(o)
	assert.Equal(t, v, o.timeout)
}

type mockRouter struct{}

func (m *mockRouter) Handle(protocol uint32, handler HandlerFunc) {}
func (m *mockRouter) do(ctx *Context)                         {}

func TestWithRouter(t *testing.T) {
	o := &options{}
	v := &mockRouter{}

	WithRouter(v)(o)
	assert.Equal(t, v, o.router)
}
