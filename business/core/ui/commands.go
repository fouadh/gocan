package ui

import (
	"com.fha.gocan/business/core/scene"
	scene2 "com.fha.gocan/business/data/store/scene"
	context "com.fha.gocan/foundation"
	"embed"
	"encoding/json"
	"expvar"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

//go:embed dist
var app embed.FS

func NewStartUiCommand(ctx *context.Context) *cobra.Command {
	var serverPort string

	cmd := cobra.Command{
		Use: "ui",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Ui.Say("Starting the UI...")
			mux := CreateServeMux()


			fsys, _ := fs.Sub(app, "dist")
			webapp := http.FileServer(http.FS(fsys))
			mux.Handle("/", webapp)

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			sceneCore := scene.NewCore(connection)

			mux.HandleFunc("/api/scenes", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				scenes, err := sceneCore.QueryAll()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				list := struct {
					Scenes []scene2.Scene `json:"scenes"`
				}{
					Scenes: scenes,
				}
				payload, err := json.Marshal(list)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(payload)
			})


			srv := &http.Server{
				Handler:      mux,
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

func CreateServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}