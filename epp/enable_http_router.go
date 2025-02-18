package epp

import (
	"context"

	"pegasus-cloud.com/aes/eventpublishpluginclient/pb"
	"pegasus-cloud.com/aes/eventpublishpluginclient/utility"
)

func enableHttpRouter(c client, input *pb.HttpRequestInfo, ctxs ...context.Context) (output *pb.HttpResponseInfo, err error) {
	if c.err != nil {
		return output, c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	output, err = pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).EnableHttpRouter(ctx, input)
	return output, err
}

// EnableHttpRouter ...
func EnableHttpRouter(input *pb.HttpRequestInfo, ctxs ...context.Context) (output *pb.HttpResponseInfo, err error) {
	return eppclient.EnableHttpRouter(input, ctxs...)
}

// EnableHttpRouter ...
func (ph *PoolHandler) EnableHttpRouter(input *pb.HttpRequestInfo, ctxs ...context.Context) (output *pb.HttpResponseInfo, err error) {
	c := &client{}
	defer ph.use(c)()
	return enableHttpRouter(*c, input, ctxs...)
}
