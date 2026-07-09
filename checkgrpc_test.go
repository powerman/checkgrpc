package checkgrpc_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/powerman/check"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// causeError implements Cause() for testing the cause() unwrap path.
type causeError struct {
	err error
}

func (e *causeError) Error() string { return e.err.Error() }
func (e *causeError) Cause() error  { return e.err }

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

func TestErrWithCause(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt)

	grpcErr := status.Error(codes.Unknown, "unknown")
	grpcErrDiff := status.Error(codes.Internal, "internal")

	// Error with Cause() wrapping a gRPC error.
	// This exercises the cause() function's c.Cause() path (line 74).
	t.Err(&causeError{err: grpcErr}, grpcErr)
	t.NotErr(&causeError{err: grpcErr}, grpcErrDiff)
	t.NotErr(&causeError{err: grpcErr}, io.EOF)

	// Error with Cause() wrapping a non-gRPC error.
	t.Err(&causeError{err: io.EOF}, io.EOF)
	t.NotErr(&causeError{err: io.EOF}, io.ErrUnexpectedEOF)

	// Chain of Cause() errors.
	t.Err(&causeError{err: &causeError{err: grpcErr}}, grpcErr)
	t.NotErr(&causeError{err: &causeError{err: grpcErr}}, grpcErrDiff)
}

func TestErrJoinedGRPC(tt *testing.T) {
	tt.Parallel()
	t := check.T(tt)

	grpcErr := status.Error(codes.Unknown, "unknown")
	grpcErrDiff := status.Error(codes.Internal, "internal")

	// errors.Join with multiple errors where first is gRPC.
	// This exercises the interface{ Unwrap() []error } path (lines 53-57).
	joined := errors.Join(grpcErr, io.EOF)
	t.Err(joined, grpcErr)
	t.NotErr(joined, grpcErrDiff)

	// errors.Join with single gRPC error.
	joinedSingle := errors.Join(grpcErr)
	t.Err(joinedSingle, grpcErr)

	// errors.Join with gRPC error wrapped in another Join.
	doubleJoined := errors.Join(errors.Join(grpcErr, io.EOF), io.ErrUnexpectedEOF)
	t.Err(doubleJoined, grpcErr)
	t.NotErr(doubleJoined, grpcErrDiff)

	// Wrap gRPC error in both Cause() and Join for combined coverage.
	combined := errors.Join(&causeError{err: grpcErr}, io.EOF)
	t.Err(combined, grpcErr)
	t.NotErr(combined, grpcErrDiff)
}
