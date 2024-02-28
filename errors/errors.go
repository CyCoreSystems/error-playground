package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nats-io/nuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error interface {
	// Implement the basic Go errors interface
	Error() string

	// GRPCStatus implement the gRPC StatusCode (internal) interface
	// See: https://pkg.go.dev/google.golang.org/grpc/status#FromError
	GRPCStatus() *status.Status

	// InternalID is an internal unique identifier
	InternalID() string

	// Special is a special, internal struct
	Special() *Special
}

// Special is a special, internal struct
type Special struct {
	Name string

	Description string

	Code int
}

// SpecialError is an Error which contains a special struct.
type SpecialError struct {
	id   string
	data *Special
}

// NewSpecial returns a new SpecialError.
func NewSpecial(data *Special) *SpecialError {
	return &SpecialError{
		id: nuid.Next(),
		data: data,
	}
}

// Error implements Error.
func (err *SpecialError) Error() string {
	return fmt.Sprintf("special error: %s (%d)", err.data.Name, err.data.Code)
}

// GRPCStatus implements Error.
func (err *SpecialError) GPRCStatus() *status.Status {
	var code codes.Code

	switch err.data.Code {
	case http.StatusOK:
		code = codes.OK
	case http.StatusNotFound:
		code = codes.NotFound
	case http.StatusBadGateway:
		code = codes.Unavailable
	default:
		code = codes.Internal
	}

	return status.New(code, err.data.Name)
}

// InternalID implements Error.
func (err *SpecialError) InternalID() string {
	return err.id
}

// Special implements Error.
func (err *SpecialError) Special() *Special {
	return err.data
}


// NewHappy returns a new HappyError.
func NewHappy() *HappyError {
	return new(HappyError)
}

// HappyError is an Error which contains happiness.
type HappyError struct{}

// Implement the basic Go errors interface
func (happyerror *HappyError) Error() string {
	return "Happy"
}

// GRPCStatus implement the gRPC StatusCode (internal) interface
// See: https://pkg.go.dev/google.golang.org/grpc/status#FromError
func (happyerror *HappyError) GRPCStatus() *status.Status {
	return status.New(codes.OK, "happy")
}

// InternalID is an internal unique identifier
func (happyerror *HappyError) InternalID() string {
	return "∞"
}

// Special is a special, internal struct
func (happyerror *HappyError) Special() *Special {
	return &Special{
		Name:        "Happiness",
		Description: "Happy",
		Code:        http.StatusOK,
	}
}

// Specialize makes any error Special.
func Specialize(err error) *Special {
	if err == nil {
		return nil
	}

	var packageError Error

	if errors.As(err, &packageError) {
		return packageError.Special()
	}

	return &Special{
		Name:        "unhandled error",
		Description: "unhandled error",
		Code:        int(codes.Unknown),
	}
}
