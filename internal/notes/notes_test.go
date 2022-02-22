package notes

import (
	"strings"
	"testing"
)

func Test_checkAcceptableLength(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "Test 1",
			args: strings.Repeat("a", maxLength+1),
			want: false,
		},
		{
			name: "Test 2",
			args: "a short note",
			want: true,
		},
		{
			name: "Test 3",
			args: strings.Repeat("a", maxLength),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkAcceptableLength(tt.args); got != tt.want {
				t.Errorf("checkAcceptableLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_returnData(t *testing.T) {
	type args struct {
		useragent string
		header    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test 1",
			args: args{
				useragent: "curl",
				header:    "",
			},
			want: true,
		},
		{
			name: "Test 2",
			args: args{
				useragent: "Chrome",
				header:    "",
			},
			want: false,
		},
		{
			name: "Test 3",
			args: args{
				useragent: "Mozilla",
				header:    "Destroy",
			},
			want: true,
		},
		{
			name: "Test 4",
			args: args{
				useragent: "wget",
				header:    "Destroy",
			},
			want: true,
		},
		{
			name: "Test 5",
			args: args{
				useragent: "random",
				header:    "other",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnData(tt.args.useragent, tt.args.header); got != tt.want {
				t.Errorf("returnData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_textResponse(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "Test 1",
			args: "curl",
			want: true,
		},
		{
			name: "Test 2",
			args: "other",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := textResponse(tt.args); got != tt.want {
				t.Errorf("textResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
