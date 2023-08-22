package utils

import (
	"context"
	"testing"
)

func Test_UIDContext(t *testing.T) {
	tests := []uint{uint(100), uint(900), uint(0)}
	for _, test := range tests {
		want := test
		ctx := SetUIDInContext(context.Background(), test)
		got := GetUidFromContext(ctx)
		if want != got {
			t.Errorf("uid: got %d, want %d", got, want)
		}
	}
}
