package utils

import (
	"context"
	"testing"
)

func Test_UUIDContext(t *testing.T) {
	type ParamPair struct {
		UUID  string
		Email string
	}
	tests := []ParamPair{
		{"uuid1", "uuid1@gmail.com"},
		{"uuid100", "uuid100@gmail.com"},
		{"uuid999", "uuid999@gmail.com"},
		{"uuid0", "uuid0@gmail.com"},
	}
	for _, test := range tests {
		want := test
		ctx := SetParamsInContext(context.Background(), test.UUID, test.Email)
		got_uuid, got_email := GetParamsFromContext(ctx)
		if want.UUID != got_uuid {
			t.Errorf("UUID: got %s, want %s", got_uuid, want.UUID)
		}
		if want.Email != got_email {
			t.Errorf("Email: got %s, want %s", got_email, want.Email)
		}
	}
}
