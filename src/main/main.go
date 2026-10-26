package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Lixsp11/inspire-go-sdk/src/cli"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	commandContext := cli.NewCommandContext(os.Stdout, os.Stderr)
	commandContext.Context = ctx
	code := openapi.Run(commandContext, os.Args[1:])
	cancel()
	os.Exit(code)
}
