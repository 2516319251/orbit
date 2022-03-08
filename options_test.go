package orbit

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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

type mockRouter struct{}

func (m *mockRouter) Handle(protocol uint32, handler HandlerFunc) {}
func (m *mockRouter) exec(ctx *Context)                         {}

func TestWithRouter(t *testing.T) {
	o := &options{}
	v := &mockRouter{}

	WithRouter(v)(o)
	assert.Equal(t, v, o.router)
}
