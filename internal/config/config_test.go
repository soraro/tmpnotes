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
			want: specification{Port: 7000, RedisUrl: "redis://localhost:1234", EnableHsts: false, MaxLength: 1000, UiMaxLength: 512, MaxExpire: 24},
		},
		{
			name: "Test 2",
			args: map[string]string{"TMPNOTES_PORT": "6000", "TMPNOTES_REDIS_URL": "rediss://someserver:1234", "TMPNOTES_ENABLE_HSTS": "true", "TMPNOTES_MAX_LENGTH": "2000", "TMPNOTES_UI_MAX_LENGTH": "600", "TMPNOTES_MAX_EXPIRE": "48"},
			want: specification{Port: 6000, RedisUrl: "rediss://someserver:1234", EnableHsts: true, MaxLength: 2000, UiMaxLength: 600, MaxExpire: 48},
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

func Test_GetConfig_Error(t *testing.T) {
	tests := []struct {
		name string
		args map[string]string
	}{
		{
			name: "Test 1",
			args: map[string]string{"TMPNOTES_PORT": "6000", "TMPNOTES_REDIS_URL": "rediss://someserver:1234", "TMPNOTES_ENABLE_HSTS": "true", "TMPNOTES_MAX_LENGTH": "2000", "TMPNOTES_UI_MAX_LENGTH": "6000", "TMPNOTES_MAX_EXPIRE": "48"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for env, val := range tt.args {
				os.Setenv(env, val)
				defer os.Unsetenv(env)
			}
			err := GetConfig()
			if err.Error() != "UiMaxLength 6000 should not be greater than MaxLength 2000" {
				t.Errorf("UiMaxLength %v should cause an error since it is larger than MaxLength %v", Config.UiMaxLength, Config.MaxLength)
			}
		})
	}
}
