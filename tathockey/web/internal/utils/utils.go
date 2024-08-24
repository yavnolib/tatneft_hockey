package utils

import "context"

type contextKey string

const RequestIDKey contextKey = "requestID"

func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return "-"
	}
	return requestID
}
