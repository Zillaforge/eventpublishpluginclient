package epp

import (
	"context"

	"github.com/Zillaforge/eventpublishpluginclient/pb"
	"github.com/Zillaforge/eventpublishpluginclient/utility"
)

func callGRPCRouter(c client, input *pb.RPCRouterRequest, ctxs ...context.Context) (output *pb.RPCRouterResponse, err error) {
	if c.err != nil {
		return output, c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	output, err = pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).CallGRPCRouter(ctx, input)
	return output, err
}

// CallGRPCRouter ...
func CallGRPCRouter(input *pb.RPCRouterRequest, ctxs ...context.Context) (output *pb.RPCRouterResponse, err error) {
	return eppclient.CallGRPCRouter(input, ctxs...)
}

// CallGRPCRouter ...
func (ph *PoolHandler) CallGRPCRouter(input *pb.RPCRouterRequest, ctxs ...context.Context) (output *pb.RPCRouterResponse, err error) {
	c := &client{}
	defer ph.use(c)()
	return callGRPCRouter(*c, input, ctxs...)
}
