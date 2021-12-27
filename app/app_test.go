package app_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iomarmochtar/content-plus-totp/app"
	cfg "github.com/iomarmochtar/content-plus-totp/config"
)

var (
	sampleKey              = "somekey"
	sampleEncryptedContent = "zPiTPXkDWf6+pdI1YBoAO6cYDNl5CQJy50ilzgY95xFS7XpahWWd"
	// MJQA2ESNXLOWX3AFOGDXFTTGAFNBPT32
	sampleEncryptedTotpMaster = "neMH0e9aKPKGi/0j2x8JatttjNE3VzZWk9Xjayhb92Lsurq8J1aVYCwAsJhcYGcQ85oR2M+3imp6jn9P"
)

func TestAppGetCombination(t *testing.T) {
	tcs := map[string]struct {
		config            *cfg.Config
		key               string
		combination       string
		expectingErrMsg   string
		getSimulationTime func() time.Time
	}{
		"error decrypting content, not valid b64 value": {
			key: sampleKey,
			config: &cfg.Config{
				ContentEnc: "not the valid one",
			},
			expectingErrMsg: "while decrypting content: illegal base64 data at input byte 3",
		},
		"error decrypting content, not a valid encrypted one": {
			key: sampleKey,
			config: &cfg.Config{
				ContentEnc: "a2FjYW5naXRlbQo=",
			},
			expectingErrMsg: "while decrypting content: not a valid length decoded b64 value (11) with nonce (12), possibly not a valid encrypted content",
		},
		"error decrypting totp master, not valid b64 value": {
			key: sampleKey,
			config: &cfg.Config{
				ContentEnc:    sampleEncryptedContent,
				TotpMasterEnc: "not a valid b64",
			},
			expectingErrMsg: "while decrypting totp master: illegal base64 data at input byte 3",
		},
		"error decrypting totp master, not a valid encrypted one": {
			key: sampleKey,
			config: &cfg.Config{
				ContentEnc:    sampleEncryptedContent,
				TotpMasterEnc: "a2FjYW5naXRlbQo=",
			},
			expectingErrMsg: "while decrypting totp master: not a valid length decoded b64 value (11) with nonce (12), possibly not a valid encrypted content",
		},
		"error generating totp token, not a valid totp secret": {
			key: sampleKey,
			config: &cfg.Config{
				ContentEnc:    sampleEncryptedContent,
				TotpMasterEnc: sampleEncryptedContent,
			},
			expectingErrMsg: "while generating totp code: Decoding of secret as base32 failed.",
		},
		"success": {
			key: sampleKey,
			config: &cfg.Config{
				ContentEnc:    sampleEncryptedContent,
				TotpMasterEnc: sampleEncryptedTotpMaster,
			},
			getSimulationTime: func() time.Time {
				// RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
				sTime, err := time.Parse(time.RFC1123Z, "Mon, 25 Dec 2021 00:00:00 +0700")
				if err != nil {
					panic(err)
				}
				return sTime
			},
			combination: "hello world831173",
		},
	}

	for title, tc := range tcs {
		t.Run(title, func(t *testing.T) {
			tApp := app.New(tc.config)
			var simulateTime time.Time
			if tc.getSimulationTime == nil {
				simulateTime = time.Now()
			} else {
				simulateTime = tc.getSimulationTime()
			}

			result, err := tApp.GetCombination(tc.key, simulateTime)
			if tc.expectingErrMsg != "" {
				assert.Equal(t, "", result)
				require.EqualError(t, err, tc.expectingErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tc.combination, result)
		})
	}
}
