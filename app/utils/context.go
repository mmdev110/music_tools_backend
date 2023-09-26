package utils

import (
	"context"
)

type UUID string

var UUID_key UUID = "UUID"

type Email string

var Email_key Email = "Email"

func SetParamsInContext(ctx context.Context, uuid, email string) context.Context {
	ctx1 := context.WithValue(ctx, UUID_key, uuid)
	ctx2 := context.WithValue(ctx1, Email_key, email)
	return ctx2
}

func GetParamsFromContext(ctx context.Context) (uuid, email string) {
	uuid = ctx.Value(UUID_key).(string)
	email = ctx.Value(Email_key).(string)

	return uuid, email
}
