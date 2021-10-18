package storyboard

import (
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strconv"
	"time"
)

func Commands(ctx foundation.Context) []*cobra.Command {
	return []*cobra.Command{
		create(ctx),
	}
}

func create(ctx foundation.Context) *cobra.Command {
	cmd := cobra.Command{
		Use: "storyboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Ui.Print("storyboard...")

			ctxt, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			ctxt, cancel = context.WithTimeout(ctxt, 15*time.Second)
			defer cancel()

			endpoint := "http://0.0.0.0:1233/"
			sceneId := "a3b28fd2-2fe8-11ec-838e-acde48001122"
			appId := "a74b58cc-2fe8-11ec-b713-acde48001122"
			after, err := date.ParseDay("2020-01-04")
			if err != nil {
				return err
			}
			before, err := date.ParseDay("2020-01-10")
			if err != nil {
				return err
			}

			daysInRange := before.Sub(after).Hours() / 24
			var buffers = make([][]byte, int(daysInRange))

			for i := 0; i < int(daysInRange); i++ {
				max := after.AddDate(0, 0, i)
				if err := chromedp.Run(ctxt, tasks(endpoint, sceneId, appId, date.FormatDay(after), date.FormatDay(max), &buffers[i])); err != nil {
					return errors.Wrap(err, "Unable to browse data")
				}
			}

			for i := 0; i < len(buffers); i++ {
				if err := ioutil.WriteFile("screenshot-" + strconv.Itoa(i) + ".jpeg", buffers[i], 0644); err != nil {
					// todo
					fmt.Println(err)
				}
				fmt.Println("wrote screenshot-" + strconv.Itoa(i) + ".jpeg")
			}


			return nil
		},
	}

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
