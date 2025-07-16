package epp

import (
	"context"

	"github.com/Zillaforge/eventpublishpluginclient/pb"
	"github.com/Zillaforge/eventpublishpluginclient/utility"
	"google.golang.org/protobuf/types/known/emptypb"
)

func initPlugin(c client, ctxs ...context.Context) (output *pb.InitPluginResponse, err error) {
	if c.err != nil {
		return nil, c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	return pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).InitPlugin(ctx, &emptypb.Empty{})
}

// InitPlugin ...
func InitPlugin(ctxs ...context.Context) (output *pb.InitPluginResponse, err error) {
	return eppclient.InitPlugin(ctxs...)
}

// InitPlugin ...
func (ph *PoolHandler) InitPlugin(ctxs ...context.Context) (output *pb.InitPluginResponse, err error) {
	c := &client{}
	defer ph.use(c)()
	return initPlugin(*c, ctxs...)
}
