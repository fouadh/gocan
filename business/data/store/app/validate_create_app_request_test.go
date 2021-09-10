package app

import (
	"com.fha.gocan/business/core/test_support"
	"testing"
)

func TestEmptyNameReturnsError(t *testing.T) {
	t.Log("\tGiven a create app request")
	{
		t.Log("\tWhen the provided name is empty")
		{
			request := NewApp{
				Name:      "",
				SceneId: "sceneId",
			}
			test_support.AssertInvalidField(t, request, "Name")
		}
	}
}

func TestEmptySceneIdReturnsError(t *testing.T) {
	t.Log("\tGiven a create app request")
	{
		t.Log("\tWhen the provided scene name is empty")
		{
			request := NewApp{
				Name:      "app name",
				SceneId: "",
			}
			test_support.AssertInvalidField(t, request, "SceneId")
		}
	}
}

func TestMaxAppNameLength(t *testing.T) {
	t.Log("\tGiven a create app request")
	{
		t.Log("\tWhen the provided app name is longer than 255 characters")
		{
			request := NewApp{
				Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				SceneId: "sceneId",
			}
			test_support.AssertInvalidField(t, request, "Name")
		}
	}
}



