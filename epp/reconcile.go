package epp

import (
	"context"

	"github.com/Zillaforge/eventpublishpluginclient/pb"
	"github.com/Zillaforge/eventpublishpluginclient/utility"
)

func reconcile(c client, input *pb.ReconcileRequest, ctxs ...context.Context) (err error) {
	if c.err != nil {
		return c.err
	}
	ctx, cancel := context.WithTimeout(utility.GetContext(ctxs...), c.timeout)
	defer cancel()
	_, err = pb.NewEventPublishPluginInterfaceCRUDControllerClient(c.conn).Reconcile(ctx, input)
	return err
}

// Reconcile ...
func Reconcile(input *pb.ReconcileRequest, ctxs ...context.Context) (err error) {
	return eppclient.Reconcile(input, ctxs...)
}

// Reconcile ...
func (ph *PoolHandler) Reconcile(input *pb.ReconcileRequest, ctxs ...context.Context) (err error) {
	c := &client{}
	defer ph.use(c)()
	return reconcile(*c, input, ctxs...)
}
