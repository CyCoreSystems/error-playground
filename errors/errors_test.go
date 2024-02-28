package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/CyCoreSystems/error-playground/errors"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

func TestSpecialErrors(t *testing.T) {
	rootErr := errors.NewSpecial(&errors.Special{
		Name:        "root",
		Description: "root error",
		Code:        int(codes.PermissionDenied),
	})

	midErr := fmt.Errorf("I am some metadata: %w", rootErr)

	outerErr := pkgerrors.Wrap(midErr, "I think I know what ID I have")

	t.Logf("outerErr: %s", outerErr.Error())

	specialError := new(errors.SpecialError)

	if stderrors.As(outerErr, &specialError) {
		t.Logf("wrapped rootErr: %s", specialError.Error())

		t.Logf("special data: %+v", specialError.Special())
	}

	var localError errors.Error

	if stderrors.As(outerErr, &localError) {
		t.Logf("wrapped (B) rootErr: %s", specialError.Error())

		t.Logf("special (B) data: %+v", specialError.Special())
	}
}

func TestHappyErrors(t *testing.T) {
	rootErr := errors.NewHappy()

	midErr := fmt.Errorf("I think I have something joyful: %w", rootErr)

	outerErr := pkgerrors.Wrap(midErr, "if I container goodness, good for you; I do not know")

	t.Logf("outerErr: %s", outerErr.Error())

	t.Logf("special: %+v", errors.Specialize(outerErr))

	var localError errors.Error

	if stderrors.As(outerErr, &localError) {
		t.Logf("unwrapped happines is: %s", localError.InternalID())
	}
}
