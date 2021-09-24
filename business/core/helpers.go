package core

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/commit"
	app2 "com.fha.gocan/business/data/store/app"
	"com.fha.gocan/foundation/date"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

func ExtractDateRangeAndAppFromArgs(connection *sqlx.DB, sceneName string, appName string, before string, after string) (app2.App, time.Time, time.Time, error) {
	a, err := app.FindAppBySceneNameAndAppName(connection, sceneName, appName)
	if err != nil {
		return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Application cannot be retrieved")
	}

	c := commit.NewCore(connection)
	cr, rangeErr := c.QueryCommitRange(a.Id)

	if before == "" {
		if rangeErr != nil {
			return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Commit range cannot be retrieved")
		}
		before = date.FormatDay(cr.MaxDate)
	}
	beforeTime, err := date.ParseDay(before)
	if err != nil {
		return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Invalid before date")
	}

	if after == "" {
		if rangeErr != nil {
			return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Commit range cannot be retrieved")
		}
		after = date.FormatDay(cr.MinDate)
	}
	afterTime, err := date.ParseDay(after)
	if err != nil {
		return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Invalid after date")
	}

	return a, beforeTime, afterTime, nil
}
