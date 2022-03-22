package spec

type Spec struct {
	Tables []TableSpec
	Tests  []TestSpec
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
	Transformation string
	ExpectedResult ExpectedResult `yaml:"expected_result"`
}

type SourceSpec struct {
	TableName string `yaml:"table_name"`
	Csv       string
}

type ExpectedResult struct {
	Csv string
}
