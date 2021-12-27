package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	cbLib "github.com/atotto/clipboard"
	"github.com/iomarmochtar/content-plus-totp/config"
	"github.com/iomarmochtar/content-plus-totp/enc"
	"github.com/urfave/cli/v2"
)

func NewCmd() cli.App {
	return cli.App{
		Name:                 "content-plus-totp",
		Usage:                "an easy way to combine static content such as password with totp token",
		Version:              VERSION,
		Compiled:             time.Now(),
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "configuration path",
			},
			&cli.BoolFlag{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate encrypted config content",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "copy-to-clipboard",
				Aliases: []string{"b"},
				Usage:   "copy the combined to clipboard",
				Value:   false,
			},
		},
		Before: func(c *cli.Context) error {
			if !c.Bool("generate") && c.String("config") == "" {
				return fmt.Errorf("you must set configuration path")
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			if c.Bool("generate") {
				return actionGenerate(c)
			}

			return actionGetCombined(c)
		},
	}
}

func readPwdStdin(msg string) (string, error) {
	// send to stderr so the prompt message still appeared for redirecting the result output
	fmt.Fprint(os.Stderr, msg)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePassword)), nil
}

func printToUser(content string, c *cli.Context) error {
	if c.Bool("copy-to-clipboard") {
		if err := cbLib.WriteAll(content); err != nil {
			return err
		}
		fmt.Println("copied to clipboard")
	} else {
		fmt.Print(content)
	}

	return nil
}

func actionGenerate(c *cli.Context) error {
	inputKey, err := readPwdStdin("encryption key =>")
	if err != nil {
		return err
	}

	inputContentValue, err := readPwdStdin("content => ")
	if err != nil {
		return err
	}

	resultContentValue, err := enc.Encrypt(inputKey, inputContentValue)
	if err != nil {
		return err
	}

	inputTotpMasterValue, err := readPwdStdin("totp master => ")
	if err != nil {
		return err
	}

	resultTotpMasterValue, err := enc.Encrypt(inputKey, inputTotpMasterValue)
	if err != nil {
		return err
	}

	cfg := config.Config{ContentEnc: resultContentValue, TotpMasterEnc: resultTotpMasterValue}

	jsonConfig, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return printToUser(string(jsonConfig), c)
}

func actionGetCombined(c *cli.Context) error {
	cfg, err := config.NewByPath(c.String("config"))
	if err != nil {
		return err
	}

	inputKey, err := readPwdStdin("key =>")
	if err != nil {
		return err
	}

	combined, err := New(cfg).GetCombination(inputKey, time.Now())
	if err != nil {
		return err
	}

	return printToUser(combined, c)
}
