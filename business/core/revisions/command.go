package revisions

import (
	context "com.fha.gocan/foundation"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"time"
)

type Revision struct {
	Entity                      string
	NumberOfRevisions           int
	NumberOfAuthors             int
	NormalizedNumberOfRevisions float64
	Code                        int
}

func NewCommand(ctx *context.Context) *cobra.Command {
	var sceneName string

	cmd := cobra.Command{
		Use:  "revisions",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui := ctx.Ui
			ui.Say("Getting app revisions...")

			datasource := ctx.DataSource
			connection, err := datasource.GetConnection()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("The connection to the dabase could not be established: %v", err.Error()))
			}

			appName := args[0]
			appId, err := getAppId(appName, connection, sceneName)
			if err != nil {
				return err
			}

			sql := fmt.Sprintf(`select entity,
				   numberOfRevisions,
				   numberOfAuthors,
				   cast(numberOfRevisions as decimal) / MAX(numberOfRevisions) OVER () as normalizedNumberOfRevisions,
				   coalesce((select lines from cloc where file=entity and app_id=$1), 0) as code
					FROM (
					 SELECT file     as entity,
								count(distinct c.author) as numberOfAuthors,
							COUNT(*) as numberOfRevisions
					 FROM stats inner join commits c on c.id = stats.commit_id
					 and c.date between $3 and $2
					  WHERE stats.app_id=$1
					  AND FILE NOT LIKE '%%=>%%'
					 GROUP BY file
				 ) a
				ORDER BY %s desc             
			`, "numberOfRevisions")

			revisions := []Revision{}
			before := time.Now().AddDate(0, 0, 1)
			after, _ := time.Parse("2006-01-02", "1970-01-01")

			err = connection.Select(&revisions, sql, appId, before, after)

			if err != nil {
				ui.Failed("Cannot fetch revisions: " + err.Error())
				return err
			}

			ui.Ok()

			table := ui.Table([]string{
				"entity",
				"n-revs",
			})
			for _, revision := range revisions {
				table.Add(revision.Entity, fmt.Sprint(revision.NumberOfRevisions))
			}
			table.Print()
			return nil
		},
	}

	cmd.Flags().StringVarP(&sceneName, "scene", "s", "", "Scene name")
	return &cmd
}

func getAppId(appName string, connection *sqlx.DB, sceneName string) (string, error) {
	appId := ""
	if err := connection.Get(&appId, "select id from apps where name=$1 and scene_id=(select id from scenes where name=$2)", appName, sceneName); err != nil {
		return "", errors.Wrap(err, "Unable to retrieve matching app id")
	}
	if appId == "" {
		return "", errors.New("Application not found.")
	}
	return appId, nil
}
