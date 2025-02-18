package epp

import (
	"context"

	"pegasus-cloud.com/aes/eventpublishpluginclient/pb"
)

var (
	eppclient EPPFuncs
)

func ReplaceGlobals(funcs EPPFuncs) {
	eppclient = funcs
}

type EPPFuncs interface {
	GetName(...context.Context) (*pb.GetNameResponse, error)
	GetVersion(...context.Context) (*pb.GetVersionResponse, error)
	SetConfig(*pb.SetConfigRequest,...context.Context) (error)
	CheckPluginVersion(...context.Context) (*pb.CheckVersionResponse, error)
	InitPlugin(...context.Context) (*pb.InitPluginResponse, error)
	Reconcile(*pb.ReconcileRequest,...context.Context) (error)
	CallGRPCRouter(*pb.RPCRouterRequest, ...context.Context) (*pb.RPCRouterResponse, error)
	EnableHttpRouter(*pb.HttpRequestInfo, ...context.Context) (*pb.HttpResponseInfo, error)
	GetRouter(...context.Context) (*pb.GetRouterResponseList, error)
}
