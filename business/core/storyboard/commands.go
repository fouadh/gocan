package storyboard

import (
	"bytes"
	"com.fha.gocan/business/core"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"context"
	"github.com/chromedp/chromedp"
	"github.com/icza/mjpeg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
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
	var filename string
	var fps int

	cmd := cobra.Command{
		Use:   "storyboard",
		Args:  cobra.ExactArgs(1),
		Short: "Create a storyboard of visualizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.SetVerbose(verbose)

			ctxt, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			ctxt, cancel = context.WithTimeout(ctxt, 15*time.Second)
			defer cancel()

			connection, err := ctx.GetConnection()
			if err != nil {
				return err
			}
			defer connection.Close()
			a, beforeTime, afterTime, err := core.ExtractDateRangeAndAppFromArgs(connection, sceneName, args[0], before, after)
			if err != nil {
				return errors.Wrap(err, "Invalid argument(s)")
			}

			daysInRange := beforeTime.Sub(afterTime).Hours() / 24
			var pngs = make([][]byte, int(daysInRange))
			for i := 0; i < int(daysInRange); i++ {
				max := afterTime.AddDate(0, 0, i)
				ui.Log("Getting data between " + date.FormatDay(afterTime) + " and " + date.FormatDay(max))
				if err := chromedp.Run(ctxt, tasks(endpoint, a.SceneId, a.Id, date.FormatDay(afterTime), date.FormatDay(max), &pngs[i])); err != nil {
					return errors.Wrap(err, "Unable to browse data")
				}
			}

			width, height, err := calculateImageDimension(pngs[0])
			if err != nil {
				return errors.Wrap(err, "Unable to calculate image dimensions")
			}

			if err := createVideo(width, height, pngs, filename, fps); err != nil {
				return errors.Wrap(err, "Unable to create video")
			}
			ui.Ok()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:1233/", "Endpoint of the UI")
	cmd.Flags().StringVarP(&before, "before", "", "", "Fetch coupling before this day")
	cmd.Flags().StringVarP(&after, "after", "", "", "Fetch coupling after this day")
	cmd.Flags().StringVarP(&filename, "filename", "f", "storyboard" + date.Today() + ".avi", "storyboard file name")
	cmd.Flags().IntVarP(&fps, "fps", "", 8, "number of frames per second")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display the log information")

	return &cmd
}

func createVideo(width int, height int, pngs [][]byte, filename string, fps int) (error) {
	aw, err := mjpeg.New(filename, int32(width), int32(height), int32(fps))
	if err != nil {
		return errors.Wrap(err, "Unable to build video")
	}
	defer aw.Close()

	for i := 0; i < len(pngs); i++ {
		jpg, err := pngToJpeg(pngs[i])
		if err != nil {
			return errors.Wrap(err, "Unable to convert png to jpeg")
		}
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, jpg, nil); err != nil {
			return errors.Wrap(err, "Unable to encode jpeg image")
		}
		if err := aw.AddFrame(buf.Bytes()); err != nil {
			return errors.Wrap(err, "Unable to add frame to video")
		}
	}
	return nil
}

func calculateImageDimension(img []byte) (int, int, error) {
	jpg, err := pngToJpeg(img)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Unable to convert first png to jpeg")
	}
	width, height := jpg.Bounds().Size().X, jpg.Bounds().Size().Y
	return width, height, nil
}

func pngToJpeg(buf []byte) (*image.RGBA, error) {
	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode png image")
	}
	jpg := image.NewRGBA(img.Bounds())
	draw.Draw(jpg, jpg.Bounds(), img, img.Bounds().Min, draw.Src)
	return jpg, nil
}

func tasks(endpoint string, sceneId string, appId string, min string, max string, buf *[]byte) chromedp.Tasks {
	url := endpoint + `scenes/` + sceneId + `/apps/` + appId + `?after=` + min + `&before=` + max

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(".Chart", chromedp.ByQuery),
		chromedp.Screenshot(".Chart", buf, chromedp.NodeVisible),
	}

	return tasks
}
