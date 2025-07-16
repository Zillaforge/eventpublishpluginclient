package epp

import (
	"context"

	"github.com/Zillaforge/eventpublishpluginclient/pb"
	"github.com/Zillaforge/eventpublishpluginclient/utility"
	"google.golang.org/protobuf/types/known/emptypb"
)

func getRouter(c client, ctxs ...context.Context) (output *pb.GetRouterResponseList, err error) {
	if c.err != nil {
		return nil, c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	return pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).GetRouter(ctx, &emptypb.Empty{})
}

// GetRouter ...
func GetRouter(ctxs ...context.Context) (output *pb.GetRouterResponseList, err error) {
	return eppclient.GetRouter(ctxs...)
}

// GetRouter ...
func (ph *PoolHandler) GetRouter(ctxs ...context.Context) (output *pb.GetRouterResponseList, err error) {
	c := &client{}
	defer ph.use(c)()
	return getRouter(*c, ctxs...)
}
