package logger

import (
	"context"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// 提取 Trace ID
		propagator := otel.GetTextMapPropagator()
		md, _ := metadata.FromIncomingContext(ctx)
		ctx = propagator.Extract(ctx, metadataCarrier(md))

		// 创建新的 Span（不会被采样）
		ctx, span := Tracer.Start(ctx, info.FullMethod)
		defer span.End()

		// 获取 Trace ID
		traceID := span.SpanContext().TraceID().String()

		if IsDebugEnable() {
			js, _ := protojson.Marshal(req.(proto.Message))
			Debugf("%s grpc server %s input:%s", traceID, info.FullMethod, string(js))
		}
		reply, err := handler(ctx, req)

		if IsDebugEnable() && err == nil {
			js, _ := proto.Marshal(reply.(proto.Message))
			Debugf("%s grpc client %s output:%s", traceID, info.FullMethod, string(js))
		}

		if IsDebugEnable() && err != nil {
			Debugf("%s grpc client %s error:%s", traceID, info.FullMethod, err.Error())
		}

		return reply, err
	}
}

type metadataCarrier metadata.MD

func (m metadataCarrier) Get(key string) string {
	vals := metadata.MD(m).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (m metadataCarrier) Set(key, value string) {
	metadata.MD(m).Set(key, value)
}

func (m metadataCarrier) Keys() []string {
	keys := make([]string, 0, metadata.MD(m).Len())
	for k := range metadata.MD(m) {
		keys = append(keys, k)
	}
	return keys
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		ctx, span := Tracer.Start(ctx, method)
		defer span.End()

		traceID := span.SpanContext().TraceID().String()

		if IsDebugEnable() {
			js, _ := protojson.Marshal(req.(proto.Message))
			Debugf("%s grpc client %s input:%s", traceID, method, string(js))
		}

		propagator := otel.GetTextMapPropagator()
		md := metadata.New(nil)
		propagator.Inject(ctx, metadataCarrier(md))

		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, method, req, reply, cc, opts...)

		if IsDebugEnable() && err == nil {
			js, _ := protojson.Marshal(reply.(proto.Message))
			Debugf("%s grpc client %s output:%s", traceID, method, string(js))
		}

		if IsDebugEnable() && err != nil {
			Debugf("%s grpc client %s error:%s", traceID, method, err.Error())
		}

		return err
	}
}
