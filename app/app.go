package app

import (
	"fmt"
	"time"

	"github.com/iomarmochtar/content-plus-totp/config"
	"github.com/iomarmochtar/content-plus-totp/enc"
	"github.com/pkg/errors"
	totpLib "github.com/pquerna/otp/totp"
)

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
	return &App{config: config}
}

func (a App) GetCombination(key string, t time.Time) (string, error) {
	var content, masterTotp, totp string
	var err error

	// get the encrypted values
	if content, err = enc.Decrypt(key, a.config.ContentEnc); err != nil {
		return "", errors.Wrap(err, "while decrypting content")
	}

	if masterTotp, err = enc.Decrypt(key, a.config.TotpMasterEnc); err != nil {
		return "", errors.Wrap(err, "while decrypting totp master")
	}

	// generate totp value
	if totp, err = totpLib.GenerateCode(masterTotp, t); err != nil {
		return "", errors.Wrap(err, "while generating totp code")
	}

	return fmt.Sprintf("%s%s", content, totp), nil
}
