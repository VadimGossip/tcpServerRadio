package app

import (
	"context"
	"github.com/VadimGossip/tcpServerRadio/internal/api/server/radio"
	"github.com/VadimGossip/tcpServerRadio/internal/config"
	"github.com/VadimGossip/tcpServerRadio/internal/domain"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
	name         string
	configDir    string
	appStartedAt time.Time
	cfg          *domain.Config
	tcpServer    *TcpServer
}

func NewApp(name, configDir string, appStartedAt time.Time) *App {
	return &App{
		name:         name,
		configDir:    configDir,
		appStartedAt: appStartedAt,
	}
}

func (app *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.Init(app.configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}
	app.cfg = cfg

	go func() {
		tcpController := radio.NewController()
		app.tcpServer = NewTcpServer(app.cfg.RadioTcpServer.Port, tcpController)
		if err := app.tcpServer.Run(ctx); err != nil {
			logrus.Fatalf("error occured while running tcp server: %s", err.Error())
		}
	}()

	logrus.Infof("%s started", app.name)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()

	logrus.Infof("%s stopped", app.name)
}
