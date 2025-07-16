package epp

import (
	"context"

	"github.com/Zillaforge/eventpublishpluginclient/pb"
	"github.com/Zillaforge/eventpublishpluginclient/utility"
	"google.golang.org/protobuf/types/known/emptypb"
)

func checkPluginVersion(c client, ctxs ...context.Context) (output *pb.CheckVersionResponse, err error) {
	if c.err != nil {
		return nil, c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	return pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).CheckPluginVersion(ctx, &emptypb.Empty{})
}

// CheckPluginVersion ...
func CheckPluginVersion(ctxs ...context.Context) (output *pb.CheckVersionResponse, err error) {
	return eppclient.CheckPluginVersion(ctxs...)
}

// CheckPluginVersion ...
func (ph *PoolHandler) CheckPluginVersion(ctxs ...context.Context) (output *pb.CheckVersionResponse, err error) {
	c := &client{}
	defer ph.use(c)()
	return checkPluginVersion(*c, ctxs...)
}
