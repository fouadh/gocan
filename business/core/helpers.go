package core

import (
	"com.fha.gocan/business/core/app"
	app2 "com.fha.gocan/business/data/store/app"
	"com.fha.gocan/foundation/date"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

func ExtractDateRangeAndAppFromArgs(connection *sqlx.DB, sceneName string, appName string, before string, after string) (app2.App, time.Time, time.Time, error) {
	beforeTime, err := date.ParseDay(before)
	if err != nil {
		return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Invalid before date")
	}

	afterTime, err := date.ParseDay(after)
	if err != nil {
		return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Invalid after date")
	}

	a, err := app.FindAppBySceneNameAndAppName(connection, sceneName, appName)
	if err != nil {
		return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Application cannot be retrieved")
	}

	return a, beforeTime, afterTime, nil
}
