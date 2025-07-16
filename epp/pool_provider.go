package epp

import (
	"context"
	"fmt"
	"time"

	cnt "github.com/Zillaforge/eventpublishpluginclient/constants"
	"github.com/Zillaforge/eventpublishpluginclient/utility"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var _ EPPFuncs = (*PoolHandler)(nil)

type (
	PoolProvider struct {
		Mode              GRPCMode
		TCPProvider       TCPProvider
		UnixProvider      UnixProvider
		RouteResponseType utility.ResponseType // deprecated
		// Sets the timeout for connection establishment (up to and including HTTP/2 handshaking) for all new connections.
		// If this is not set, the default is 120 seconds.
		// A zero or negative value will result in an immediate timeout.
		Timeout time.Duration
		// Sets the maximum message size in bytes the client can receive.
		// If this is not set, gRPC uses the default 4MB.
		MaxReceiveMessageSize int
		// Sets the maximum message size in bytes the client can send.
		// If this is not set, gRPC uses the default `math.MaxInt32`.
		MaxSendMessageSize int
		// How much data can be batched before doing a write on the wire.
		// The corresponding memory allocation for this buffer will be twice the size to keep syscalls low.
		// The default value for this buffer is 32KB.
		// Zero or negative values will disable the write buffer such that each write
		// will be on underlying connection. Note: A Send call may not directly
		// translate to a write.
		WriteBufferSize int
		// Set the size of read buffer, this determines how much data can be read at most for each read syscall.
		// The default value for this buffer is 32KB.
		// Zero or negative values will disable read buffer for a connection so data framer can access the
		// underlying conn directly.
		ReadBufferSize int
		_              struct{}
	}
	PoolHandler struct {
		host []string
		pool
	}

	pool struct {
		clients chan client
		_       struct{}
	}
)

// Init ...
func Init(provider PoolProvider) (err error) {
	p, err := new(provider)
	if err != nil {
		return err
	}

	if provider.Mode == TCPMode {
		p.host = provider.TCPProvider.Hosts
	}

	ReplaceGlobals(p)

	return nil
}

func (ph *PoolProvider) getTimeout() (timeout time.Duration) {
	timeout = defaultConnectionTimeout
	if ph.Timeout != 0 {
		timeout = ph.Timeout * time.Second
	}
	return timeout
}

func (ph *PoolProvider) getMaxReceiveMessageSize() (maxReceiveMessageSize int) {
	maxReceiveMessageSize = defaultMaxReceiveMessageSize
	if ph.MaxReceiveMessageSize != 0 {
		maxReceiveMessageSize = ph.MaxReceiveMessageSize
	}
	return maxReceiveMessageSize
}

func (ph *PoolProvider) getMaxSendMessageSize() (maxSendMessageSize int) {
	maxSendMessageSize = defaultMaxReceiveMessageSize
	if ph.MaxSendMessageSize != 0 {
		maxSendMessageSize = ph.MaxSendMessageSize
	}
	return maxSendMessageSize
}

func (ph *PoolProvider) getWriteBufferSize() (writeBufferSize int) {
	writeBufferSize = defaultWriteBufSize
	if ph.WriteBufferSize != 0 {
		writeBufferSize = ph.WriteBufferSize
	}
	return writeBufferSize
}

func (ph *PoolProvider) getReadBufferSize() (readBufferSize int) {
	readBufferSize = defaultReadBufSize
	if ph.ReadBufferSize != 0 {
		readBufferSize = ph.ReadBufferSize
	}
	return readBufferSize
}

func new(provider PoolProvider) (poolHandler *PoolHandler, err error) {
	poolHandler = &PoolHandler{}
	utility.RouteResponseType = provider.RouteResponseType // deprecated
	targets, options := []string{}, []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(tracer.NewGRPCUnaryClientInterceptor()),
		grpc.WithWriteBufferSize(provider.getWriteBufferSize()),
		grpc.WithReadBufferSize(provider.getReadBufferSize()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(provider.getMaxReceiveMessageSize()),
			grpc.MaxCallSendMsgSize(provider.getMaxSendMessageSize()),
		),
	}
	switch provider.Mode {
	case UnixMode:
		count := provider.UnixProvider.ConnCount
		poolHandler.clients = make(chan client, count)
		connTerm := fmt.Sprintf("unix://%s", provider.UnixProvider.SocketPath)
		for i := 0; i < count; i++ {
			targets = append(targets, connTerm)
		}
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	case TCPMode:
		fallthrough
	default:
		count := len(provider.TCPProvider.Hosts) * provider.TCPProvider.ConnPerHost
		poolHandler.clients = make(chan client, count)
		for i := 0; i < count; i++ {
			targets = append(targets, provider.TCPProvider.Hosts[i%len(provider.TCPProvider.Hosts)])
		}
		if provider.TCPProvider.TLS.Enable {
			cert, err := credentials.NewClientTLSFromFile(provider.TCPProvider.TLS.CertPath, "")
			if err != nil {
				return nil, err
			}
			options = append(options, grpc.WithTransportCredentials(cert))
		} else {
			options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
	}
	if err := poolHandler.dials(provider.Mode, targets, provider.getTimeout(), options...); err != nil {
		return nil, err
	}
	return poolHandler, nil

}

// Close ...
func (ph *PoolHandler) Close() {
	if err := ph.close(); err != nil {
		panic(err)
	}
}

func (ph *PoolHandler) dials(mode GRPCMode, targets []string, timeout time.Duration, options ...grpc.DialOption) (err error) {
	for _, target := range targets {
		// the connection still be established and store into ph.clients channel If any errors are appeared.
		c, _ := ph.dial(mode, target, timeout, options...)
		ph.clients <- *c
	}
	return nil
}

func (ph *PoolHandler) dial(mode GRPCMode, target string, timeout time.Duration, options ...grpc.DialOption) (cc *client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	c, err := grpc.DialContext(ctx, target, options...)
	defer cancel()
	cc = &client{
		target:  target,
		timeout: timeout,
		conn:    c,
		mode:    mode,
		opts:    options,
	}
	if err != nil {
		cancel()
		return cc, err
	}
	return cc, nil
}

func (ph *PoolHandler) use(c *client) func() {
	var counter int = 1
	total := len(ph.clients)
	for {
		*c = ph.get()
		if ph.check(*c) {
			break
		} else {
			cc, err := ph.dial(c.mode, c.target, c.timeout, c.opts...)
			if err == nil {
				c = cc
				break
			}
		}
		ph.recycle(*c)
		counter++
		// if the all of connection are not ready
		// and value of counter is bigger than total
		// then returns 14000001 error
		if counter > total {
			c.err = tkErr.New(cnt.GRPCAllOfConnectionsCanNotBeUsedErr)
			return func() {}
		}
	}
	return func() {
		ph.recycle(*c)
	}
}

// New ...
func New(provider PoolProvider) (poolHandler *PoolHandler, err error) {
	return new(provider)
}

func (ph *PoolHandler) get() (c client) {
	c = <-ph.clients
	return c
}

func (ph *PoolHandler) check(c client) bool {
	if c.conn == nil {
		return false
	}
	return c.conn.GetState() == connectivity.Ready
}

func (ph *PoolHandler) recycle(c client) {
	ph.clients <- c
}

func (ph *PoolHandler) close() (err error) {
	for i := 0; i < len(ph.clients); i++ {
		client := <-ph.clients
		if err := client.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
