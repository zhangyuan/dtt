package db

import (
	"dtt/db/drivers/postgres"
	"dtt/db/sql_builder"
	"fmt"
)

func NewSQLBuilder(database *Database) (sql_builder.SQLBuilder, error) {
	if database.Driver == "postgres" {
		return &postgres.PostgresBuilder{}, nil
	} else {
		return nil, fmt.Errorf("database %s is not supported yet", database.Driver)
	}

}
