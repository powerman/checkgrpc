package checkgrpc_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/powerman/check"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrGRPCStatus(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt)
	todo := t.TODO()

	grpcErr := status.Error(codes.Unknown, "unknown")
	grpcErrSame := status.Error(codes.Unknown, "unknown")
	grpcErrDiff := status.Error(codes.Internal, "internal")

	// Equal gRPC errors.
	t.Err(grpcErr, grpcErr)
	t.Err(grpcErr, grpcErrSame)
	todo.NotErr(grpcErr, grpcErr)
	todo.NotErr(grpcErr, grpcErrSame)

	// Different gRPC errors.
	todo.Err(grpcErr, grpcErrDiff)
	t.NotErr(grpcErr, grpcErrDiff)

	// gRPC error vs nil.
	todo.Err(grpcErr, nil)
	t.NotErr(grpcErr, nil)

	// gRPC error vs regular error.
	todo.Err(grpcErr, io.EOF)
	t.NotErr(grpcErr, io.EOF)

	// Regular errors are unaffected.
	t.Err(io.EOF, io.EOF)
	todo.Err(io.EOF, io.ErrUnexpectedEOF)
}

func TestNotErrGRPCStatus(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt)
	todo := t.TODO()

	grpcErr := status.Error(codes.Unknown, "unknown")

	// Two equal gRPC errors should not be not-equal.
	todo.NotErr(grpcErr, grpcErr)

	// gRPC vs non-gRPC.
	t.NotErr(grpcErr, io.EOF)

	// gRPC vs nil.
	t.NotErr(grpcErr, nil)
}

func TestErrWrappedGRPC(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt)
	todo := t.TODO()

	grpcErr := status.Error(codes.Unknown, "unknown")
	wrapped := fmt.Errorf("wrapped: %w", grpcErr)
	wrapped2 := fmt.Errorf("double wrapped: %w", fmt.Errorf("wrapped: %w", grpcErr))

	// Wrapped gRPC errors should be detected.
	t.Err(wrapped, grpcErr)
	t.Err(wrapped2, grpcErr)
	todo.NotErr(wrapped, grpcErr)
	todo.NotErr(wrapped2, grpcErr)

	// Different wrapped gRPC errors.
	grpcErrDiff := status.Error(codes.Internal, "internal")
	wrappedDiff := fmt.Errorf("wrapped: %w", grpcErrDiff)
	todo.Err(wrapped, wrappedDiff)
	t.NotErr(wrapped, wrappedDiff)
}

func TestNonGRPCErrors(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt)
	todo := t.TODO()

	// Non-gRPC errors should work normally.
	t.Err(io.EOF, io.EOF)
	todo.Err(io.EOF, io.ErrUnexpectedEOF)
	t.NotErr(io.EOF, io.ErrUnexpectedEOF)
	todo.NotErr(io.EOF, io.EOF)

	// Nil error.
	t.Err(nil, nil)
	todo.NotErr(nil, nil)
	todo.Err(nil, io.EOF)
	t.NotErr(nil, io.EOF)
}
