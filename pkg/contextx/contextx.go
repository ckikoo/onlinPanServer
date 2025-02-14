package contextx

import (
	"context"
	"onlineCLoud/pkg/util/uuid"
)

type (
	userIDCtx    struct{}
	userEmailCtx struct{}
	uuidCtx      struct{}
	middileCtx   struct{}
	AdminCtx     struct{}
)

func NewMiddle(ctx context.Context, reason string) context.Context {
	return context.WithValue(ctx, middileCtx{}, reason)
}

func FromMiddle(ctx context.Context) string {
	v := ctx.Value(middileCtx{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
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
func NewAdmin(ctx context.Context, admin string) context.Context {
	return context.WithValue(ctx, AdminCtx{}, admin)
}

func GetAdmin(ctx context.Context) string {
	v := ctx.Value(AdminCtx{})
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
