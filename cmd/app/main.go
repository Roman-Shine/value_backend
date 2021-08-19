package main

import "github.com/Roman-Shine/value_backend/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
