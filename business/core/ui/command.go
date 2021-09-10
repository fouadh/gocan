package ui

import (
	context "com.fha.gocan/business/platform"
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

func NewCommand(ctx *context.Context) *cobra.Command {
	var serverPort string

	cmd := cobra.Command{
		Use: "ui",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Ui.Say("Starting the UI...")
			r := http.NewServeMux()
			fsys, _ := fs.Sub(app, "dist")
			webapp := http.FileServer(http.FS(fsys))
			r.HandleFunc("/api/scenes", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "{\"test\": 123}")
			})
			r.Handle("/", webapp)
			srv := &http.Server{
				Handler:      r,
				Addr:         "localhost:" + serverPort,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}
			ctx.Ui.Ok()

			ctx.Ui.Say("Application running on http://localhost:" + serverPort)
			log.Fatal(srv.ListenAndServe())
			return nil
		},
	}

	cmd.Flags().StringVarP(&serverPort, "port", "p", "1233", "Port to use for the server")
	return &cmd
}