package create_app

import (
	"com.fha.gocan/internal/test_support"
	"testing"
)

func TestEmptyNameReturnsError(t *testing.T) {
	t.Log("\tGiven a create app request")
	{
		t.Log("\tWhen the provided name is empty")
		{
			request := CreateAppRequest{
				Name:      "",
				SceneName: "scene name",
			}
			test_support.AssertInvalidField(t, request, "Name")
		}
	}
}

func TestEmptySceneNameReturnsError(t *testing.T) {
	t.Log("\tGiven a create app request")
	{
		t.Log("\tWhen the provided scene name is empty")
		{
			request := CreateAppRequest{
				Name:      "app name",
				SceneName: "",
			}
			test_support.AssertInvalidField(t, request, "SceneName")
		}
	}
}





