package ui

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed dist
var app embed.FS

func BuildUiCommand() *cobra.Command {
	return &cobra.Command{
		Use: "ui",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("run ui")
			r := http.NewServeMux()
			fsys, _ := fs.Sub(app, "dist")
			webapp := http.FileServer(http.FS(fsys))
			r.Handle("/", webapp)
			srv := &http.Server{
				Handler: r,
				Addr: "localhost:1234",
				WriteTimeout: 15 * time.Second,
				ReadTimeout: 15 * time.Second,
			}
			log.Fatal(srv.ListenAndServe())
			return nil
		},
	}
}