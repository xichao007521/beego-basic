package b_globals

import (
	"context"
)

type requestIdKey string

var (
	reqID = requestIdKey("req-id")
)

func GetRequestID(ctx context.Context) (int, bool) {
	id, exists := ctx.Value(reqID).(int)
	return id, exists
}

func WithRequestID(ctx context.Context, reqId int) context.Context {
	return context.WithValue(ctx, reqID, reqId)
}
