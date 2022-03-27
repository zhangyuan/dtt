package main

import (
	"dtt/internal/spec"
	"dtt/internal/sql_builder"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"flag"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	specPath := flag.String("spec", "", "path to spec")

	flag.Parse()

	if strings.TrimSpace(*specPath) == "" {
		log.Fatalln("spec is required.")
	}

	database, err := NewDatabaseFromEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err := run(database, *specPath); err != nil {
		log.Fatalln(err)
	}

}

type Database struct {
	Driver         string // "postgres"
	DataSourceName string // "host=192.168.64.6 user=postgres password=postgres dbname=postgres sslmode=disable"
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

type Run struct {
	Spec         spec.TestSpec
	ExpectedData [][]string
	ActualData   [][]string
}

func (run *Run) IsOk() bool {
	return reflect.DeepEqual(run.ActualData, run.ExpectedData)
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

func run(database *Database, specPath string) error {
	db, err := sqlx.Connect(database.Driver, database.DataSourceName)

	if err != nil {
		return err
	}

	spec, err := loadSpec(specPath)
	if err != nil {
		return err
	}

	reports := []Run{}

	for _, testSpec := range spec.TestSpecs {
		sql, err := sql_builder.BuildSQL(spec.Tables, testSpec)
		if err != nil {
			return err
		}

		fmt.Println(">>>>>" + testSpec.Name + "")
		fmt.Println(sql)

		actualData, err := Sql2Rows(db, sql)
		if err != nil {
			return err
		}

		report, err := BuildReport(testSpec, actualData)
		if err != nil {
			return err
		}
		reports = append(reports, *report)
	}

	isOk := true
	for _, report := range reports {
		fmt.Println("===================>" + report.Spec.Name)
		if report.IsOk() {
			fmt.Println("[Ok]")
		} else {
			isOk = false
			fmt.Println("[Failed]")
			fmt.Println("Expected:")
			fmt.Println(report.ExpectedData)

			fmt.Println("Actual:")
			fmt.Println(report.ActualData)
		}
	}

	if isOk {
		return nil
	}
	return errors.New("some of the tests failed")
}

func BuildReport(testSpec spec.TestSpec, actualData [][]string) (*Run, error) {
	expectedDataFile, err := os.Open(testSpec.ExpectedResult.Csv)
	if err != nil {
		return nil, err
	}

	expectedData, err := csv.NewReader(expectedDataFile).ReadAll()
	if err != nil {
		return nil, err
	}

	return &Run{
		Spec:         testSpec,
		ExpectedData: expectedData,
		ActualData:   actualData,
	}, nil
}

func Sql2Rows(db *sqlx.DB, sql string) ([][]string, error) {
	sqlRows, err := db.Query(sql)

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
