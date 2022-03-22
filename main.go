package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"trt/internal/spec"
	"trt/internal/sql_builder"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"gopkg.in/yaml.v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
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

func run() error {
	db, err := sqlx.Connect("postgres", "host=192.168.64.6 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	spec, err := loadSpec("./fixtures/spec.yaml")

	if err != nil {
		return err
	}

	reports := []Run{}

	for _, testSpec := range spec.TestSpecs {
		sql, err := sql_builder.BuildSQL(spec.Tables, testSpec)
		if err != nil {
			return err
		}

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
		if report.IsOk() {
			fmt.Print("[OK]")
		} else {
			isOk = false
			fmt.Print("[Failed]")
		}

		fmt.Printf("%s \n", report.Spec.Name)
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

	file, err := ioutil.TempFile("", "output.*.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var rows [][]string

	rows = append(rows, columnNames)

	if sqlRows.Next() {
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
