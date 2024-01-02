package main

import (
	"github.com/tty-monkey/auth-server/internal/app"
	"github.com/tty-monkey/auth-server/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.Start(cfg)
}
