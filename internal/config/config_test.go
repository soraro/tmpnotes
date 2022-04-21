package config

import (
	"os"
	"reflect"
	"testing"
)

func Test_GetConfig(t *testing.T) {
	tests := []struct {
		name string
		args map[string]string
		want specification
	}{
		{
			name: "Test 1",
			args: map[string]string{"PORT": "7000", "REDIS_URL": "redis://localhost:1234"},
			want: specification{Port: 7000, RedisUrl: "redis://localhost:1234"},
		},
		{
			name: "Test 2",
			args: map[string]string{"TMPNOTES_PORT": "6000", "TMPNOTES_REDIS_URL": "rediss://someserver:1234"},
			want: specification{Port: 6000, RedisUrl: "rediss://someserver:1234"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for env, val := range tt.args {
				os.Setenv(env, val)
				defer os.Unsetenv(env)
			}
			GetConfig()
			if !reflect.DeepEqual(Config, tt.want) {
				t.Errorf("%s: GetConfig() = %v, want %v", tt.name, Config, tt.want)
			}
		})
	}
}
