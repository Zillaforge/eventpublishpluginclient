package epp

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"pegasus-cloud.com/aes/eventpublishpluginclient/pb"
	"pegasus-cloud.com/aes/eventpublishpluginclient/utility"
)

func getVersion(c client, ctxs ...context.Context) (output *pb.GetVersionResponse, err error) {
	if c.err != nil {
		return nil, c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	return pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).GetVersion(ctx, &emptypb.Empty{})
}

// GetVersion ...
func GetVersion(ctxs ...context.Context) (output *pb.GetVersionResponse, err error) {
	return eppclient.GetVersion(ctxs...)
}

// GetVersion ...
func (ph *PoolHandler) GetVersion(ctxs ...context.Context) (output *pb.GetVersionResponse, err error) {
	c := &client{}
	defer ph.use(c)()
	return getVersion(*c, ctxs...)
}
