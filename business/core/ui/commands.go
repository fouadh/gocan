package ui

import (
	app3 "com.fha.gocan/app/api/app"
	churn2 "com.fha.gocan/app/api/churn"
	"com.fha.gocan/app/api/coupling"
	"com.fha.gocan/app/api/revision"
	scene2 "com.fha.gocan/app/api/scene"
	app2 "com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/churn"
	coupling2 "com.fha.gocan/business/core/coupling"
	revision2 "com.fha.gocan/business/core/revision"
	"com.fha.gocan/business/core/scene"
	context "com.fha.gocan/foundation"
	"embed"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
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
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}

			mux := httptreemux.New()
			mux.PathSource = httptreemux.URLPath


			fsys, _ := fs.Sub(app, "dist")
			fs := http.FS(fsys)
			server := http.FileServer(fs)
			mux.GET("/", func(w http.ResponseWriter, r *http.Request, m map[string]string) {
				server.ServeHTTP(w, r)
				return
			})

			mux.GET("/*path", func(w http.ResponseWriter, r *http.Request, m map[string]string) {
				// todo not very nice approach: can we do simpler than opening the file to check if it exists ?
				fullPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
				f, err := fs.Open(fullPath)
				if err != nil {
					r.URL.Path = "/"
				} else {
					f.Close()
				}
				server.ServeHTTP(w, r)
				return
			})

			group := mux.NewGroup("/api")
			sceneCore := scene.NewCore(connection)
			appCore := app2.NewCore(connection)
			revisionCore := revision2.NewCore(connection)
			couplingCore := coupling2.NewCore(connection)
			churnCore := churn.NewCore(connection)

			sceneHandlers := scene2.Handlers{Scene: sceneCore, App: appCore}
			appHandlers := app3.Handlers{App: appCore}
			revisionHandlers := revision.Handlers{Revision: revisionCore}
			couplingHandlers := coupling.Handlers{Coupling: couplingCore, App: appCore}
			churnHandlers := churn2.Handlers{Churn: churnCore}

			group.GET("/scenes",  func(writer http.ResponseWriter, request *http.Request, params map[string]string) {
				err := sceneHandlers.QueryAll(writer, request)
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:id", func(writer http.ResponseWriter, request *http.Request, params map[string]string) {
				err := sceneHandlers.QueryById(writer, request, params)
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps", func(writer http.ResponseWriter, request *http.Request, params map[string]string) {
				err := appHandlers.QueryAll(writer, request, params)
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := appHandlers.QueryById(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}

			})

			group.GET("/scenes/:sceneId/apps/:appId/revisions", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := revisionHandlers.Query(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/coupling-hierarchy", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := couplingHandlers.BuildCouplingHierarchy(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/code-churn", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := churnHandlers.Query(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
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