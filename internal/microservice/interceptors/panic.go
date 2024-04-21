package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"runtime/debug"
)

func PanicRecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if errPanic := recover(); errPanic != nil {
			log.Println(
				"Panic, ",
				"Method: ", info.FullMethod,
				//"Error: ", errPanic.(string),
				"Message: ", string(debug.Stack()),
			)
		}
	}()

	return handler(ctx, req)
}
