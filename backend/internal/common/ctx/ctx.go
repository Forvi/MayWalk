package ctx

import "context"

type key string

const keyRequestID key = "requestID"

func RequestID(ctx context.Context) string {
	reqID, _ := ctx.Value(keyRequestID).(string)
	return reqID
}

func SetRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, keyRequestID, reqID)
}
