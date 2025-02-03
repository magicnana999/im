package grpc

import (
	"flag"
	"github.com/magicnana999/im/common/pb"
	"github.com/magicnana999/im/logger"
	"github.com/magicnana999/im/service/impl"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start() {
	flag.Parse()
	lis, err := net.Listen("tcp", "7540")
	if err != nil {
		logger.FatalF("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(ErrorHandlingInterceptor))
	pb.RegisterUserApiServer(s, &impl.UserAPIImpl{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
