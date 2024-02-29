package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/CyCoreSystems/error-playground/errors"
	pkgerrors "github.com/pkg/errors"
	"github.com/test-go/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestExternalError(t *testing.T) {
	rootErr := errors.External(
		codes.PermissionDenied,
		"root error",
	)

	midErr := fmt.Errorf("I am some metadata: %w", rootErr)

	outerErr := pkgerrors.Wrap(midErr, "I think I know what ID I have")

	t.Logf("outerErr: %s", outerErr.Error())

	var extError errors.ExternalError

	if stderrors.As(outerErr, &extError) {
		t.Logf("wrapped rootErr: %s", extError.Error())

		t.Logf("special data: %+v", extError.GRPCStatus())
	}

	// Add another detail
	outerErr = errors.AddDetails(outerErr, &anypb.Any{
		TypeUrl: "https://dummy.com/bogus-text",
		Value: []byte("I am bogus"),
	})

	var detailedError errors.DetailedError

	if stderrors.As(outerErr, &detailedError) {
		t.Logf("wrapped (B) rootErr: %s", detailedError.Error())

		t.Logf("special (B) data: %+v", detailedError.Details())

		assert.Len(t, detailedError.Details(), 2)
	}
}
