package support

import (
	"com.fha.gocan/internal/platform"
	"com.fha.gocan/internal/platform/config"
	"fmt"
	"io/ioutil"
)

func CreateContext() *context.Context {
	ui := FakeUI{}
	c := config.DefaultConfig
	dir, err := ioutil.TempDir("", "gocan")
	if err != nil {
		fmt.Printf("Cannot create temp directory: %s", err.Error())
		return nil
	}
	c.EmbeddedDataPath = dir
	return context.New(&ui, &c)
}

