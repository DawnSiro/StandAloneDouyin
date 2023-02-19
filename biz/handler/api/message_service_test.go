package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"testing"
)

func TestGetMessageChat(t *testing.T) {
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
			GetMessageChat(tt.args.ctx, tt.args.c)
		})
	}
}

func TestSendMessage(t *testing.T) {
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
			SendMessage(tt.args.ctx, tt.args.c)
		})
	}
}
