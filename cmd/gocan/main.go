package main

import (
	create_scene "com.fha.gocan/internal/create-scene"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gocan",
}

var createCmd = create_scene.BuildCreateSceneCmd()

func main() {
	rootCmd.AddCommand(createCmd)
	rootCmd.Execute()
}