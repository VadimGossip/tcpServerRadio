package app

import (
	"context"
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
	*Factory
	name              string
	configDir         string
	appStartedAt      time.Time
	cfg               *domain.Config
	cPool             cpool.Pool
	metricsHttpServer *HttpServer
	tcpServer         *TcpServer
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

	app.Factory, err = newFactory(app.cfg)
	if err != nil {
		logrus.Fatalf("Fail to create factory %s", err)
	}

	app.cPool = cpool.NewPool(cfg.ConnectionPoolCfg)

	go func() {
		tcpController := router.NewController(app.routerManager, app.protocolService, app.prometheusService)
		app.tcpServer = NewTcpServer(app.cfg.RouterTcpServer.Port, app.cPool, tcpController)
		if err := app.tcpServer.Run(ctx); err != nil {
			logrus.Fatalf("error occured while running tcp server: %s", err.Error())
		}
	}()

	go func() {
		app.metricsHttpServer = NewHttpServer(app.cfg.RouterMetricsHttpServer.Port)
		initMetricsHttpRouter(app)
		if err := app.metricsHttpServer.Run(); err != nil {
			logrus.Fatalf("error occured while running drs metrics http server: %s", err.Error())
		}
	}()

	logrus.Infof("Drs router started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()

	if err := app.storageClient.Disconnect(); err != nil {
		logrus.Infof("grpc client disconnect error %s", err)
	}

	logrus.Infof("Drs router stopped")
}
