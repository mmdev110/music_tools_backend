package utils

import (
	"context"
)

type UID string

var UID_key UID = "UserID"

func SetUIDInContext(ctx context.Context, value uint) context.Context {
	return context.WithValue(ctx, UID_key, value) //UID_keyを直接"UserID"にするとwarning出ます
}

func GetUidFromContext(ctx context.Context) uint {
	userId := ctx.Value(UID_key).(uint)

	return userId
}
