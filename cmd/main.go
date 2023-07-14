package main

import (
	"github.com/VadimGossip/tcpServerRadio/internal/app"
	"time"
)

var configDir = "config"

func main() {
	tcpRadio := app.NewApp("Tcp Radio", configDir, time.Now())
	tcpRadio.Run()
}
