package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	Driver         string // e.g. "postgres"
	DataSourceName string // e.g. "host=192.168.64.6 user=postgres password=postgres dbname=postgres sslmode=disable"
	DbConnection   *sqlx.DB
}

func (database *Database) Connect() error {
	dbConnection, err := sqlx.Connect(database.Driver, database.DataSourceName)
	if err != nil {
		return err
	}

	database.DbConnection = dbConnection

	return nil
}

func NewDatabaseFromEnv() (*Database, error) {
	databaseDriver, ok := os.LookupEnv("DATABASE_DRIVER")
	if ok == false {
		return nil, errors.New("DATABASE_DRIVER is empty")
	}

	dataSourceName, ok := os.LookupEnv("DATA_SOURCE_NAME")
	if ok == false {
		return nil, errors.New("DATA_SOURCE_NAME is empty")
	}

	return &Database{
		Driver:         databaseDriver,
		DataSourceName: dataSourceName,
	}, nil
}

func (db *Database) Sql2Rows(sql string) ([][]string, error) {
	sqlRows, err := db.DbConnection.Query(sql)

	if err != nil {
		return nil, err
	}

	columnNames, err := sqlRows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columnNames)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	var rows [][]string

	rows = append(rows, columnNames)

	for sqlRows.Next() {
		for i := range columnNames {
			valuePtrs[i] = &values[i]
		}

		if err := sqlRows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowValues := []string{}
		for _, value := range values {
			rowValues = append(rowValues, fmt.Sprint(value))
		}

		rows = append(rows, rowValues)
	}
	return rows, nil
}
