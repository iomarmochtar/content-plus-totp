package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/iomarmochtar/content-plus-totp/config"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	tcs := map[string]struct {
		jsonContent  string
		expectErrMsg string
	}{
		"valid json content": {
			jsonContent: `{"content_enc": "aGVsbG8K", "totp_master_enc": "d29ybGQK"}`,
		},
		"not a valid json content": {
			jsonContent:  `hello world`,
			expectErrMsg: "invalid character 'h' looking for beginning of value",
		},
		"content_enc is not set": {
			jsonContent:  `{"totp_master_enc": "d29ybGQK"}`,
			expectErrMsg: `you must set "content_enc" field`,
		},
		"totp_master_enc is not set": {
			jsonContent:  `{"content_enc": "d29ybGQK"}`,
			expectErrMsg: `you must set "totp_master_enc" field`,
		},
	}

	for title, tc := range tcs {
		t.Run(title, func(t *testing.T) {
			_, err := config.New([]byte(tc.jsonContent))

			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			}
		})
	}
}

func TestReadConfigByPath(t *testing.T) {
	getTmpFile := func() string {
		tempFile, err := os.CreateTemp(os.TempDir(), "test_read_path")
		if err != nil {
			panic(err)
		}
		return tempFile.Name()
	}

	tcs := map[string]struct {
		preExec      func() string
		postExec     func(string)
		expectErrMsg string
	}{
		"valid content": {
			preExec: func() string {
				tmpFile := getTmpFile()
				jsonContent := `{"content_enc": "aGVsbG8K", "totp_master_enc": "d29ybGQK"}`
				if err := os.WriteFile(tmpFile, []byte(jsonContent), 0600); err != nil {
					panic(err)
				}

				return tmpFile
			},
			postExec: func(path string) {
				if err := os.Remove(path); err != nil {
					panic(err)
				}
			},
		},
		"config is not exists": {
			preExec: func() string {
				return "/this/is/not/exists.json"
			},
			expectErrMsg: "stat /this/is/not/exists.json: no such file or directory",
		},
		"the provided path is a directory": {
			preExec:      os.TempDir,
			expectErrMsg: fmt.Sprintf("%s is a directory", os.TempDir()),
		},
	}

	for title, tc := range tcs {
		t.Run(title, func(t *testing.T) {
			path := tc.preExec()
			_, err := config.NewByPath(path)

			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			}

			if tc.postExec != nil {
				tc.postExec(path)
			}
		})
	}
}
