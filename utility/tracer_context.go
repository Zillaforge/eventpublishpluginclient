package utility

import (
	"context"

	"github.com/Zillaforge/toolkits/tracer"
	"google.golang.org/grpc/metadata"
)

func GetContext(ctxs ...context.Context) context.Context {
	ctx := context.Background()
	if len(ctxs) != 0 {
		ctx = ctxs[0]
		if ctx.Value(tracer.TracerContext) != nil {
			ctx = ctx.Value(tracer.TracerContext).(context.Context)
		}
		if ctx.Value(tracer.RequestID) != nil {
			ctx = metadata.NewOutgoingContext(
				ctx,
				metadata.Pairs(tracer.RequestID, ctx.Value(tracer.RequestID).(string)),
			)
		}

	}
	return ctx
}
