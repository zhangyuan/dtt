package spec

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Spec struct {
	Tables    []TableSpec
	TestSpecs []TestSpec `yaml:"tests"`
}

type TableSpec struct {
	Name    string
	Columns []ColumnSpec
}

type ColumnSpec struct {
	Name     string
	DataType string `yaml:"data_type"`
}

type TestSpec struct {
	Name           string
	Sources        []SourceSpec
	Transformation TransformationSpec
	ExpectedResult ExpectedResult `yaml:"expected_result"`
}

type TransformationSpec struct {
	Query     string
	QueryPath string
}

type SourceSpec struct {
	TableName string `yaml:"table_name"`
	Csv       string
}

type ExpectedResult struct {
	Csv string
}

func NewSpecFromPath(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	spec := Spec{}

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}
