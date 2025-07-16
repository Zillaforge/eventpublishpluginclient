package epp

import (
	"context"

	"github.com/Zillaforge/eventpublishpluginclient/pb"
	"github.com/Zillaforge/eventpublishpluginclient/utility"
)

func setConfig(c client, input *pb.SetConfigRequest, ctxs ...context.Context) (err error) {
	if c.err != nil {
		return c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	_, err = pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).SetConfig(ctx, input)
	return err
}

// SetConfig ...
func SetConfig(input *pb.SetConfigRequest, ctxs ...context.Context) (err error) {
	return eppclient.SetConfig(input, ctxs...)
}

// SetConfig ...
func (ph *PoolHandler) SetConfig(input *pb.SetConfigRequest, ctxs ...context.Context) (err error) {
	c := &client{}
	defer ph.use(c)()
	return setConfig(*c, input, ctxs...)
}
