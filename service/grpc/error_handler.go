package grpc

import (
	"context"
	"github.com/magicnana999/im/common/pb"
	"github.com/magicnana999/im/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ErrorHandlingInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	// 处理请求并获取响应
	resp, err := handler(ctx, req)
	if err != nil {
		// 错误发生时，记录日志并返回 gRPC 错误
		st, ok := status.FromError(err)
		if !ok {
			// 如果是自定义错误类型，处理自定义错误
			appErr := pb.FromError(err)
			if appErr != nil {
				logger.InfoF("Error occurred: %v", appErr.Error())
				return nil, status.Errorf(appErr.Code, appErr.Message)
			}
		} else {
			// 记录 gRPC 错误
			logger.InfoF("gRPC Error: %v", st.Message())
		}
	}

	// 返回正常响应
	return resp, err
}
