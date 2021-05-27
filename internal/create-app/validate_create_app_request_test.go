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

func TestMaxAppNameLength(t *testing.T) {
	t.Log("\tGiven a create app request")
	{
		t.Log("\tWhen the provided app name is longer than 255 characters")
		{
			request := CreateAppRequest{
				Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				SceneName: "some scene",
			}
			test_support.AssertInvalidField(t, request, "Name")
		}
	}
}

func TestMaxSceneNameLength(t *testing.T) {
	t.Log("\tGiven a create app request")
	{
		t.Log("\tWhen the provided scene name is longer than 255 characters")
		{
			request := CreateAppRequest{
				Name: "an app",
				SceneName: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			}
			test_support.AssertInvalidField(t, request, "SceneName")
		}
	}

}




