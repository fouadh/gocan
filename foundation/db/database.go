package db

import (
	"com.fha.gocan/foundation/terminal"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type DataSource interface {
	GetConnection() (*sqlx.DB, error)
}

type SqlxDataSource struct {
	Dsn string
	Ui terminal.UI
}

func (ds *SqlxDataSource) GetConnection() (*sqlx.DB, error) {
	ds.Ui.Say("Connecting to the database...")
	db, err := sqlx.Connect("postgres", ds.Dsn)
	if err != nil {
		return nil, err
	}
	ds.Ui.Ok()
	return db, nil
}

func GetBulkInsertSQL(table string, columns []string, rowCount int) string {
	var b strings.Builder
	var cnt int

	columnCount := len(columns)

	b.Grow(40000) // Need to calculate, I'm too lazy))

	b.WriteString("INSERT INTO " + table + "(" + strings.Join(columns, ", ") + ") VALUES ")

	for i := 0; i < rowCount; i++ {
		b.WriteString("(")
		for j := 0; j < columnCount; j++ {
			cnt++
			b.WriteString("$")
			b.WriteString(strconv.Itoa(cnt))
			if j != columnCount-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(")")
		if i != rowCount-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}