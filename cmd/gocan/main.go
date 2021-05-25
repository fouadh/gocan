package main

import (
	create_scene "com.fha.gocan/internal/create-scene"
	"com.fha.gocan/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gocan",
}

var createCmd = create_scene.BuildCreateSceneCmd()
var uiCmd = ui.BuildUiCommand()

func main() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.Execute()
}