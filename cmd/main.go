package main

import (
	"os"

	"github.com/Corray333/authbot/internal/app"
	"github.com/Corray333/authbot/internal/config"
)

func main() {
	config.MustInit(os.Args[1])
	app.New().Run()
}
