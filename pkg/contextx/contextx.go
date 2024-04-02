package contextx

import (
	"context"
	"onlineCLoud/pkg/util/uuid"
)

type (
	userIDCtx    struct{}
	userEmailCtx struct{}
	uuidCtx      struct{}
)

func NewUUID(ctx context.Context) context.Context {
	return context.WithValue(ctx, uuidCtx{}, uuid.MustString())
}

func FromUUID(ctx context.Context) string {
	v := ctx.Value(uuidCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

func FromUserID(ctx context.Context) string {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewUserEmail(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, userEmailCtx{}, userName)
}

func FromUserEmail(ctx context.Context) string {
	v := ctx.Value(userEmailCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
