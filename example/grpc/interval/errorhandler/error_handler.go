package errorhandler

import (
	"fmt"
	"log/slog"

	"github.com/syntaxfa/quick-connect/protobuf/shared/golang/errdetailspb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGRPCError(err error, logger *slog.Logger) {
	st, ok := status.FromError(err)
	if !ok {
		logger.Error("Non-gRPC error occurred", slog.String("error", err.Error()))
		fmt.Printf("An unexpected error occurred: %v\n", err)

		return
	}

	logger.Warn("gRPC request failed",
		slog.String("code", st.Code().String()),
		slog.String("message", st.Message()))
	fmt.Printf("gRPC Error:\n Code: %s\n message: %s\n", st.Code(), st.Message())

	if st.Code() == codes.InvalidArgument {
		fmt.Println(" Details:")
		foundDetails := false

		for _, detail := range st.Details() {
			switch d := detail.(type) {
			case *errdetailspb.BadRequest:
				foundDetails = true
				fmt.Println(" BadRequest Details:")
				for _, violation := range d.GetFieldViolations() {
					fmt.Printf("  - Field: %s, description: %s\n", violation.GetField(), violation.GetDescription())
				}
			default:
				fmt.Printf(" Unknown Detail Type: %T\n", d)
			}
		}
		if !foundDetails {
			fmt.Println("No specific BadRequest details found")
		}
	}
}
