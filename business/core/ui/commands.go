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
					fmt.Println(err)
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

			group.GET("/scenes/:sceneId/apps/:appId/hotspots", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := revisionHandlers.QueryHotspots(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/revisions-trends", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := revisionHandlers.QueryRevisionsTrends(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println(err)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/boundaries", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := boundaryHandlers.QueryByAppId(w, r, params)
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

			group.GET("/scenes/:sceneId/apps/:appId/modus-operandi", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := modusOperandiHandlers.Query(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/active-set", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := activeSetHandlers.Query(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/developers", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := developerHandlers.QueryDevelopers(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/knowledge-map", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := developerHandlers.BuildKnowledgeMap(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/complexity-analyses", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := complexityHandlers.QueryAnalyses(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			group.GET("/scenes/:sceneId/apps/:appId/complexity-analyses/:complexityId", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
				err := complexityHandlers.QueryAnalysisEntriesById(w, r, params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})


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