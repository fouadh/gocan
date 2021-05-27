package create_app

import (
	"github.com/go-playground/validator"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func TestEmptyNameReturnsError(t *testing.T) {
	t.Log("\tGiven a create app request")
	{
		t.Log("\tWhen the provided name is empty")
		{
			request := CreateAppRequest{
				Name:      "",
				SceneName: "scene name",
			}
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
			if invalidField != "Name" {
				t.Fatalf("\t%s The invalid field is %s instead of Name", failed, invalidField)
			}

			t.Logf("\t%s Then an error should have be returned", succeed)
		}
	}
}



