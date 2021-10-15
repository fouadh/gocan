package core

import (
	"com.fha.gocan/business/core/app"
	"com.fha.gocan/business/core/commit"
	app2 "com.fha.gocan/business/data/store/app"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

func ExtractDateRangeAndAppFromArgs(connection *sqlx.DB, sceneName string, appName string, before string, after string) (app2.App, time.Time, time.Time, error) {
	a, err := app.FindAppByAppNameAndSceneName(connection, appName, sceneName)
	if err != nil {
		return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Application cannot be retrieved")
	}

	c := commit.NewCore(connection)
	beforeTime, afterTime, err := c.ExtractDateRangeFromArgs(a.Id, before, after)
	if err != nil {
		return app2.App{}, time.Time{}, time.Time{}, errors.Wrap(err, "Date range cannot be retrieved")
	}

	return a, beforeTime, afterTime, nil
}
