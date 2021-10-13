package ui

import (
	active_set "com.fha.gocan/business/api/active-set"
	app2 "com.fha.gocan/business/api/app"
	"com.fha.gocan/business/api/boundary"
	"com.fha.gocan/business/api/churn"
	"com.fha.gocan/business/api/complexity"
	"com.fha.gocan/business/api/coupling"
	"com.fha.gocan/business/api/developer"
	modus_operandi "com.fha.gocan/business/api/modus-operandi"
	"com.fha.gocan/business/api/revision"
	"com.fha.gocan/business/api/scene"
	context "com.fha.gocan/foundation"
	"embed"
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/jmoiron/sqlx"
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

func NewStartUiCommand(ctx context.Context) *cobra.Command {
	var serverPort string
	var verbose bool

	cmd := cobra.Command{
		Use: "ui",
		Args: cobra.NoArgs,
		Short: "Start the gocan ui",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Ui.SetVerbose(verbose)
			ctx.Ui.Log("Starting the UI...")
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

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
			handlers := createHandlers(connection)
			for path, h := range handlers {
				handler := h
				group.GET(path, func(writer http.ResponseWriter, request *http.Request, params map[string]string) {
					ctx.Ui.Log("Query " + path)
					err := handler(writer, request, params)
					if err != nil {
						ctx.Ui.Failed(err.Error())
						if verbose {
							fmt.Println(err)
						}
						writer.WriteHeader(http.StatusInternalServerError)
					}

				})
			}

			srv := &http.Server{
				Handler:      mux,
				Addr:         "0.0.0.0:" + serverPort,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}
			ctx.Ui.Ok()

			ctx.Ui.Print("Application running on http://0.0.0.0:" + serverPort)
			log.Fatal(srv.ListenAndServe())
			connection.Close()
			return nil
		},
	}

	cmd.Flags().StringVarP(&serverPort, "port", "p", "1233", "Port to use for the server")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}

func createHandlers(connection *sqlx.DB) map[string]func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	sceneHandlers := scene.NewHandlers(connection)
	appHandlers := app2.NewHandlers(connection)
	revisionHandlers := revision.NewHandlers(connection)
	couplingHandlers := coupling.NewHandlers(connection)
	churnHandlers := churn.NewHandlers(connection)
	modusOperandiHandlers := modus_operandi.NewHandlers(connection)
	activeSetHandlers := active_set.NewHandlers(connection)
	developerHandlers := developer.NewHandlers(connection)
	boundaryHandlers := boundary.NewHandlers(connection)
	complexityHandlers := complexity.NewHandlers(connection)

	handlers := make(map[string](func(w http.ResponseWriter, r *http.Request, params map[string]string) error))
	handlers["/scenes"] = sceneHandlers.QueryAll
	handlers["/scenes/:id"] = sceneHandlers.QueryById
	handlers["/scenes/:sceneId/apps"] = appHandlers.QueryAll
	handlers["/scenes/:sceneId/apps/:appId"] = appHandlers.QueryById
	handlers["/scenes/:sceneId/apps/:appId/revisions"] = revisionHandlers.Query
	handlers["/scenes/:sceneId/apps/:appId/hotspots"] = revisionHandlers.QueryHotspots
	handlers["/scenes/:sceneId/apps/:appId/revisions-trends"] = revisionHandlers.QueryRevisionsTrends
	handlers["/scenes/:sceneId/apps/:appId/boundaries"] = boundaryHandlers.QueryByAppId
	handlers["/scenes/:sceneId/apps/:appId/coupling-hierarchy"] = couplingHandlers.BuildCouplingHierarchy
	handlers["/scenes/:sceneId/apps/:appId/code-churn"] = churnHandlers.Query
	handlers["/scenes/:sceneId/apps/:appId/modus-operandi"] = modusOperandiHandlers.Query
	handlers["/scenes/:sceneId/apps/:appId/active-set"] = activeSetHandlers.Query
	handlers["/scenes/:sceneId/apps/:appId/developers"] = developerHandlers.QueryDevelopers
	handlers["/scenes/:sceneId/apps/:appId/knowledge-map"] = developerHandlers.BuildKnowledgeMap
	handlers["/scenes/:sceneId/apps/:appId/complexity-analyses"] = complexityHandlers.QueryAnalyses
	handlers["/scenes/:sceneId/apps/:appId/complexity-analyses/:complexityId"] = complexityHandlers.QueryAnalysisEntriesById
	return handlers
}