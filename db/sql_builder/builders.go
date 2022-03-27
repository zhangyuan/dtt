package sql_builder

import (
	"dtt/spec"
)

type SQLBuilder interface {
	BuildSQL(tables []spec.TableSpec, test spec.TestSpec) (string, error)
}
