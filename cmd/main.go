package main

import (
	"context"
	"onlineCLoud/internel/app"
)

func main() {
	err := app.Run(context.Background(), app.SetConfigFile("config/config.bak.toml"),
		app.SetVersion("v1.1"))
	if err != nil {
		panic(err)
	}
}
