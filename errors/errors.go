package errors

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/nats-io/nuid"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

// StatusUnknown is a grpc Status with unknown details.
var StatusUnknown = status.New(codes.Unknown, "unknown error")

// Detailer defines something from which details may be obtained.
type Detailer interface {
	// Details returns the set of free-form details.
	Details() []*anypb.Any
}

// DetailedError defines an error which contains additional details.
type DetailedError interface {
	Detailer

	// Implement the basic Go errors interface.
	Error() string

   // AddDetails appends details to the DetailedError.
	AddDetails(...*anypb.Any)
}

type ExternalError interface {
	DetailedError

	// GRPCStatus implement the gRPC StatusCode (internal) interface
	// See: https://pkg.go.dev/google.golang.org/grpc/status#FromError
	GRPCStatus() *status.Status

	// InternalID is an internal unique identifier
	InternalID() string
}

// detailedError is an error with Details.
type detailedError struct {
	details []*anypb.Any
}

// Implement the basic Go errors interface
func (de *detailedError) Error() string {
	return "(details attached)"
}

// Details retrieves any attached free-form details.
func (de *detailedError) Details() []*anypb.Any {
	return de.details
}

// AddDetail adds a detail to the set of free-form details which may attached to this error.
func (de *detailedError) AddDetail(detail *anypb.Any) {
	de.details = append(de.details, detail)
}

// ExternalError is an Error which is suitable for display externally.
type externalError struct {
	id   string
	status *spb.Status
}

// External returns a new ExternalError.
func External(code codes.Code, message string, details ...*anypb.Any) ExternalError {
	s := &spb.Status{
		Code: int32(code),
		Message: message,
		Details: details,
	}

	// always attach a stack trace to a special error
	s.Details = append(s.Details, &anypb.Any{
		TypeUrl: "https://dummy.com/stacktrace",
		Value:   debug.Stack(),
	})

	return &externalError{
		id:   nuid.Next(),
		status: s,
	}
}

func (err *externalError) Details() []*anypb.Any {
	return err.status.Details
}

func (err *externalError) AddDetails(details ...*anypb.Any) {
	err.status.Details = append(err.status.Details, details...)
}

// Error implements Error.
func (err *externalError) Error() string {
	return fmt.Sprintf("external error %q: %s (%d)", err.id, err.status.Message, err.status.Code)
}

// GRPCStatus implements Error.
func (err *externalError) GRPCStatus() *status.Status {
	return status.FromProto(err.status)
}

// InternalID implements Error.
func (err *externalError) InternalID() string {
	return err.id
}

// AddDetail attaches 
func AddDetails(err error, details ...*anypb.Any) error {
	var detailedErr DetailedError

	// If the original error implements DetailedError already add the detail and return it directly.
	if errors.As(err, &detailedErr) {
		detailedErr.AddDetails(details...)

		return err
	}

	// Otherwise, we need to construct a new detailed error.
	
	// First, check to see if we can extract details from the received error.
	var detailer Detailer

	if errors.As(err, &detailer) {
		// NB: we _prepend_ the existing details here, to preserve order
		details = append(detailer.Details(), details...)
	}

	// Return the new detailed error, with or without 
	return errors.Join(err, &detailedError{
		details: details,
	})
}
