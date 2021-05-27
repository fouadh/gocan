package create_scene

import (
	"com.fha.gocan/internal/test_support"
	"testing"
)

func TestEmptyNameReturnsError(t *testing.T) {
	t.Log("\tGiven a create scene request")
	{
		t.Log("\tWhen the provided name is empty")
		{
			request := CreateSceneRequest{
				Name: "",
			}
			test_support.AssertInvalidField(t, request, "Name")
		}
	}
}

