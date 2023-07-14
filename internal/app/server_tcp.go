package app

import (
	"context"
	"fmt"
	"github.com/VadimGossip/tcpServerRadio/internal/api/server/radio"
	"github.com/sirupsen/logrus"
	"net"
)

type TcpServer struct {
	tcpPort       int
	tcpController radio.Controller
}

func NewTcpServer(tcpPort int, tcpController radio.Controller) *TcpServer {
	return &TcpServer{tcpPort: tcpPort, tcpController: tcpController}
}

func (s *TcpServer) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.tcpPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logrus.Infof("Tcp server started at %d", s.tcpPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Errorf("error occurred while running tcp server: %s", err.Error())
		}

		go s.tcpController.HandleConnection(ctx, conn)
	}
}
