package util

import "testing"

func TestBcryptHash(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BcryptHash(tt.args.password); got != tt.want {
				t.Errorf("BcryptHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBcryptCheck(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BcryptCheck(tt.args.password, tt.args.hash); got != tt.want {
				t.Errorf("BcryptCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBcryptHash1(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BcryptHash(tt.args.password); got != tt.want {
				t.Errorf("BcryptHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
