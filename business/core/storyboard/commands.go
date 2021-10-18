package storyboard

import (
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		create(ctx),
	}
}

func create(ctx foundation.Context) *cobra.Command {
	var endpoint string
	var before string
	var after string
	var verbose bool
	var sceneName string

	cmd := cobra.Command{
		Use: "storyboard",
		Args: cobra.ExactArgs(1),
		Short: "Create a storyboard of visualizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)
			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()

			ctxt, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			ctxt, cancel = context.WithTimeout(ctxt, 15*time.Second)
			defer cancel()

			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			daysInRange := beforeTime.Sub(afterTime).Hours() / 24
			var buffers = make([][]byte, int(daysInRange))

			for i := 0; i < int(daysInRange); i++ {
				max := afterTime.AddDate(0, 0, i)
				ui.Log("Getting data between " + date.FormatDay(afterTime) + " and " + date.FormatDay(max))
				if err := chromedp.Run(ctxt, tasks(endpoint, a.SceneId, a.Id, date.FormatDay(afterTime), date.FormatDay(max), &buffers[i])); err != nil {
					return errors.Wrap(err, "Unable to browse data")
				}
			}

			dir, err := ioutil.TempDir("", "gocan-storyboard-")
			defer os.RemoveAll(dir)

			if err != nil {
				return errors.Wrap(err, "Unable to build temp folder")
			}
			for i := 0; i < len(buffers); i++ {
				filename := dir + "/screenshot-" + strconv.Itoa(i) + ".jpeg"
				if err := ioutil.WriteFile(filename, buffers[i], 0644); err != nil {
					// todo
					fmt.Println(err)
					ui.Failed(err.Error())
				}
				ui.Log("wrote " + filename)
			}

			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:1233/", "Endpoint of the UI")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch coupling before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch coupling after this day")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}

func tasks(endpoint string, sceneId string, appId string, min string, max string, buf *[]byte) chromedp.Tasks {
	url := endpoint + `scenes/`+ sceneId + `/apps/` + appId + `?after=` + min + `&before=` + max

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(".Chart", chromedp.ByQuery),
		chromedp.Screenshot(".Chart", buf, chromedp.NodeVisible),
	}

	return tasks
}
