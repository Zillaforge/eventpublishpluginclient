package epp

import (
	"math"
	"time"

	"google.golang.org/grpc"
)

const (
	// default value from grpc server
	defaultMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultMaxSendMessageSize    = math.MaxInt32
	defaultConnectionTimeout     = 120 * time.Second
	defaultWriteBufSize          = 32 * 1024
	defaultReadBufSize           = 32 * 1024
)

const (
	TCPMode GRPCMode = iota
	UnixMode
)

type (
	// GRPCMode ...
	GRPCMode int

	// TCPProvider ...
	TCPProvider struct {
		// 當使用 ConnProvider時，Hosts不作用
		Hosts []string
		// 當使用 ConnProvider時，請使用Host
		Host        string
		TLS         TLSConfig
		ConnPerHost int
		_           struct{}
	}

	// TLSConfig ...
	TLSConfig struct {
		Enable   bool
		CertPath string
	}

	// UnixProvider ...
	UnixProvider struct {
		SocketPath string
		ConnCount  int
		_          struct{}
	}

	// Client ...
	client struct {
		target  string
		conn    *grpc.ClientConn
		timeout time.Duration
		err     error
		mode    GRPCMode
		opts    []grpc.DialOption
		_       struct{}
	}
)
