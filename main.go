package main

import (
	"fmt"
	"os"

	"github.com/iomarmochtar/content-plus-totp/app"
)

func main() {
	a := app.NewCmd()
	if err := a.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
