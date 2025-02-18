package service

import (
	"github.com/magicnana999/im/apiimpl"
	"github.com/magicnana999/im/conf"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start() {
	lis, err := net.Listen("tcp", conf.Global.Service.Addr)
	if err != nil {
		logger.FatalF("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserApiServer(s, apiimpl.InitUserApi())
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
