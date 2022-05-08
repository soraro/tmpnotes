package config

import "testing"

func TestGetTemplates(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		path    string
	}{
		{
			name:    "Test 1",
			wantErr: false,
			path:    "../../templates/*",
		},
		{
			name:    "Test 2",
			wantErr: true,
			path:    "../../test_data/bad_templates/*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path = tt.path
			if err := GetTemplates(); (err != nil) != tt.wantErr {
				t.Errorf("GetTemplates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
