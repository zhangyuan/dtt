package sql_builder

import (
	"bytes"
	"dtt/internal/spec"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
)

func BuildSQL(tables []spec.TableSpec, test spec.TestSpec) (string, error) {
	valueTables := []string{}
	for _, source := range test.Sources {
		valueTable, err := buildTable(&source, tables, test.Transformation.Query)
		if err != nil {
			return "", err
		}
		valueTables = append(valueTables, valueTable)
	}

	var buffer bytes.Buffer

	buffer.WriteString("WITH\n")
	buffer.WriteString(strings.Join(valueTables, ", \n"))
	buffer.WriteString("\n")

	query := strings.TrimSpace(test.Transformation.Query)
	if strings.ToUpper(query[0:4]) == "WITH" {
		buffer.WriteString(",")
		buffer.WriteString(query[4:])
	} else {
		buffer.WriteString(query)
	}

	return buffer.String(), nil
}

func findTable(tables *[]spec.TableSpec, tableName string) (*spec.TableSpec, error) {
	for _, table := range *tables {
		if table.Name == tableName {
			return &table, nil
		}
	}
	return nil, fmt.Errorf("table '%s' is not found", tableName)
}

func findColumn(columns *[]spec.ColumnSpec, columnName string) (*spec.ColumnSpec, error) {
	for _, column := range *columns {
		if columnName == column.Name {
			return &column, nil
		}
	}
	return nil, fmt.Errorf("column '%s' is not found", columnName)
}

func buildTable(source *spec.SourceSpec, tables []spec.TableSpec, transformation string) (string, error) {
	csvFile, err := os.Open(source.Csv)
	if err != nil {
		return "", err
	}
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()

	if err != nil {
		return "", err
	}

	defer csvFile.Close()

	if tableSpec, ok := lo.Find(tables, func(t spec.TableSpec) bool { return t.Name == source.TableName }); ok {
		return buildTableWithDataType(tableSpec, records)
	} else {
		return buildTableWithoutDataType(source.TableName, records), nil
	}
}

func buildTableWithDataType(tableSpec spec.TableSpec, records [][]string) (string, error) {
	valueExpressions := []string{}
	headersValues := []string{}
	for rowIndex, row := range records {
		if rowIndex == 0 {
			headersValues = row
			continue
		}

		rowValues := []string{}
		for rawValueIndex, rawValue := range row {
			if rowIndex == 1 {
				column, err := findColumn(&tableSpec.Columns, headersValues[rawValueIndex])

				if err != nil {
					return "", nil
				}

				value := fmt.Sprintf("'%s'::%s", rawValue, column.DataType)
				rowValues = append(rowValues, value)
			} else {
				value := fmt.Sprintf("'%s'", rawValue)
				rowValues = append(rowValues, value)
			}
		}

		valueExpression := fmt.Sprintf("(%s)", strings.Join(rowValues, ","))
		valueExpressions = append(valueExpressions, valueExpression)
	}

	template := `__dt_table_name__ (__dt_fields__) AS (
VALUES
__dt_values__)`

	template = strings.ReplaceAll(template, "__dt_table_name__", tableSpec.Name)
	template = strings.ReplaceAll(template, "__dt_fields__", strings.Join(headersValues, ","))
	template = strings.ReplaceAll(template, "__dt_values__", strings.Join(valueExpressions, ", \n"))

	return template, nil
}

func buildTableWithoutDataType(tableName string, rawRecords [][]string) string {
	records := []string{}
	columnNames := []string{}
	for rowIndex, rawRecord := range rawRecords {
		if rowIndex == 0 {
			columnNames = rawRecord
			continue
		}

		recordFields := []string{}
		for _, rawField := range rawRecord {
			recordFields = append(recordFields, rawField)
		}

		record := fmt.Sprintf("(%s)", strings.Join(recordFields, ", "))
		records = append(records, record)
	}

	template := `__dt_table_name__ (__dt_column_names__) AS (
VALUES
__dt_records__)`

	template = strings.ReplaceAll(template, "__dt_table_name__", tableName)
	template = strings.ReplaceAll(template, "__dt_column_names__", strings.Join(columnNames, ","))
	template = strings.ReplaceAll(template, "__dt_records__", strings.Join(records, ", \n"))

	return template
}
