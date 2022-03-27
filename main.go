package main

import (
	"dtt/db"
	"dtt/run"
	"dtt/spec"
	"log"
	"strings"

	_ "github.com/lib/pq"

	"flag"
)

func main() {
	specPath := flag.String("spec", "", "path to spec")

	flag.Parse()

	if strings.TrimSpace(*specPath) == "" {
		log.Fatalln("spec is required.")
	}

	spec, err := spec.NewSpecFromPath(*specPath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	database, err := NewDatabase()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err := run.Execute(database, spec); err != nil {
		log.Fatalln(err)
	}
}

func NewDatabase() (*db.Database, error) {
	database, err := db.NewDatabaseFromEnv()

	if err != nil {
		return nil, err
	}

	if err := database.Connect(); err != nil {
		return nil, err
	}

	return database, nil

}
