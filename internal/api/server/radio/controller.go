package radio

import (
	"context"
	"fmt"
	"github.com/VadimGossip/tcpServerRadio/internal/domain"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Controller interface {
	HandleConnection(ctx context.Context, conn net.Conn)
}

type controller struct {
	cfg domain.RadioLogicConfig
}

var _ Controller = (*controller)(nil)

func NewController(cfg domain.RadioLogicConfig) *controller {
	return &controller{cfg: cfg}
}

func (c *controller) runConnectionWriter(ctx context.Context, conn net.Conn) {
	ticker := time.NewTicker(time.Second * time.Duration(c.cfg.ResponseRate))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msg := fmt.Sprintf("Actual time: %s\n", time.Now())
			_, err := conn.Write([]byte(msg))
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"handler": "runConnectionWriter",
					"problem": "write conn",
				}).Error(err)
				return
			}
		}
	}
}

func (c *controller) HandleConnection(ctx context.Context, conn net.Conn) {
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			logrus.Infof("Connection closed with err %s", err)
		}
	}(conn)
	c.runConnectionWriter(ctx, conn)
}
