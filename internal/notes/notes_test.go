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

func Test_noteType(t *testing.T) {
	tests := []struct {
		name string
		note string
		want string
	}{
		{
			name: "Test 1",
			note: "test",
			want: "noteCount",
		},
		{
			name: "Test 2",
			note: "This is a test note",
			want: "noteCount",
		},
		{
			name: "Test 3",
			note: "[ENC]abc123",
			want: "encNoteCount",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := noteType(tt.note); got != tt.want {
				t.Errorf("noteType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encryptDecrypt(t *testing.T) {
	cipherText, err := encryptNote("test", "44ec4dd68dd0aee62aedf766")
	if err != nil {
		t.Errorf("Encryption failed: %s", err)
	}

	plainText, err := decryptNote(cipherText, "44ec4dd68dd0aee62aedf766")
	if err != nil {
		t.Errorf("Decryption failed: %s", err)
	}

	if plainText != "test" {
		t.Errorf("Decrypted text should be \"test\" but instead got: %s", plainText)
	}
}
