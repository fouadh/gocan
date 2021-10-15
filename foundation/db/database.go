package db

import (
	"com.fha.gocan/foundation/terminal"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"reflect"
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
	ds.Ui.Log("Connecting to the database...")
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
	b.WriteString(" ON CONFLICT DO NOTHING")
	return b.String()
}

func NamedQuerySlice(connection *sqlx.DB, query string, data interface{}, dest interface{}) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	rows, err := connection.NamedQuery(query, data)
	if err != nil {
		return err
	}
	defer rows.Close()

	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err := rows.StructScan(v.Interface()); err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}

	return nil
}

func NamedQueryStruct(connection *sqlx.DB, query string, data interface{}, dest interface{}) error {
	rows, err := connection.NamedQuery(query, data)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New("Unable to retrieve record")
	}
	defer rows.Close()

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}
