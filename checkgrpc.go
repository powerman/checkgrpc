// Package checkgrpc provides gRPC status error comparison for github.com/powerman/check.
//
// Import this package as a blank import in your test file or TestMain
// to enable gRPC status comparison via check.Err/NotErr:
//
//	import _ "github.com/powerman/checkgrpc"
//
// It also imports github.com/powerman/checkproto,
// so a single blank import covers both protobuf message and gRPC status comparison.
package checkgrpc

import (
	"github.com/powerman/check"
	_ "github.com/powerman/checkproto" // Enable proto.Equal via DeepEqual.
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

//nolint:gochecknoinits // Required to register ErrChecker.
func init() {
	check.RegisterErrChecker(CheckGRPCStatus)
}

// CheckGRPCStatus compares two errors using gRPC status comparison.
// It matches the built-in behavior:
// the actual error is unwrapped to its root cause (as [check.Err] does internally),
// and if either the unwrapped actual or the expected error implements GRPCStatus,
// the comparison uses [proto.Equal] on the converted status protos.
func CheckGRPCStatus(actual, expected error) (equal, ok bool) {
	actual2 := unwrapErr(actual)
	_, grpc1 := actual2.(interface{ GRPCStatus() *status.Status })
	_, grpc2 := expected.(interface{ GRPCStatus() *status.Status })
	if !grpc1 && !grpc2 {
		return false, false
	}
	return proto.Equal(
		status.Convert(actual2).Proto(),
		status.Convert(expected).Proto(),
	), true
}

// unwrapErr recursively unwraps err using Cause and Unwrap,
// replicating the same logic as check.unwrapErr.
func unwrapErr(err error) (actual error) {
	defer func() { _ = recover() }()
	actual = err
	for {
		actual = cause(actual)
		var unwrapped error
		switch wrapped := actual.(type) { //nolint:errorlint // False positive.
		case interface{ Unwrap() error }:
			unwrapped = wrapped.Unwrap()
		case interface{ Unwrap() []error }:
			unwrappeds := wrapped.Unwrap()
			if len(unwrappeds) > 0 {
				unwrapped = unwrappeds[0]
			}
		}
		if unwrapped == nil {
			break
		}
		actual = unwrapped
	}
	return actual
}

// cause walks the Cause chain (for github.com/pkg/errors compatibility).
func cause(err error) error {
	for err != nil {
		c, ok := err.(interface{ Cause() error })
		if !ok {
			break
		}
		err = c.Cause()
	}
	return err
}
