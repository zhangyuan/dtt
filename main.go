package main

import (
	"datatest/internal/spec"
	"datatest/internal/sql_builder"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func loadSpec(path string) (*spec.Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	spec := spec.Spec{}

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}

func run() error {
	spec, err := loadSpec("./fixtures/spec.yaml")

	if err != nil {
		return err
	}

	for _, test := range spec.Tests {
		sql, err := sql_builder.BuildSQL(spec.Tables, test)

		if err != nil {
			return err
		}
		fmt.Printf("SQL: %+v\n", sql)
	}

	fmt.Printf("%+v\n", spec)

	return nil
}
