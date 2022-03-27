package run

import (
	"dtt/db"
	"dtt/spec"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"reflect"
)

type Run struct {
	Spec         spec.TestSpec
	ExpectedData [][]string
	ActualData   [][]string
}

func (run *Run) IsOk() bool {
	return reflect.DeepEqual(run.ActualData, run.ExpectedData)
}

func Execute(database *db.Database, spec *spec.Spec) error {
	sqlBuilder, err := db.NewSQLBuilder(database)
	if err != nil {
		return err
	}

	reports := []Run{}

	for _, testSpec := range spec.TestSpecs {
		sql, err := sqlBuilder.BuildSQL(spec.Tables, testSpec)
		if err != nil {
			return err
		}

		fmt.Println(">>>>>" + testSpec.Name + "")
		fmt.Println(sql)

		actualData, err := database.Sql2Rows(sql)
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
