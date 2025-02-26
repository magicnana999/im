package training

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func initTracer() *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New())

	return tp
}

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	// 获取并打印 traceId
	traceID := oteltrace.SpanFromContext(ctx).SpanContext().TraceID().String()
	fmt.Printf("Received request, traceId: %s\n", traceID)

	// 返回响应
	return &HelloResponse{Message: "Hello, " + req.Name}, nil
}

func startServer() {
	initTracer()

	// 设置 gRPC 服务
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器并启用 OpenTelemetry 拦截器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()), // 使用 OpenTelemetry 中间件
	)

	// 注册服务
	RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	// 启动服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startClient() {
	initTracer()
	// 连接到 gRPC 服务
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建 Greeter 客户端
	client := NewGreeterClient(conn)

	// 创建上下文并启动 Span
	ctx, span := otel.Tracer("client").Start(context.Background(), "SayHello")
	defer span.End()

	fmt.Println(oteltrace.SpanFromContext(ctx).SpanContext().TraceID().String())

	// 调用服务端方法
	resp, err := client.SayHello(ctx, &HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// 打印响应
	fmt.Printf("Greeting: %s\n", resp.Message)
}
