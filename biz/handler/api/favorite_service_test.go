package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"testing"
)

func TestFavoriteVideo(t *testing.T) {
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
			FavoriteVideo(tt.args.ctx, tt.args.c)
		})
	}
}

func TestGetFavoriteList(t *testing.T) {
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
			GetFavoriteList(tt.args.ctx, tt.args.c)
		})
	}
}
