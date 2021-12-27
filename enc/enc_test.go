package enc_test

import (
	"testing"

	"github.com/iomarmochtar/content-plus-totp/enc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	sampleKey = "somekey"
)

func TestEncrypt(t *testing.T) {
	content := "hello world"

	encrypted, err := enc.Encrypt(sampleKey, content)
	assert.Nil(t, err)

	// it can be reversed (decryptable)
	decrypted, err := enc.Decrypt(sampleKey, encrypted)
	assert.Nil(t, err)

	assert.Equal(t, content, decrypted)
}

func TestDecrypt(t *testing.T) {
	sampleEncryptedContent := "zPiTPXkDWf6+pdI1YBoAO6cYDNl5CQJy50ilzgY95xFS7XpahWWd"

	tsc := map[string]struct {
		key              string
		content          string
		decryptedContent string
		expectingErrMsg  string
	}{
		"not a valid base64 content": {
			key:             sampleKey,
			content:         "this is not a valid b64",
			expectingErrMsg: "illegal base64 data at input byte 4",
		},
		"wrong encryption key": {
			key:             "wrong key",
			content:         sampleEncryptedContent,
			expectingErrMsg: "cipher: message authentication failed",
		},
		"success": {
			key:              sampleKey,
			content:          sampleEncryptedContent,
			decryptedContent: "hello world",
		},
	}

	for title, tc := range tsc {
		t.Run(title, func(t *testing.T) {
			decrypted, err := enc.Decrypt(tc.key, tc.content)

			if tc.expectingErrMsg != "" {
				assert.EqualError(t, err, tc.expectingErrMsg)
				assert.Equal(t, "", decrypted)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.decryptedContent, decrypted)
		})
	}
}
