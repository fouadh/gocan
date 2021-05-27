package test_support

import (
	"github.com/go-playground/validator"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func AssertInvalidField(t *testing.T, request interface{}, fieldName string) {
	v := validator.New()
	err := v.Struct(request)

	if err == nil {
		t.Fatalf("\t%s There is no invalid field", failed)
	}

	errors := err.(validator.ValidationErrors)

	if len(errors) > 1 {
		t.Fatalf("\t%s More than one field is invalid", failed)
	}

	invalidField := errors[0].Field()
	if invalidField != fieldName {
		t.Fatalf("\t%s The invalid field is %s instead of %s", failed, invalidField, fieldName)
	}

	t.Logf("\t%s Then an error should have be returned", succeed)
}

