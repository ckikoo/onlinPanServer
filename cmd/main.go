package main

import (
	"context"
	"onlineCLoud/internel/app"
)

func main() {
	err := app.Run(context.Background(), app.SetConfigFile("config/config.toml"),
		app.SetVersion("v1.0"))
	if err != nil {
		panic(err)
	}
}
