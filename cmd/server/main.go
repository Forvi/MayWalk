package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Forvi/MayWalk/internal"
)

func main() {
	app, err := internal.InitApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init app, err: %v\n", err)
		return
	}

	if err := app.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
