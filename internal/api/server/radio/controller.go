package radio

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Controller interface {
	HandleConnection(ctx context.Context, conn net.Conn)
}

type controller struct {
}

var _ Controller = (*controller)(nil)

func NewController() *controller {
	return &controller{}
}

func (c *controller) runConnectionWriter(ctx context.Context, conn net.Conn) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msg := fmt.Sprintf("время: %s", time.Now())

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
