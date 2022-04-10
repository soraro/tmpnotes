// This is taken from https://github.com/gtank/cryptopasta

package crypto

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func TestEncryptDecryptGCM(t *testing.T) {
	randomKey := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, randomKey[:])
	if err != nil {
		t.Fatal(err)
	}

	gcmTests := []struct {
		plaintext []byte
		key       *[32]byte
	}{
		{
			plaintext: []byte("Hello, world!"),
			key:       randomKey,
		},
	}

	for _, tt := range gcmTests {
		ciphertext, err := Encrypt(tt.plaintext, tt.key)
		if err != nil {
			t.Fatal(err)
		}

		plaintext, err := Decrypt(ciphertext, tt.key)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(plaintext, tt.plaintext) {
			t.Errorf("plaintexts don't match")
		}

		ciphertext[0] ^= 0xff
		_, err = Decrypt(ciphertext, tt.key)
		if err == nil {
			t.Errorf("gcmOpen should not have worked, but did")
		}
	}
}
