package epp

import (
	"context"

	"github.com/Zillaforge/eventpublishpluginclient/pb"
	"github.com/Zillaforge/eventpublishpluginclient/utility"
	"google.golang.org/protobuf/types/known/emptypb"
)

func getName(c client, ctxs ...context.Context) (output *pb.GetNameResponse, err error) {
	if c.err != nil {
		return nil, c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	return pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).GetName(ctx, &emptypb.Empty{})
}

// GetName ...
func GetName(ctxs ...context.Context) (output *pb.GetNameResponse, err error) {
	return eppclient.GetName(ctxs...)
}

// GetName ...
func (ph *PoolHandler) GetName(ctxs ...context.Context) (output *pb.GetNameResponse, err error) {
	c := &client{}
	defer ph.use(c)()
	return getName(*c, ctxs...)
}
