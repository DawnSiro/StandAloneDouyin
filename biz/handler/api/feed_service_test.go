package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"testing"
)

func TestGetFeed(t *testing.T) {
	type args struct {
		ctx context.Context
		c   *app.RequestContext
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetFeed(tt.args.ctx, tt.args.c)
		})
	}
}
